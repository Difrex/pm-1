package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func IsDmenuInstalled() (bool, error) {
	out, err := exec.Command("whereis", "dmenu").Output()
	if err != nil {
		return false, err
	}

	data := strings.Split(string(out), ":")
	return len(data[1]) > 1, nil
}

func DmenuShow(passwords string) string {
	out, err := exec.Command("bash", "-c", "echo '"+passwords+"' | dmenu").Output()
	if err != nil {
		fmt.Println("failed to spawn dmenu:", err)
		return ""
	}

	return strings.TrimSpace(string(out))
}
