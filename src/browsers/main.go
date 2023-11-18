package browsers

import "os"

type Credential struct {
	Browser  string `json:"browser"`
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	APPDATA      string = os.Getenv("APPDATA")
	PROGRAMFILES string = os.Getenv("PROGRAMFILES")
	LOCALAPPDATA string = os.Getenv("LOCALAPPDATA")
)
