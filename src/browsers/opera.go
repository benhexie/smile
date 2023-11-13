package browsers

func Opera() BrowserConfig {
	credentials := BrowserConfig{
		Browser:     "opera",
		Credentials: []Credential{},
	}

	return credentials
}
