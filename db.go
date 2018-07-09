package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// DB implement database struct
type DB struct {
	// Database file path
	Path string
}

// NewDB populate and return new *DB instance
func NewDB() *DB {
	return &DB{
		Path: os.Getenv("HOME") + "/.PM/db.sqlite",
	}
}

func connect(path string) (*sql.DB, error) {
	return sql.Open("sqlite3", path)
}

func checkConfig() {
	dbPath := os.Getenv("HOME") + "/.PM/db.sqlite"
	if !pathExists(dbPath) {
		init_base()
		os.Exit(0)
	}
}

func init_base() {
	pmDir := os.Getenv("HOME") + "/.PM"
	if pathExists(pmDir) {
		fmt.Println("removing old directory...")
		err := cmd("rm", "-rf", pmDir)
		if err != nil {
			fmt.Println("failed to remove directory", err)
			return
		}
	}

	fmt.Println("creating configuration directory...")
	err := mkdir(pmDir)
	if err != nil {
		fmt.Println("failed to create configuration directory", err)
		return
	}

	pass := generate(16)
	dbFile := "/tmp/" + pass
	err = cmd("touch", dbFile)
	if err != nil {
		fmt.Println("failed to create database file", err)
		return
	}

	conn, err := connect(dbFile)
	if err != nil {
		fmt.Println("failed to open the database file", err)
		return
	}
	defer conn.Close()

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
	_, err = conn.Exec(cmd)
	if err != nil {
		fmt.Println("failed to create db scheme", err)
		return
	}

	fmt.Println("encrypting database...")
	encrypt(dbFile)
}

func dbQuery(path string, query string, args ...interface{}) error {
	conn, err := connect(path)
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.Begin()
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

	return nil
}

func selectByName(name string) (*password, error) {
	decrypted, err := decrypt()
	if err != nil {
		return nil, err
	}

	conn, err := connect(decrypted)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.Query("select id, name, username, resource, password,"+
		" comment, `group` from passwords where name=?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, err
	}

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

	return passwd, nil
}

func selectByGroup(group string) ([]*password, error) {
	decrypted, err := decrypt()
	if err != nil {
		return nil, err
	}

	conn, err := connect(decrypted)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.Query("select id, name, username, resource, password,"+
		" comment, `group` from passwords where `group`=?", group)
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

func selectAll() ([]*password, error) {
	decrypted, err := decrypt()
	if err != nil {
		return nil, err
	}

	conn, err := connect(decrypted)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.Query("select * from passwords")

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
