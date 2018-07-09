package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// DB implement database struct
type DB struct {
	Conn *sql.DB
	Path string
}

// NewDB populate and return new *DB instance
func NewDB(path string) (*DB, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return &DB{conn, path}, nil
}

func checkConfig() {
	dbPath := os.Getenv("HOME") + "/.PM/db.sqlite"
	if !pathExists(dbPath) {
		initBase()
		os.Exit(0)
	}
}

func initBase() error {
	pmDir := os.Getenv("HOME") + "/.PM"
	fmt.Println("creating configuration directory...")
	err := mkdir(pmDir)
	if err != nil {
		return err
	}

	pass := generate(16)
	dbFile := getPrefix() + pass
	err = mkfile(dbFile)
	if err != nil {
		return err
	}

	db, err := NewDB(dbFile)
	if err != nil {
		return err
	}
	defer db.Conn.Close()

	fmt.Println("creating database scheme...")
	cmd := `
CREATE TABLE passwords(
'id' INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
name VARCHAR(32) NOT NULL,
username VARCHAR(32) NOT NULL,
resource TEXT NOT NULL,
password VARCHAR(32) NOT NULL,
comment TEXT NOT NULL,
'group' VARCHAR(32) NOT NULL
)`
	_, err = db.Conn.Exec(cmd)
	if err != nil {
		return err
	}

	fmt.Println("encrypting database...")
	return encrypt(dbFile)
}

func (db *DB) doQuery(query string, args ...interface{}) error {
	defer db.Conn.Close()

	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}

	cmd, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer cmd.Close()

	_, err = cmd.Exec(args...)
	if err != nil {
		return err
	}
	tx.Commit()

	return encrypt(db.Path)
}

func (db *DB) doSelect(query string, args ...interface{}) ([]*password, error) {
	defer func() {
		db.Conn.Close()
		err := rmfile(db.Path)
		if err != nil {
			fmt.Println("failed to remove unencrypted database:", err)
		}
	}()

	rows, err := db.Conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passwords []*password
	for rows.Next() {
		passwd := &password{}
		err = rows.Scan(
			&passwd.id,
			&passwd.name,
			&passwd.username,
			&passwd.resource,
			&passwd.password,
			&passwd.comment,
			&passwd.group,
		)

		if err != nil {
			return nil, err
		}

		passwords = append(passwords, passwd)
	}

	return passwords, nil
}
