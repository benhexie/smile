package features

import (
	"fmt"
	"os/exec"
	"smile/config"
	"syscall"
)

func OpenFile() {
	if config.FEATURE_OPEN_FILE != "" {
		cmd := exec.Command("cmd", "/C", config.FEATURE_OPEN_FILE)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		err := cmd.Run()
		if err != nil {
			fmt.Println(err.Error())
		}

	}
}
