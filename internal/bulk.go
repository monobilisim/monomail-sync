package internal

import (
	"sync"
	"time"
)

var (
	bulkMigrations   = make(map[int]*BulkMigration)
	bulkMigrationsMu sync.RWMutex
	bulkMigrationID  int
)

func AddBulkMigration(bm *BulkMigration) int {
	bulkMigrationsMu.Lock()
	defer bulkMigrationsMu.Unlock()
	bulkMigrationID++
	bm.ID = bulkMigrationID
	bm.CreatedAt = time.Now().Unix()
	bm.Status = "pending"
	bulkMigrations[bulkMigrationID] = bm
	return bulkMigrationID
}

func GetBulkMigration(id int) *BulkMigration {
	bulkMigrationsMu.RLock()
	defer bulkMigrationsMu.RUnlock()
	return bulkMigrations[id]
}

func UpdateBulkAccount(migrationID int, accountIndex int, status string, progress int, total int, copied int, errMsg string) {
	bulkMigrationsMu.Lock()
	defer bulkMigrationsMu.Unlock()

	if bm, ok := bulkMigrations[migrationID]; ok {
		if accountIndex < len(bm.Accounts) {
			bm.Accounts[accountIndex].Status = status
			bm.Accounts[accountIndex].Progress = progress
			bm.Accounts[accountIndex].TotalMessages = total
			bm.Accounts[accountIndex].CopiedMessages = copied
			bm.Accounts[accountIndex].Error = errMsg
		}
	}
}

func GetAllBulkMigrations() []BulkMigration {
	bulkMigrationsMu.RLock()
	defer bulkMigrationsMu.RUnlock()

	var result []BulkMigration
	for _, bm := range bulkMigrations {
		result = append(result, *bm)
	}
	return result
}
