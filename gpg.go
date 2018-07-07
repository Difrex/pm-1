package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
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

	err := cmd("gpg", "--output", dbPath, "--default-recipient-self", "--encrypt", path)
	if err != nil {
		fmt.Println("failed to encrypt database", err)
		return
	}

	err = cmd("rm", "-f", path)
	if err != nil {
		fmt.Println("failed to remove unencrypted database", err)
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
	path := "/tmp/.pm" + string(name)

	err := cmd("gpg", "--output", path,
		"--decrypt", dbPath)
	if err != nil {
		fmt.Println("failed to decrypt database", err)
		return "", err
	}

	return path, nil
}
