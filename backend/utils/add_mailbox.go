package utils

import (
	"email-wizard/backend/clients"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func prepare_mailboxes(mailboxes interface{}) ([]map[string]interface{}, error) {
	raw_mailboxes, ok := mailboxes.([]interface{})
	ret_mailboxes := make([]map[string]interface{}, 0)
	if !ok {
		return ret_mailboxes, nil
	}
	for _, mbox := range raw_mailboxes {
		ret_mailboxes = append(ret_mailboxes, mbox.(map[string]interface{}))
	}
	return ret_mailboxes, nil
}

/*
Do necessary preprocessing for user credentials, for example:
 * For Gmail, exchange auth code for access token and refresh token
 * For other mailbox, do nothing
 */
func prepare_credentials(mailbox_type string, mailbox_address string, 
		credentials map[string]interface{}) (map[string]interface{}, error) {
	if (mailbox_type == "gmail") {
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
		req_url := fmt.Sprintf("https://oauth2.googleapis.com/token?client_id=%v&client_secret=%v&code=%v&grant_type=authorization_code&redirect_uri=%v",
							   app_creds["client_id"], app_creds["client_secret"], 
							   credentials["auth_code"], app_creds["redirect_uris"].([]interface{})[0])
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
		// set expire timestamp
		creds["expire_timestamp"] = time.Now().Unix() + int64(creds["expires_in"].(float64))
		creds["client_id"] = app_creds["client_id"]
		creds["client_secret"] = app_creds["client_secret"]
		return creds, nil
	}
	return credentials, nil
}

func AddUserMailbox(user_id int, mailbox_type string, mailbox_address string, credentials map[string]interface{}) error {
	if _, err := GetUserEmailAccountFromAddress(user_id, mailbox_address); err == nil {
		return fmt.Errorf("mailbox %v already added", mailbox_address)
	}
	creds, err := prepare_credentials(mailbox_type, mailbox_address, credentials)
	if err != nil {
		return err
	}
	res, err := clients.Query([]string{"mailboxes"}, map[string]interface{}{
		"user_id": user_id,
	}, "users")
	if err != nil {
		return err
	}
	mailboxes, err := prepare_mailboxes(res[0]["mailboxes"])
	if err != nil {
		return err
	}
	mailboxes = append(mailboxes, map[string]interface{}{
		"username": mailbox_address,
		"protocol": mailbox_type,
		"credentials": creds,
	})
	err = clients.UpdateValue("mailboxes", mailboxes, map[string]interface{}{
		"user_id": user_id,
	}, "users")
	return err
}