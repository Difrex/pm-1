package main

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/ogier/pflag"
)

var (
	show        bool
	name        string
	group       string
	new         bool
	link        string
	user        string
	comment     string
	pass        string
	remove      bool
	id          int
	open        bool
	interactive bool
)

func printUsage() {
	fmt.Println(`Simple password manager written in Go

-s                      show password
-n [Name of resource]   name of resource
-g [Group name]         group name
-w                      store new password
-l [Link]               link to resource
-u                      username
-c                      comment
-p [Password]           password
                        (if password is omitted PM will
                        generate a secure password)
-r                      remove password
-i                      password ID
-o                      open link
-I                      interactive mode for adding new password
-h                      show help`)
}

func initArgs() {
	pflag.BoolVarP(&show, "show", "s", false, "show password")
	pflag.StringVarP(&name, "name", "n", "", "name of the resource")
	pflag.StringVarP(&group, "group", "g", "", "name of the group")
	pflag.BoolVarP(&new, "write", "w", false, "add new password")
	pflag.StringVarP(&link, "link", "l", "", "link to the resource")
	pflag.StringVarP(&user, "user", "u", "", "username of the resource")
	pflag.StringVarP(&comment, "comment", "c", "", "comment")
	pflag.StringVarP(&pass, "password", "p", "", "password")
	pflag.BoolVarP(&remove, "remove", "r", false, "remove password")
	pflag.IntVarP(&id, "id", "i", -1, "password id")
	pflag.BoolVarP(&open, "open", "o", false, "open link in browser")
	pflag.BoolVarP(&interactive, "interactive", "I", false, "interactive mode")
	pflag.Usage = printUsage

	pflag.Parse()
}

func parseArgs() {
	if !show && !new && !remove {
		printUsage()
		return
	}

	if show {
		if name == "" && group == "" {
			printUsage()
			return
		}

		if name != "" && group == "" {
			passwd, err := selectByName(name)
			if err != nil {
				fmt.Println("failed to get password:", err)
				return
			}
			if passwd == nil {
				fmt.Println("no password found for name", name)
				return
			}

			if len(passwd) > 1 {
				printPaswords(passwd)
				return
			}

			err = clipboard.WriteAll(passwd[0].password)
			if err != nil {
				fmt.Println("failed to copy password to the clipboard")
			} else {
				fmt.Println("password was copied to the clipboard!")
			}

			fmt.Print("URL: ")
			color.Blue(passwd[0].resource)
			fmt.Print("User: ")
			color.Yellow(passwd[0].username)
			if passwd[0].group != "" {
				fmt.Print("Group: ")
				color.Magenta(passwd[0].group)
			}

			if open {
				openURL(passwd[0].resource)
			}
		}

		if name == "" && group != "" {
			passwords, err := selectByGroup(group)
			if err != nil {
				fmt.Println("failed to get passwords:", err)
				return
			}

			if passwords == nil {
				fmt.Println("no passwords found for group", group)
				return
			}

			fmt.Print("Group: ")
			color.Magenta(group)
			printPaswords(passwords)
		}
	}

	if remove {
		if id == -1 {
			printUsage()
			return
		}

		err := removePassword(id)
		if err != nil {
			fmt.Println("failed to remove password:", err)
			return
		}

		fmt.Println("successfuly removed password with id", id)
	}

	if new {
		if interactive {
			addInteractive()
			return
		}

		if name == "" || link == "" {
			printUsage()
			return
		}

		passwd := ""

		if pass != "" {
			passwd = pass
		} else {
			passwd = generate(16)
		}

		err := addPassword(&password{
			name:     name,
			resource: link,
			password: passwd,
			username: user,
			comment:  comment,
			group:    group,
		})

		if err != nil {
			fmt.Println("failed to add password:", err)
			return
		}

		fmt.Println("successfuly added new password!")
	}
}

func addInteractive() {
	fmt.Print("name: ")
	name, err := readLine()
	if err != nil {
		fmt.Println("failed to read line:", err)
		return
	}

	fmt.Print("resource: ")
	resource, err := readLine()
	if err != nil {
		fmt.Println("failed to read line:", err)
		return
	}

	fmt.Print("password (leave empty to generate): ")
	passwd, err := readLine()
	if err != nil {
		fmt.Println("failed to read line:", err)
		return
	}

	fmt.Print("username: ")
	username, err := readLine()
	if err != nil {
		fmt.Println("failed to read line:", err)
		return
	}

	fmt.Print("comment: ")
	comment, err := readLine()
	if err != nil {
		fmt.Println("failed to read line:", err)
		return
	}

	fmt.Print("group: ")
	grp, err := readLine()
	if err != nil {
		fmt.Println("failed to read line:", err)
	}

	if passwd == "" {
		passwd = generate(16)
	}

	err = addPassword(&password{
		name:     name,
		resource: resource,
		password: passwd,
		username: username,
		comment:  comment,
		group:    grp,
	})

	if err != nil {
		fmt.Println("failed to add password:", err)
		return
	}

	fmt.Println("successfuly added password to the database!")
}
