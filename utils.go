package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func pathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func mkdir(path string) error {
	return os.MkdirAll(path, 0755)
}

func mkfile(name string) error {
	_, err := os.Create(name)
	return err
}

func rmfile(path string) error {
	return os.Remove(path)
}

func getPrefix() string {
	if runtime.GOOS == "darvin" {
		return "/tmp/"
	}
	return "/dev/shm/"
}

func cmd(name string, args ...string) error {
	_, err := exec.Command(name, args...).Output()
	return err
}

func openURL(url string) {
	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()

	case "darwin":
		err = exec.Command("open", url).Start()
	}

	if err != nil {
		fmt.Printf("failed to open url %s: %s\n", url, err)
	}
}

func readLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(line), nil
}

func printPaswords(passwords []*password) {
	for _, p := range passwords {
		fmt.Printf("%d | %s | %s | %s | %s\n",
			p.id, p.name, p.resource, p.username, p.comment)
	}
}
