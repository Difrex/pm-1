package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func PathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func Mkdir(path string) error {
	return os.MkdirAll(path, 0755)
}

func Mkfile(name string) error {
	_, err := os.Create(name)
	return err
}

func Rmfile(path string) error {
	return os.Remove(path)
}

func GetPrefix() string {
	if runtime.GOOS == "darvin" {
		return "/tmp/"
	}
	return "/dev/shm/"
}

func OpenURL(url string) {
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

func ReadLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(line), nil
}
