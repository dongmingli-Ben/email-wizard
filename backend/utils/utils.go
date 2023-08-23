package utils

import (
	"crypto/sha256"
	"email-wizard/backend/clients"
	"encoding/hex"
	"fmt"
)

var N_EMAIL_RETREIVAL int32 = 15

func ParsedUserEmailIDs(user_id int) ([]map[string]interface{}, error) {
	res, err := clients.Query([]string{"email_id", "email_address"},
		map[string]interface{}{"user_id": user_id}, "emails")
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetUserUnparsedEmails(emails []map[string]interface{}, email_address string,
	user_id int) ([]map[string]interface{}, error) {
	parsed_email_info, err := ParsedUserEmailIDs(user_id)
	if err != nil {
		return nil, err
	}
	email_ids := make(map[string]bool)
	for _, info := range parsed_email_info {
		if info["email_address"] != email_address {
			continue
		}
		email_ids[info["email_id"].(string)] = true
	}
	unparsed_emails := make([]map[string]interface{}, 0)
	for _, email := range emails {
		if _, ok := email_ids[email["email_id"].(string)]; !ok {
			unparsed_emails = append(unparsed_emails, email)
		}
	}
	return unparsed_emails, nil
}

func StoreUserEmails(emails []map[string]interface{}, account map[string]interface{}, user_id int) error {
	email_address := account["username"]
	for _, email := range emails {
		_, err := clients.AddRow(map[string]interface{}{
			"user_id":          user_id,
			"email_id":         email["email_id"],
			"email_address":    email_address,
			"mailbox_type":     account["protocol"],
			"email_subject":    email["subject"],
			"email_sender":     email["sender"],
			"email_recipients": email["recipient"],
			"email_date":       email["date"],
			"email_content":    email["content"],
			"event_ids":        []int32{},
		}, "emails")
		if err != nil {
			return err
		}
	}
	return nil
}

func GetUserEmailsFromAccount(account map[string]interface{}) ([]map[string]interface{}, error) {
	all_emails := make([]map[string]interface{}, 0)

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
	}
	return all_emails, nil
}

func GetUserProfile(user_id int) (map[string]interface{}, error) {
	res, err := clients.Query([]string{"user_name", "mailboxes"},
		map[string]interface{}{"user_id": user_id}, "users")
	if err != nil {
		return nil, err
	}
	if len(res) != 1 {
		return nil, fmt.Errorf("cannot find user with id: %v", user_id)
	}
	return res[0], nil
}

func GetUserEmailAccounts(user_id int) ([]map[string]interface{}, error) {
	user_profile, err := GetUserProfile(user_id)
	if err != nil {
		return nil, err
	}
	raw_accounts, ok := user_profile["mailboxes"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("fail to convert mailboxes from %v", user_profile)
	}
	accounts := make([]map[string]interface{}, 0)
	for _, account := range raw_accounts {
		acc, ok := account.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("fail to convert mailboxes from %v", account)
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func GetUserEmailAccountFromAddress(user_id int, email_address string) (map[string]interface{}, error) {
	accounts, err := GetUserEmailAccounts(user_id)
	if err != nil {
		return nil, err
	}
	for _, account := range accounts {
		if account["username"].(string) == email_address {
			return account, nil
		}
	}
	return nil, fmt.Errorf("cannot find mailbox with address %v for user with ID %v", 
		email_address, user_id)
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
		pk_values, err := clients.AddRow(map[string]interface{}{
			"user_id":       user_id,
			"email_id":      email_id,
			"email_address": email_address,
			"event_content": event,
		}, "events")
		if err != nil {
			return err
		}
		event_ids[i] = int(pk_values["event_id"].(float64))
	}
	// update email
	condition := map[string]interface{}{
		"email_id": email_id, "email_address": email_address,
	}
	data, err := clients.Query([]string{"event_ids"}, condition, "emails")
	if err != nil {
		return err
	}
	old_event_ids := data[0]["event_ids"].([]interface{})
	new_ids := make([]int, 0)
	for _, id := range old_event_ids {
		new_ids = append(new_ids, id.(int))
	}
	new_ids = append(new_ids, event_ids...)
	err = clients.UpdateValue("event_ids", new_ids, condition, "emails")
	return err
}

func ValidateUserSecret(user_id int, secret string) (bool, error) {
	res, err := clients.Query([]string{"user_id"},
		map[string]interface{}{
			"user_id":     user_id,
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
		map[string]interface{}{
			"user_name":     username,
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
		map[string]interface{}{"user_name": username}, "users")
	if err != nil {
		return 0, "", err
	}
	if len(res) != 1 {
		return 0, "", fmt.Errorf("fail to read user info for %v, not in DB", username)
	}
	return int(res[0]["user_id"].(float64)), res[0]["user_secret"].(string), nil
}

func GetUserEvents(user_id int) ([]map[string]interface{}, error) {
	events := make([]map[string]interface{}, 0)
	res, err := clients.Query([]string{"email_id", "email_address", "event_content"},
		map[string]interface{}{"user_id": user_id}, "events")
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
	_, err := clients.AddRow(map[string]interface{}{
		"user_secret":   secret,
		"user_name":     username,
		"user_password": password,
	}, "users")
	return err
}
