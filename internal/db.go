package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"imap-sync/config"
	"imap-sync/logger"
	"reflect"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var (
	admin_name string
	admin_pass string
	DB_path    string
)

var db *sql.DB
var log = logger.Log

func InitDb() error {
	if _, err := getEncryptionKey(); err != nil {
		return fmt.Errorf("encryption key required: %w", err)
	}

	admin_name = config.Conf.DatabaseInfo.AdminName
	admin_pass = config.Conf.DatabaseInfo.AdminPass
	DB_path = config.Conf.DatabaseInfo.DatabasePath

	var err error
	db, err = sql.Open("sqlite3", DB_path)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	// Initialize Login Table

	initStmt := `
	CREATE TABLE IF NOT EXISTS Users (
	id INTEGER PRIMARY KEY,
	username VARCHAR(64) NULL,
	password VARCHAR(64) NULL
	);
	`
	_, err = db.Exec(initStmt)
	if err != nil {
		return fmt.Errorf("error creating login table: %w", err)
	}

	initTaskTable := `
	CREATE TABLE IF NOT EXISTS Tasks (
	id INTEGER PRIMARY KEY,
	source_account VARCHAR(255) NULL,
	source_server VARCHAR(255) NULL,
	source_password VARCHAR(255) NULL,
	destination_account VARCHAR(255) NULL,
	destination_server VARCHAR(255) NULL,
	destination_password VARCHAR(255) NULL,
	started_at INTEGER NULL,
	ended_at INTEGER NULL,
	status VARCHAR(64) NULL,
	logfile VARCHAR(255) NULL,
	messages_copied INTEGER DEFAULT 0,
	bytes_transferred BIGINT DEFAULT 0,
	error_detail TEXT
	);
	`

	_, err = db.Exec(initTaskTable)

	if err != nil {
		return fmt.Errorf("error creating task table: %w", err)
	}

	initAuditLogTable := `
	CREATE TABLE IF NOT EXISTS AuditLog (
	id INTEGER PRIMARY KEY,
	user VARCHAR(64) NULL,
	action VARCHAR(64) NULL,
	details TEXT,
	ip_address VARCHAR(64) NULL,
	created_at INTEGER NULL
	);
	`
	_, err = db.Exec(initAuditLogTable)
	if err != nil {
		return fmt.Errorf("error creating audit log table: %w", err)
	}

	initSessionTable := `
	CREATE TABLE IF NOT EXISTS Sessions (
	id INTEGER PRIMARY KEY,
	user VARCHAR(64) NULL,
	ip_address VARCHAR(64) NULL,
	user_agent VARCHAR(255) NULL,
	created_at INTEGER NULL,
	last_activity INTEGER NULL,
	active INTEGER DEFAULT 1
	);
	`
	_, err = db.Exec(initSessionTable)
	if err != nil {
		return fmt.Errorf("error creating sessions table: %w", err)
	}

	var exists bool
	err = db.QueryRow("SELECT exists (SELECT 1 FROM users WHERE username = ?)", admin_name).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if admin exists: %w", err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(admin_pass), 14)
	if err != nil {
		return fmt.Errorf("error hashing admin password: %w", err)
	}

	if !exists {
		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", admin_name, password)
		if err != nil {
			return fmt.Errorf("error creating admin user: %w", err)
		}
	} else {
		dbPass, err := GetPassword(admin_name)
		if err != nil {
			return fmt.Errorf("error getting admin password: %w", err)
		}
		if !reflect.DeepEqual(password, []byte(dbPass)) {
			changePassword(admin_name, string(password))
		}
	}

	return nil
}

func changePassword(username string, newPassword string) error {
	stmt, err := db.Prepare("UPDATE users SET password = ? WHERE username = ?")
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(newPassword, username)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func GetPassword(username string) (string, error) {
	var password string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("no user found with username %s", username)
		}
		return "", fmt.Errorf("error getting password: %w", err)
	}
	return password, nil
}

func AddTaskToDB(task *Task) error {
	log.Info("Adding task to database")
	encSourcePass, err := EncryptString(task.SourcePassword)
	if err != nil {
		return fmt.Errorf("error encrypting source password: %w", err)
	}
	encDestPass, err := EncryptString(task.DestinationPassword)
	if err != nil {
		return fmt.Errorf("error encrypting destination password: %w", err)
	}

	stmt, err := db.Prepare("INSERT INTO tasks(source_account, source_server, source_password, destination_account, destination_server, destination_password, started_at, ended_at, status, logfile) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.SourceAccount, task.SourceServer, encSourcePass, task.DestinationAccount, task.DestinationServer, encDestPass, task.StartedAt, task.EndedAt, task.Status, task.LogFile)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func updateTaskStatus(task *Task, status string) error {
	timeUnix := time.Now().Unix()

	if status == "In Progress" {
		stmt, err := db.Prepare("UPDATE tasks SET started_at = ?, status = ? WHERE id = ?")
		if err != nil {
			return fmt.Errorf("error preparing statement: %w", err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(timeUnix, status, task.ID)
		if err != nil {
			return fmt.Errorf("error executing statement: %w", err)
		}
		task.StartedAt = timeUnix
	} else {
		if task.Status == "In Progress" {
			stmt, err := db.Prepare("UPDATE tasks SET ended_at = ?, status = ? WHERE id = ?")
			if err != nil {
				return fmt.Errorf("error preparing statement: %w", err)
			}
			defer stmt.Close()
			_, err = stmt.Exec(timeUnix, status, task.ID)
			if err != nil {
				return fmt.Errorf("error executing statement: %w", err)
			}
			task.EndedAt = timeUnix
		} else {
			stmt, err := db.Prepare("UPDATE tasks SET status = ? WHERE id = ?")
			if err != nil {
				return fmt.Errorf("error preparing statement: %w", err)
			}
			defer stmt.Close()
			_, err = stmt.Exec(status, task.ID)
			if err != nil {
				return fmt.Errorf("error executing statement: %w", err)
			}
		}
	}

	task.Status = status

	return nil
}

func updateTaskLogFile(task *Task, logFile string) error {
	stmt, err := db.Prepare("UPDATE tasks SET logfile = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(logFile, task.ID)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	task.LogFile = logFile

	return nil
}

func InitializeQueueFromDB() error {
	log.Info("Initializing queue from database")
	rows, err := db.Query("SELECT id, source_account, source_server, source_password, destination_account, destination_server, destination_password, status, logfile FROM tasks")
	if err != nil {
		return fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.SourceAccount, &task.SourceServer, &task.SourcePassword, &task.DestinationAccount, &task.DestinationServer, &task.DestinationPassword, &task.Status, &task.LogFile)
		if err != nil {
			return fmt.Errorf("error scanning row: %w", err)
		}

		sourcePass, err := DecryptString(task.SourcePassword)
		if err != nil {
			if errors.Is(err, errValueNotEncrypted) {
				encSource, encErr := EncryptString(task.SourcePassword)
				if encErr != nil {
					return fmt.Errorf("error encrypting legacy source password: %w", encErr)
				}
				encDest, encErr := EncryptString(task.DestinationPassword)
				if encErr != nil {
					return fmt.Errorf("error encrypting legacy destination password: %w", encErr)
				}
				if err := updateTaskPasswords(task.ID, encSource, encDest); err != nil {
					return fmt.Errorf("error updating legacy passwords: %w", err)
				}
			} else {
				return fmt.Errorf("error decrypting source password: %w", err)
			}
		} else {
			task.SourcePassword = sourcePass
		}

		destPass, err := DecryptString(task.DestinationPassword)
		if err != nil {
			if errors.Is(err, errValueNotEncrypted) {
				encSource, encErr := EncryptString(task.SourcePassword)
				if encErr != nil {
					return fmt.Errorf("error encrypting legacy source password: %w", encErr)
				}
				encDest, encErr := EncryptString(task.DestinationPassword)
				if encErr != nil {
					return fmt.Errorf("error encrypting legacy destination password: %w", encErr)
				}
				if err := updateTaskPasswords(task.ID, encSource, encDest); err != nil {
					return fmt.Errorf("error updating legacy passwords: %w", err)
				}
			} else {
				return fmt.Errorf("error decrypting destination password: %w", err)
			}
		} else {
			task.DestinationPassword = destPass
		}

		if task.Status == "In Progress" {
			task.Status = "Cancelled"
		}

		queue.PushFront(&task)
	}

	return nil
}

func updateTaskPasswords(taskID int, sourcePassword string, destinationPassword string) error {
	stmt, err := db.Prepare("UPDATE tasks SET source_password = ?, destination_password = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("error preparing update: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(sourcePassword, destinationPassword, taskID)
	if err != nil {
		return fmt.Errorf("error executing update: %w", err)
	}

	return nil
}

func InitSettingsTable() error {
	initStmt := `
	CREATE TABLE IF NOT EXISTS Settings (
	id INTEGER PRIMARY KEY,
	source_server VARCHAR(255) NULL,
	source_account_prefix VARCHAR(255) NULL,
	source_use_tls INTEGER DEFAULT 0,
	destination_server VARCHAR(255) NULL,
	destination_account_prefix VARCHAR(255) NULL,
	destination_use_tls INTEGER DEFAULT 0
	);
	`
	_, err := db.Exec(initStmt)
	if err != nil {
		return fmt.Errorf("error creating settings table: %w", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM Settings").Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking settings count: %w", err)
	}

	if count == 0 {
		_, err = db.Exec("INSERT INTO Settings(source_server, source_account_prefix, source_use_tls, destination_server, destination_account_prefix, destination_use_tls) VALUES(?, ?, ?, ?, ?, ?)",
			config.Conf.SourceAndDestination.SourceServer,
			config.Conf.SourceAndDestination.SourceMail,
			0,
			config.Conf.SourceAndDestination.DestinationServer,
			config.Conf.SourceAndDestination.DestinationMail,
			0)
		if err != nil {
			return fmt.Errorf("error inserting default settings: %w", err)
		}
	}

	return nil
}

func GetSettings() (*Settings, error) {
	var s Settings
	var sourceUseTLS, destUseTLS int
	err := db.QueryRow("SELECT source_server, source_account_prefix, source_use_tls, destination_server, destination_account_prefix, destination_use_tls FROM Settings WHERE id = 1").Scan(
		&s.SourceServer, &s.SourceAccountPrefix, &sourceUseTLS, &s.DestinationServer, &s.DestinationAccountPrefix, &destUseTLS)
	if err != nil {
		return nil, fmt.Errorf("error getting settings: %w", err)
	}
	s.SourceUseTLS = sourceUseTLS == 1
	s.DestinationUseTLS = destUseTLS == 1
	return &s, nil
}

func UpdateSettings(s *Settings) error {
	stmt, err := db.Prepare("UPDATE Settings SET source_server = ?, source_account_prefix = ?, source_use_tls = ?, destination_server = ?, destination_account_prefix = ?, destination_use_tls = ? WHERE id = 1")
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	sourceTLS := 0
	if s.SourceUseTLS {
		sourceTLS = 1
	}
	destTLS := 0
	if s.DestinationUseTLS {
		destTLS = 1
	}

	_, err = stmt.Exec(s.SourceServer, s.SourceAccountPrefix, sourceTLS, s.DestinationServer, s.DestinationAccountPrefix, destTLS)
	if err != nil {
		return fmt.Errorf("error updating settings: %w", err)
	}
	return nil
}

type DashboardStats struct {
	TotalTasks      int     `json:"total_tasks"`
	CompletedTasks  int     `json:"completed_tasks"`
	FailedTasks     int     `json:"failed_tasks"`
	PendingTasks    int     `json:"pending_tasks"`
	InProgressTasks int     `json:"in_progress_tasks"`
	CancelledTasks  int     `json:"cancelled_tasks"`
	SuccessRate     float64 `json:"success_rate"`
}

func GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	rows, err := db.Query("SELECT status, COUNT(*) as count FROM tasks GROUP BY status")
	if err != nil {
		return nil, fmt.Errorf("error querying stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats.TotalTasks += count
		switch status {
		case "Done":
			stats.CompletedTasks = count
		case "Error":
			stats.FailedTasks = count
		case "Pending":
			stats.PendingTasks = count
		case "In Progress":
			stats.InProgressTasks = count
		case "Cancelled":
			stats.CancelledTasks = count
		}
	}

	if stats.TotalTasks > 0 {
		stats.SuccessRate = float64(stats.CompletedTasks) / float64(stats.TotalTasks-stats.PendingTasks-stats.InProgressTasks) * 100
	}

	return stats, nil
}

type AuditLogEntry struct {
	ID        int
	User      string
	Action    string
	Details   string
	IPAddress string
	CreatedAt int64
}

type Session struct {
	ID           int
	User         string
	IPAddress    string
	UserAgent    string
	CreatedAt    int64
	LastActivity int64
	Active       bool
}

func AddAuditLog(user, action, details, ipAddress string) error {
	stmt, err := db.Prepare("INSERT INTO AuditLog(user, action, details, ip_address, created_at) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user, action, details, ipAddress, time.Now().Unix())
	return err
}

func GetAuditLog(limit int) ([]AuditLogEntry, error) {
	rows, err := db.Query("SELECT id, user, action, details, ip_address, created_at FROM AuditLog ORDER BY created_at DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []AuditLogEntry
	for rows.Next() {
		var e AuditLogEntry
		if err := rows.Scan(&e.ID, &e.User, &e.Action, &e.Details, &e.IPAddress, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func AddSession(user, ipAddress, userAgent string) (int, error) {
	stmt, err := db.Prepare("INSERT INTO Sessions(user, ip_address, user_agent, created_at, last_activity, active) VALUES(?, ?, ?, ?, ?, 1)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	now := time.Now().Unix()
	result, err := stmt.Exec(user, ipAddress, userAgent, now, now)
	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()
	return int(id), nil
}

func UpdateSessionActivity(sessionID int) error {
	_, err := db.Exec("UPDATE Sessions SET last_activity = ? WHERE id = ?", time.Now().Unix(), sessionID)
	return err
}

func GetActiveSessions() ([]Session, error) {
	rows, err := db.Query("SELECT id, user, ip_address, user_agent, created_at, last_activity, active FROM Sessions WHERE active = 1 ORDER BY last_activity DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var s Session
		var active int
		if err := rows.Scan(&s.ID, &s.User, &s.IPAddress, &s.UserAgent, &s.CreatedAt, &s.LastActivity, &active); err != nil {
			return nil, err
		}
		s.Active = active == 1
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func TerminateSession(sessionID int) error {
	_, err := db.Exec("UPDATE Sessions SET active = 0 WHERE id = ?", sessionID)
	return err
}

func TerminateAllSessions() error {
	_, err := db.Exec("UPDATE Sessions SET active = 0")
	return err
}
