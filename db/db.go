package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/himidori/pm/utils"
	_ "github.com/mattn/go-sqlite3"
)

// DB implement database struct
type DB struct {
	Conn *sql.DB
	Path string
}

// NewDB populate and return new *DB instance
func newDB(path string) (*DB, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return &DB{conn, path}, nil
}

func InitBase() error {
	pmDir := os.Getenv("HOME") + "/.PM"
	fmt.Println("creating configuration directory...")
	err := utils.Mkdir(pmDir)
	if err != nil {
		return err
	}

	pass := GeneratePassword(16)
	dbFile := utils.GetPrefix() + pass
	err = utils.Mkfile(dbFile)
	if err != nil {
		return err
	}

	db, err := newDB(dbFile)
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

func (db *DB) doSelect(query string, args ...interface{}) ([]*Password, error) {
	defer func() {
		db.Conn.Close()
		err := utils.Rmfile(db.Path)
		if err != nil {
			fmt.Println("failed to remove unencrypted database:", err)
		}
	}()

	rows, err := db.Conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passwords []*Password
	for rows.Next() {
		passwd := &Password{}
		err = rows.Scan(
			&passwd.Id,
			&passwd.Name,
			&passwd.Username,
			&passwd.Resource,
			&passwd.Password,
			&passwd.Comment,
			&passwd.Group,
		)

		if err != nil {
			return nil, err
		}

		passwords = append(passwords, passwd)
	}

	return passwords, nil
}
