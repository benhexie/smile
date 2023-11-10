package browsers

type Credential struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type BrowserConfig struct {
	Browser     string       `json:"browser"`
	Credentials []Credential `json:"credentials"`
}