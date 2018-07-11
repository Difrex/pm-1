package db

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"math/rand"
	"os/exec"
	"time"
)

// Password struct
// holds every database field
type Password struct {
	Id       int
	Name     string
	Resource string
	Password string
	Username string
	Comment  string
	Group    string
}

// method used for inserting a new password
// into the database
func AddPassword(pass *Password) error {
	db, err := decrypt()
	if err != nil {
		return err
	}

	query := `
insert into passwords(name, resource, password, username, comment, 'group')
values (?, ?, ?, ?, ?, ?)`
	return db.doQuery(query, pass.Name, pass.Resource, pass.Password,
		pass.Username, pass.Comment, pass.Group)
}

// method used for removing a password
// from the database
func RemovePassword(id int) error {
	db, err := decrypt()
	if err != nil {
		return err
	}

	return db.doQuery("delete from passwords where id=?", id)
}

// method used for selecting passwords
// when the -n flag is provided
func SelectByName(name string) ([]*Password, error) {
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

// method used for selecting passwords
// when the -g flag is provided
func SelectByGroup(name string) ([]*Password, error) {
	db, err := decrypt()
	if err != nil {
		return nil, err
	}

	query := "select id, name, username, resource, password" +
		", comment, `group` from passwords where `group`=?"
	return db.doSelect(query, name)
}

// method used for generating a password
// of given length
func GeneratePassword(length int) string {
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

// method used for printing a passwords array
func PrintPaswords(passwords []*Password) {
	for _, p := range passwords {
		fmt.Printf("%d | %s | %s | %s | %s\n",
			p.Id, p.Name, p.Resource, p.Username, p.Comment)
	}
}
