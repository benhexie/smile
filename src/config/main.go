package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

var (
	SERVER_URL        = "https://rimiru.vercel.app"
	USER_ID           string
	SILENT            = true
	WRITE_FILE        = "offline"
	ONLINE            = true
	FEATURE_OPEN_FILE string
)

func SetConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("prop") // Explicitly set config type
	viper.AddConfigPath("..") // Check parent directory (useful for dev)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Config file not found...")
	}

	re := regexp.MustCompile("(?i)show")
	if re.MatchString(viper.GetString("MODE")) {
		SILENT = false
	}

	USER_ID = viper.GetString("USER_ID")

	re = regexp.MustCompile("(?i)forever|never")
	if re.MatchString(viper.GetString("WRITE_FILE")) {
		WRITE_FILE = strings.ToLower(viper.GetString("WRITE_FILE"))
	}

	re = regexp.MustCompile("(?i)false")
	if re.MatchString(viper.GetString("ONLINE")) {
		ONLINE = false
	}

	FEATURE_OPEN_FILE = viper.GetString("FEATURE_OPEN_FILE")
}
