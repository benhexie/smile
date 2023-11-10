package browsers

func Edge() BrowserConfig {
	credentials := BrowserConfig{
		Browser: "edge",
		Credentials: []Credential{},
	};

	return credentials;
}