package utils

import (
	"crypto/sha256"
	"email-wizard/backend/clients"
	"encoding/hex"
	"fmt"
)

var N_EMAIL_RETREIVAL int32 = 5

func ParsedUserEmailIDs(user_id int) []string {
	return make([]string, 0)
}

func GetUserUnparsedEmails(emails []map[string]interface{}, user_id int) []map[string]interface{} {
	parsed_email_ids := ParsedUserEmailIDs(user_id)
	email_ids := make(map[string]bool)
	for _, id := range parsed_email_ids {
		email_ids[id] = true
	}
	unparsed_emails := make([]map[string]interface{}, 0)
	for _, email := range emails {
		if _, ok := email_ids[email["email_id"].(string)]; !ok {
			unparsed_emails = append(unparsed_emails, email)
		}
	}
	return unparsed_emails
}

func StoreUserEmail(account map[string]string, email map[string]interface{}) error {
	return nil
}

func GetUserEmailsFromAccounts(accounts []map[string]string) ([]map[string]interface{}, error) {
	all_emails := make([]map[string]interface{}, 0)
	for _, account := range accounts {
		emails, err := clients.GetEmails(account, N_EMAIL_RETREIVAL)
		if err != nil {
			return nil, err
		}
		for _, email := range emails.Items {
			e := map[string]interface{}{
				"email_id":  email.EmailID,
				"subject":   email.Item.Subject,
				"sender":    email.Item.Sender,
				"date":      email.Item.Date,
				"recipient": email.Item.Recipient,
				"content":   email.Item.Content,
			}
			all_emails = append(all_emails, e)
			err := StoreUserEmail(account, e)
			if err != nil {
				return nil, err
			}
		}
	}
	return all_emails, nil
}

func GetUserProfile(user_id int) (map[string]interface{}, error) {
	res, err := clients.Query([]string {"user_name", "mailboxes"},
		map[string]interface{} {"user_id": user_id}, "users")
	if err != nil {
		return nil, err
	}
	if (len(res) != 1) {
		return nil, fmt.Errorf("cannot find user with id: %v", user_id)
	}
	return res[0], nil
}

func GetUserEmailAccounts(user_id int) ([]map[string]string, error) {
	user_profile, err := GetUserProfile(user_id)
	if err != nil {
		return nil, err
	}
	accounts, ok := user_profile["mailboxes"].([]map[string]string)
	if !ok {
		return nil, fmt.Errorf("fail to convert mailboxes from %v", user_profile)
	}
	return accounts, nil
}

func ParseEmailToEvents(email map[string]interface{}, retry int) ([]map[string]string, error) {
	events, err := clients.ParseEmail(email, "Asia/Shanghai", retry)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func StoreUserEvents(events []map[string]string, user_id int, email_id string, email_address string) error {
	event_ids := make([]int, len(events))
	for i, event := range events {
		pk_values, err := clients.AddRow(map[string]interface{} {
			"user_id": user_id,
			"email_id": email_id,
			"email_address": email_address,
			"email_content": event,
		}, "events")
		if err != nil {
			return err
		}
		event_ids[i] = int(pk_values["event_id"].(float64))
	}
	// update email
	condition := map[string]interface{} {
		"email_id": email_id, "email_address": email_address,
	}
	data, err := clients.Query([]string{"event_ids"}, condition, "emails")
	if err != nil {
		return err
	}
	old_event_ids := data[0]["event_ids"].([]int)
	event_ids = append(old_event_ids, event_ids...)
	err = clients.UpdateValue("event_ids", event_ids, condition, "email")
	return err
}

func ValidateUserSecret(user_id int, secret string) (bool, error) {
	res, err := clients.Query([]string{"user_id"}, 
		map[string]interface{} {
			"user_id": user_id,
			"user_secret": secret,
		}, "users")
	if err != nil {
		return false, err
	}
	if len(res) == 0 {
		return false, nil
	}
	return true, nil
}

func ValidateUserPassword(username string, password string) (bool, error) {
	res, err := clients.Query([]string{"user_id"}, 
		map[string]interface{} {
			"user_name": username,
			"user_password": password,
		}, "users")
	if err != nil {
		return false, err
	}
	if len(res) == 0 {
		return false, nil
	}
	return true, nil
}

func GetUserIdSecret(username string) (int, string, error) {
	res, err := clients.Query([]string{"user_id", "user_secret"},
		map[string]interface{} {"user_name": username}, "users")
	if err != nil {
		return 0, "", err
	}
	if len(res) != 1 {
		return 0, "", fmt.Errorf("fail to read user info for %v, not in DB", username)
	}
	return res[0]["user_id"].(int), res[0]["user_secret"].(string), nil
}

func GetUserEvents(user_id int) ([]map[string]interface{}, error) {
	events := make([]map[string]interface{}, 0)
	res, err := clients.Query([]string{"email_id", "email_address", "event_content"},
		map[string]interface{} {"user_id": user_id}, "events")
	if err != nil {
		return nil, err
	}
	for _, value := range res {
		event, ok := value["event_content"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("fail to load event content from %v", 
				value["event_content"])
		}
		events = append(events, event)
	}
	return events, nil
}

func AddUserDB(username string, password string) error {
	fmt.Printf("adding username: %v, password: %v\n", username, password)
	hash := sha256.New()
	hash.Write([]byte(username + "|" + password))
	hash_bytes := hash.Sum(nil)
	secret := hex.EncodeToString(hash_bytes)
	_, err := clients.AddRow(map[string]interface{} {
		"user_secret": secret,
		"user_name": username,
		"user_password": password,
	}, "users")
	return err
}