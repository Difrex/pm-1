package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"runtime"

	"github.com/Difrex/gpg"
)

func encrypt(path string) {
	dbPath := os.Getenv("HOME") + "/.PM/db.sqlite"
	if pathExists(dbPath) {
		err := cmd("rm", "-f", dbPath)
		if err != nil {
			fmt.Println("failed to remove old database", err)
			return
		}
	}

	err := gpg.EncryptFileRecipientSelf(path, dbPath)
	if err != nil {
		fmt.Println("failed to encrypt database", err.Error())
		return
	}

	err = cmd("rm", "-f", path)
	if err != nil {
		fmt.Println("failed to remove unencrypted database", err.Error())
		return
	}
}

func decrypt() (string, error) {
	dbPath := os.Getenv("HOME") + "/.PM/db.sqlite"
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	name := make([]byte, 10)
	rand.Seed(time.Now().UnixNano())
	for i := range name {
		name[i] = chars[rand.Intn(len(chars))]
	}

	// Detect OS and swith to path output prefix
	pathPrefix := "/dev/shm/"
	if runtime.GOOS != "linux" {
		pathPrefix = "/tmp/"
	}

	path := pathPrefix + ".pm" + string(name)

	err := gpg.DecryptFile(dbPath, path)
	if err != nil {
		fmt.Println("failed to decrypt database", err)
		return "", err
	}

	return path, nil
}
