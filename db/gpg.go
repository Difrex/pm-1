package db

import (
	"math/rand"
	"os"
	"time"

	"github.com/Difrex/gpg"
	"github.com/himidori/pm/utils"
)

func encrypt(path string) error {
	dbPath := os.Getenv("HOME") + "/.PM/db.sqlite"
	if utils.PathExists(dbPath) {
		err := utils.Rmfile(dbPath)
		if err != nil {
			return err
		}
	}

	err := gpg.EncryptFileRecipientSelf(path, dbPath)
	if err != nil {
		return err
	}

	return utils.Rmfile(path)
}

func decrypt() (*DB, error) {
	dbPath := os.Getenv("HOME") + "/.PM/db.sqlite"
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	name := make([]byte, 10)
	rand.Seed(time.Now().UnixNano())
	for i := range name {
		name[i] = chars[rand.Intn(len(chars))]
	}

	path := utils.GetPrefix() + ".pm" + string(name)

	err := gpg.DecryptFile(dbPath, path)
	if err != nil {
		return nil, err
	}

	return NewDB(path)
}
