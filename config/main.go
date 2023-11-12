package config

import (
	"fmt"
	"regexp"

	"github.com/spf13/viper"
)

var (
	SERVER_URL   = "http://localhost"
	USER_ID      string
	SILENT       = true
)

func SetConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Config file not found...")
	}

	if viper.GetString("SERVER_URL") != "" {
		SERVER_URL = viper.GetString("SERVER_URL")
	}

	re := regexp.MustCompile("(?i)show")
	if re.MatchString(viper.GetString("MODE")) {
		SILENT = false;
	}

	USER_ID = viper.GetString("USER_ID")
}
