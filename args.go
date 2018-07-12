package main

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/himidori/pm/db"
	"github.com/himidori/pm/utils"
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
	menu        bool
	rofi        bool
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
-m                      show dmenu
-R                      show rofi
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
	pflag.BoolVarP(&menu, "menu", "m", false, "show dmenu")
	pflag.BoolVarP(&rofi, "rofi", "R", false, "show rofi")
	pflag.Usage = printUsage

	pflag.Parse()
}

func parseArgs() {
	if menu {
		ok, err := utils.IsIntalled("dmenu")
		if err != nil {
			fmt.Println("failed to check dmenu installation:", err)
			return
		}
		if !ok {
			fmt.Println("dmenu is not installed")
			return
		}

		passwords, err := db.SelectAll()
		if err != nil {
			fmt.Println("failed to fetch passwords:", err)
			return
		}
		if passwords == nil {
			fmt.Println("no passwords found")
		}

		str := ""
		for _, p := range passwords {
			str += p.Name + "|" + p.Group + "|" + p.Resource + "\n"
		}
		res, err := utils.ShowMenu("dmenu", str)
		if err != nil {
			fmt.Println("failed to show menu:", err)
			return
		}
		if res == "" {
			return
		}

		res = strings.Split(res, "|")[0]
		for _, p := range passwords {
			if p.Name == res {
				err = clipboard.WriteAll(p.Password)
				if err != nil {
					utils.Notify(p.Name, "failed to copy password to the clipboard")
				}
				utils.Notify(p.Name, "copied password to the clipboard!")

				return
			}
		}

	}

	if rofi {
		ok, err := utils.IsIntalled("rofi")
		if err != nil {
			fmt.Println("failed to check rofi installation:", err)
			return
		}
		if !ok {
			fmt.Println("rofi is not installed")
			return
		}

		passwords, err := db.SelectAll()
		if err != nil {
			fmt.Println("failed to fetch passwords:", err)
			return
		}

		if passwords == nil {
			fmt.Println("no passwords found")
		}

		str := ""
		for _, p := range passwords {
			str += p.Name + "|" + p.Group + "|" + p.Resource + "\n"
		}
		res, err := utils.ShowMenu("rofi -dmenu", str)
		if err != nil {
			fmt.Println("failed to show menu:", err)
			return
		}
		if res == "" {
			return
		}

		res = strings.Split(res, "|")[0]
		for _, p := range passwords {
			if p.Name == res {
				err = clipboard.WriteAll(p.Password)
				if err != nil {
					utils.Notify(p.Name, "failed to copy password to the clipboard")
				}
				utils.Notify(p.Name, "copied password to the clipboard!")

				return
			}
		}
	}

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
			passwd, err := db.SelectByName(name)
			if err != nil {
				fmt.Println("failed to get password:", err)
				return
			}
			if passwd == nil {
				fmt.Println("no password found for name", name)
				return
			}

			if len(passwd) > 1 {
				db.PrintPaswords(passwd)
				return
			}

			err = clipboard.WriteAll(passwd[0].Password)
			if err != nil {
				fmt.Println("failed to copy password to the clipboard")
			} else {
				fmt.Println("password was copied to the clipboard!")
			}

			fmt.Print("URL: ")
			color.Blue(passwd[0].Resource)
			fmt.Print("User: ")
			color.Yellow(passwd[0].Username)
			if passwd[0].Group != "" {
				fmt.Print("Group: ")
				color.Magenta(passwd[0].Group)
			}

			if open {
				utils.OpenURL(passwd[0].Resource)
			}
		}

		if name == "" && group != "" {
			passwords, err := db.SelectByGroup(group)
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
			db.PrintPaswords(passwords)
		}

		if name != "" && group != "" {
			passwords, err := db.SelectByGroupAndName(name, group)
			if err != nil {
				fmt.Println("failed to get passwords:", err)
				return
			}

			if passwords == nil {
				fmt.Println("no password found")
				return
			}

			err = clipboard.WriteAll(passwords[0].Password)
			if err != nil {
				fmt.Println("failed to copy password to the clipboard")
			} else {
				fmt.Println("password was copied to the clipboard!")
			}

		}
	}

	if remove {
		if id == -1 {
			printUsage()
			return
		}

		err := db.RemovePassword(id)
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
			passwd = db.GeneratePassword(16)
		}

		err := db.AddPassword(&db.Password{
			Name:     name,
			Resource: link,
			Password: passwd,
			Username: user,
			Comment:  comment,
			Group:    group,
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
	name, err := utils.ReadLine()
	if err != nil {
		fmt.Println("failed to read line:", err)
		return
	}

	fmt.Print("resource: ")
	resource, err := utils.ReadLine()
	if err != nil {
		fmt.Println("failed to read line:", err)
		return
	}

	fmt.Print("password (leave empty to generate): ")
	passwd, err := utils.ReadLine()
	if err != nil {
		fmt.Println("failed to read line:", err)
		return
	}

	fmt.Print("username: ")
	username, err := utils.ReadLine()
	if err != nil {
		fmt.Println("failed to read line:", err)
		return
	}

	fmt.Print("comment: ")
	comment, err := utils.ReadLine()
	if err != nil {
		fmt.Println("failed to read line:", err)
		return
	}

	fmt.Print("group: ")
	grp, err := utils.ReadLine()
	if err != nil {
		fmt.Println("failed to read line:", err)
	}

	if passwd == "" {
		passwd = db.GeneratePassword(16)
	}

	err = db.AddPassword(&db.Password{
		Name:     name,
		Resource: resource,
		Password: passwd,
		Username: username,
		Comment:  comment,
		Group:    grp,
	})

	if err != nil {
		fmt.Println("failed to add password:", err)
		return
	}

	fmt.Println("successfuly added password to the database!")
}
