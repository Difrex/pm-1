package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"math/rand"
	"os/exec"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type password struct {
	id       int
	name     string
	resource string
	password string
	username string
	comment  string
	group    string
}

func generate(length int) string {
	bytes, _ := exec.Command("head", "-c4096", "/dev/urandom").Output()
	hash := fmt.Sprintf("%x", md5.Sum(bytes))
	b64 := []rune(base64.StdEncoding.EncodeToString([]byte(hash)))
	chars := []rune("!@()#$%^&")
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10; i++ {
		for _, c := range chars {
			b64[rand.Intn(len(b64))] = c
		}
	}

	outStr := make([]rune, length)
	for i := range outStr {
		outStr[i] = b64[rand.Intn(len(b64))]
	}

	return string(outStr)
}

func addPassword(pass *password) error {
	db, err := decrypt()
	if err != nil {
		return err
	}

	query := `
insert into passwords(name, resource, password, username, comment, 'group')
values (?, ?, ?, ?, ?, ?)`
	return db.doQuery(query, pass.name, pass.resource, pass.password,
		pass.username, pass.comment, pass.group)
}

func removePassword(id int) error {
	db, err := decrypt()
	if err != nil {
		return err
	}

	return db.doQuery("delete from passwords where id=?", id)
}

func selectByName(name string) ([]*password, error) {
	db, err := decrypt()
	if err != nil {
		return nil, err
	}

	query := "select id, name, username, resource, password" +
		", comment, `group` from passwords where name=?"
	if name == "all" {
		query = "select id, name, username, resource, password" +
			", comment, `group` from passwords"
	}

	return db.doSelect(query, name)
}

func selectByGroup(name string) ([]*password, error) {
	db, err := decrypt()
	if err != nil {
		return nil, err
	}

	query := "select id, name, username, resource, password" +
		", comment, `group` from passwords where `group`=?"
	return db.doSelect(query, name)
}
