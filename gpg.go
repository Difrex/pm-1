package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/Difrex/gpg"
)

func encrypt(path string) error {
	dbPath := os.Getenv("HOME") + "/.PM/db.sqlite"
	if pathExists(dbPath) {
		err := rmfile(dbPath)
		if err != nil {
			return err
		}
	}

	err := gpg.EncryptFileRecipientSelf(path, dbPath)
	if err != nil {
		return err
	}

	return rmfile(path)
}

func decrypt() (*DB, error) {
	dbPath := os.Getenv("HOME") + "/.PM/db.sqlite"
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	name := make([]byte, 10)
	rand.Seed(time.Now().UnixNano())
	for i := range name {
		name[i] = chars[rand.Intn(len(chars))]
	}

	path := getPrefix() + ".pm" + string(name)

	err := gpg.DecryptFile(dbPath, path)
	if err != nil {
		return nil, err
	}

	db, err := NewDB(path)
	if err != nil {
		return nil, err
	}

	return db, nil
}
