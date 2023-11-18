package device

import (
	"os"
	"os/exec"
	"runtime"
)

func GetSysInfo() map[string]interface{} {
	sysInfo := make(map[string]interface{})

	sysInfo["label"] = "system"
	sysInfo["os"] = runtime.GOOS
	sysInfo["arch"] = runtime.GOARCH
	sysInfo["cpuCores"] = runtime.NumCPU()
	sysInfo["envVars"] = os.Environ()
	sysInfo["hostname"], _ = os.Hostname()
	sysInfo["username"] = os.Getenv("USERNAME")
	sysInfo["ipconfig"] = getIpconfig()

	return sysInfo
}

func getIpconfig() string {
	cmd := exec.Command("ipconfig")
	out, err := cmd.Output()
	if err != nil {
		return err.Error()
	}
	return string(out)
}
