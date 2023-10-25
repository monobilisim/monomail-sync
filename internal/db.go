package internal

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"reflect"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var (
	admin_name = flag.String("admin_name", "admin", "Admin username")
	admin_pass = flag.String("admin_pass", "admin", "Admin password")
	DB_path    = flag.String("db_path", "./db.db", "Path to database")
)

var db *sql.DB

func InitDb() error {

	var err error
	db, err = sql.Open("sqlite3", *DB_path)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	defer db.Close()

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
	source_account VARCHAR(64) NULL,
	source_server VARCHAR(64) NULL,
	source_password VARCHAR(64) NULL,
	destination_account VARCHAR(64) NULL,
	destination_server VARCHAR(64) NULL,
	destination_password VARCHAR(64) NULL,
	started_at INTEGER NULL,
	ended_at INTEGER NULL,
	status VARCHAR(64) NULL,
	logfile VARCHAR(64) NULL
	);
	`

	_, err = db.Exec(initTaskTable)

	if err != nil {
		return fmt.Errorf("error creating task table: %w", err)
	}

	var exists bool
	err = db.QueryRow("SELECT exists (SELECT 1 FROM users WHERE username = ?)", admin_name).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if admin exists: %w", err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(*admin_pass), 14)
	if err != nil {
		return fmt.Errorf("error hashing admin password: %w", err)
	}

	if !exists {
		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", admin_name, password)
		if err != nil {
			return fmt.Errorf("error creating admin user: %w", err)
		}
	} else {
		dbPass, err := GetPassword(*admin_name)
		if err != nil {
			return fmt.Errorf("error getting admin password: %w", err)
		}
		if !reflect.DeepEqual(password, []byte(dbPass)) {
			changePassword(*admin_name, string(password))
		}
	}

	return nil
}

func changePassword(username string, newPassword string) error {
	var err error
	db, err = sql.Open("sqlite3", *DB_path)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}
	defer db.Close()

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
	var err error
	db, err = sql.Open("sqlite3", *DB_path)
	if err != nil {
		return "", err
	}
	defer db.Close()
	var password string
	err = db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&password)
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
	var err error
	db, err = sql.Open("sqlite3", *DB_path)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO tasks(source_account, source_server, source_password, destination_account, destination_server, destination_password, started_at, ended_at, status, logfile) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.SourceAccount, task.SourceServer, task.SourcePassword, task.DestinationAccount, task.DestinationServer, task.DestinationPassword, task.StartedAt, task.EndedAt, task.Status, task.LogFile)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func updateTaskStatus(task *Task, status string) error {
	var err error
	db, err = sql.Open("sqlite3", *DB_path)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}
	defer db.Close()

	var stmt *sql.Stmt
	timeUnix := time.Now().Unix()

	if status == "In Progress" {
		stmt, err = db.Prepare("UPDATE tasks SET started_at = ?, status = ? WHERE id = ?")
		if err != nil {
			return fmt.Errorf("error preparing statement: %w", err)
		}
		_, err = stmt.Exec(timeUnix, status, task.ID)
		if err != nil {
			return fmt.Errorf("error executing statement: %w", err)
		}
		task.StartedAt = timeUnix
	} else {
		if task.Status == "In Progress" {
			stmt, err = db.Prepare("UPDATE tasks SET ended_at = ?, status = ? WHERE id = ?")
			if err != nil {
				return fmt.Errorf("error preparing statement: %w", err)
			}
			_, err = stmt.Exec(timeUnix, status, task.ID)
			if err != nil {
				return fmt.Errorf("error executing statement: %w", err)
			}
			task.EndedAt = timeUnix
		} else {
			stmt, err = db.Prepare("UPDATE tasks SET status = ? WHERE id = ?")
			if err != nil {
				return fmt.Errorf("error preparing statement: %w", err)
			}
			_, err = stmt.Exec(status, task.ID)
			if err != nil {
				return fmt.Errorf("error executing statement: %w", err)
			}

		}
	}

	defer stmt.Close()

	task.Status = status

	return nil
}

func updateTaskLogFile(task *Task, logFile string) error {
	var err error
	db, err = sql.Open("sqlite3", *DB_path)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}
	defer db.Close()

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
	var err error
	db, err = sql.Open("sqlite3", *DB_path)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}
	defer db.Close()

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

		if task.Status == "In Progress" {
			task.Status = "Cancelled"
		}

		queue.PushFront(&task)
	}

	return nil
}
