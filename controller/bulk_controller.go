package controller

import (
	"encoding/json"
	"imap-sync/internal"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

const MaxConcurrentMigrations = 5

type BulkMigrationRequest struct {
	SourceServer      string `json:"source_server"`
	SourceUseTLS      bool   `json:"source_use_tls"`
	DestinationServer string `json:"destination_server"`
	DestinationUseTLS bool   `json:"destination_use_tls"`
	Accounts          string `json:"accounts"`
}

type BulkAccountStatus struct {
	Index           int    `json:"index"`
	SourceUser      string `json:"source_user"`
	DestinationUser string `json:"destination_user"`
	Status          string `json:"status"`
	Error           string `json:"error"`
	Progress        int    `json:"progress"`
	TotalMessages   int    `json:"total_messages"`
	CopiedMessages  int    `json:"copied_messages"`
}

type BulkMigrationStatus struct {
	ID       int                 `json:"id"`
	Status   string              `json:"status"`
	Accounts []BulkAccountStatus `json:"accounts"`
}

func HandleBulkMigration(ctx *gin.Context) {
	var req BulkMigrationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accounts := internal.ParseBulkAccounts(req.Accounts)
	if len(accounts) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No accounts found in CSV"})
		return
	}

	bm := &internal.BulkMigration{
		SourceServer:      req.SourceServer,
		SourceUseTLS:      req.SourceUseTLS,
		DestinationServer: req.DestinationServer,
		DestinationUseTLS: req.DestinationUseTLS,
		Accounts:          accounts,
	}

	migrationID := internal.AddBulkMigration(bm)

	go runBulkMigrationParallel(migrationID, bm)

	ctx.JSON(http.StatusOK, gin.H{"migration_id": migrationID, "total_accounts": len(accounts)})
}

func runBulkMigrationParallel(migrationID int, bm *internal.BulkMigration) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, MaxConcurrentMigrations)

	for i := range bm.Accounts {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			account := bm.Accounts[index]

			validation := &internal.AccountValidation{
				SourceServer:        bm.SourceServer,
				SourceUser:          account.SourceUser,
				SourcePassword:      account.SourcePassword,
				SourceUseTLS:        bm.SourceUseTLS,
				DestinationServer:   bm.DestinationServer,
				DestinationUser:     account.DestinationUser,
				DestinationPassword: account.DestinationPassword,
				DestinationUseTLS:   bm.DestinationUseTLS,
			}

			internal.UpdateBulkAccount(migrationID, index, "validating", 0, 0, 0, "")

			result := internal.ValidateAccount(validation)

			if !result.SourceValid {
				internal.UpdateBulkAccount(migrationID, index, "failed", 0, 0, 0, "Source: "+result.SourceError)
				return
			}
			if !result.DestinationValid {
				internal.UpdateBulkAccount(migrationID, index, "failed", 0, 0, 0, "Destination: "+result.DestinationError)
				return
			}

			internal.UpdateBulkAccount(migrationID, index, "syncing", 0, 0, 0, "")

			task := &internal.Task{
				SourceAccount:       account.SourceUser,
				SourceServer:        bm.SourceServer,
				SourcePassword:      account.SourcePassword,
				DestinationAccount:  account.DestinationUser,
				DestinationServer:   bm.DestinationServer,
				DestinationPassword: account.DestinationPassword,
				Status:              "In Progress",
			}

			if err := internal.AddTaskToDB(task); err != nil {
				internal.UpdateBulkAccount(migrationID, index, "failed", 0, 0, 0, "Failed to persist task")
				return
			}
			internal.TaskChan() <- *task

			internal.UpdateBulkAccount(migrationID, index, "completed", 100, 0, 0, "")
		}(i)
	}

	wg.Wait()
}

func HandleBulkMigrationStatus(ctx *gin.Context) {
	migrationID := ctx.Query("id")
	if migrationID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "migration_id required"})
		return
	}

	var id int
	json.Unmarshal([]byte(migrationID), &id)

	bm := internal.GetBulkMigration(id)
	if bm == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Migration not found"})
		return
	}

	status := BulkMigrationStatus{
		ID:     bm.ID,
		Status: bm.Status,
	}

	for i, acc := range bm.Accounts {
		status.Accounts = append(status.Accounts, BulkAccountStatus{
			Index:           i,
			SourceUser:      acc.SourceUser,
			DestinationUser: acc.DestinationUser,
			Status:          acc.Status,
			Error:           acc.Error,
			Progress:        acc.Progress,
			TotalMessages:   acc.TotalMessages,
			CopiedMessages:  acc.CopiedMessages,
		})
	}

	ctx.JSON(http.StatusOK, status)
}
