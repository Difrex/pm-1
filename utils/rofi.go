package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func IsRofiInstalled() (bool, error) {
	out, err := exec.Command("whereis", "rofi").Output()
	if err != nil {
		return false, err
	}

	data := strings.Split(string(out), ":")
	return len(data[1]) > 1, nil
}

func RofiShow(passwords string) string {
	out, err := exec.Command("bash", "-c", "echo '"+passwords+"' | rofi -dmenu").Output()
	if err != nil {
		fmt.Println("failed to spawn dmenu:", err)
		return ""
	}

	return strings.TrimSpace(string(out))
}
