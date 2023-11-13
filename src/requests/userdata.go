package requests

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendUserData(data map[string]interface{}) (string, error) {
	url := "http://localhost/userdata";
	
	// convert data to json string
	dataJson, err := json.Marshal(data);
	if err != nil {
		return "", err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataJson))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close();

	return resp.Status, nil
}