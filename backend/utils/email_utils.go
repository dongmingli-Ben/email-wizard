package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)


func RefreshGmailToken(credentials map[string]interface{}) (map[string]interface{}, error) {
	// load app credentials
	data, err := os.ReadFile("cert/credentials.json")
	if err != nil {
		return nil, err
	}
	app_creds := make(map[string]interface{})
	if err = json.Unmarshal(data, &app_creds); err != nil {
		return nil, err
	}
	app_creds = app_creds["web"].(map[string]interface{})
	// prepare POST request for token exchange
	req_url := fmt.Sprintf("https://oauth2.googleapis.com/token?client_id=%v&client_secret=%v&refresh_token=%v&grant_type=refresh_token",
						   app_creds["client_id"], app_creds["client_secret"], 
						   credentials["refresh_token"])
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", req_url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Host", "oauth2.googleapis.com")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	creds := make(map[string]interface{})
	if err = json.Unmarshal(body, &creds); err != nil {
		return nil, err
	}
	fmt.Printf("Google response to refresh: %v\n", string(body))
	if _, ok := creds["error"]; ok {
		return nil, fmt.Errorf("%v: %v", creds["error"], creds["error_description"])
	}
	credentials["access_token"] = creds["access_token"]
	credentials["expires_in"] = creds["expires_in"]
	// set expire timestamp
	credentials["expire_timestamp"] = time.Now().Unix() + int64(creds["expires_in"].(float64))
	return credentials, nil
}

func RevokeGmailAccess(credentials map[string]interface{}) error {
	req_url := fmt.Sprintf("https://oauth2.googleapis.com/revoke?token=%v",
						   credentials["refresh_token"])
	req, err := http.NewRequest("POST", req_url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Google response to revoke: %v\n", string(body))
	data := make(map[string]interface{})
	if err = json.Unmarshal(body, &data); err != nil {
		return err
	}
	if resp.StatusCode != 200 && !(data["error"].(string) == "invalid_token" && 
								   data["error_description"].(string) == "Token expired or revoked") {
		return fmt.Errorf("encounter status code %v when revoking gmail access", resp.StatusCode)
	}
	return nil
}