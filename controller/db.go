package controller

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"imap-sync/internal"
	"reflect"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var log = internal.Log

var (
	admin_name = flag.String("admin_name", "admin", "Admin username")
	admin_pass = flag.String("admin_pass", "admin", "Admin password")
)

var db *sql.DB

func InitDb() error {
	var err error
	db, err = sql.Open("sqlite3", "./db.db")
	if err != nil {
		return err
	}
	defer db.Close()

	initStmt := `
	CREATE TABLE IF NOT EXISTS Users (
	id INTEGER PRIMARY KEY,
	username VARCHAR(64) NULL,
	password VARCHAR(64) NULL
	);
	`
	_, err = db.Exec(initStmt)
	if err != nil {
		return err
	}

	var exists bool
	err = db.QueryRow("SELECT exists (SELECT 1 FROM users WHERE username = ?)", admin_name).Scan(&exists)
	if err != nil {
		return err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(*admin_pass), 14)
	if err != nil {
		return err
	}

	if !exists {
		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", admin_name, password)
		if err != nil {
			return err
		}
	} else {
		dbPass, err := getPassword(*admin_name)
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(password, []byte(dbPass)) {
			changePassword(*admin_name, string(password))
		}
	}

	return nil
}

func changePassword(username string, newPassword string) error {
	var err error
	db, err = sql.Open("sqlite3", "./db.db")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE users SET password = ? WHERE username = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newPassword, username)
	if err != nil {
		return err
	}

	return nil
}

func getPassword(username string) (string, error) {
	var err error
	db, err = sql.Open("sqlite3", "./db.db")
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
		return "", err
	}
	return password, nil
}
