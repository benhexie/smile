package browsers

func Safari() BrowserConfig {
	credentials := BrowserConfig{
		Browser:     "safari",
		Credentials: []Credential{},
	}

	return credentials
}
