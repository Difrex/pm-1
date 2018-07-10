package db

import (
	"os"

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
	path := utils.GetPrefix() + ".pm" + utils.GenerateName()

	err := gpg.DecryptFile(dbPath, path)
	if err != nil {
		return nil, err
	}

	return newDB(path)
}
