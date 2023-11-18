package requests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"smile/config"
)

func SendReconData(data map[string]interface{}) (string, error) {
	url := config.SERVER_URL + "/smile"

	// convert data to json string
	dataJson, err := json.Marshal(data)
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
	defer resp.Body.Close()

	return resp.Status, nil
}
