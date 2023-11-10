package browsers

func Firefox() BrowserConfig {
	credentials := BrowserConfig{
		Browser: "firefox",
		Credentials: []Credential{},
	};

	return credentials;
}