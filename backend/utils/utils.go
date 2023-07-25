package utils

import (
	"email-wizard/backend/clients"
	"encoding/json"
	"errors"
	"os"
)

var N_EMAIL_RETREIVAL int32 = 5

func ParsedUserEmailIDs(user_id string) []string {
	return make([]string, 0)
}

func GetUserUnparsedEmails(emails []map[string]interface{}, user_id string) []map[string]interface{} {
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

func GetUserEmailsFromAccounts(accounts []map[string]string) ([]map[string]interface{}, error) {
	all_emails := make([]map[string]interface{}, 0)
	for _, account := range accounts {
		emails, err := clients.GetEmails(account, N_EMAIL_RETREIVAL)
		if err != nil {
			return make([]map[string]interface{}, 0), err
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
	}
	return all_emails, nil
}

func GetUserProfile(user_id string) (map[string]interface{}, error) {
	// use fake account for now
	user_name := "Jake"
	body, err := os.ReadFile("tests/outlook.json")
	if err != nil {
		return make(map[string]interface{}), err
	}
	account := make(map[string]string)
	_ = json.Unmarshal(body, &account)
	accounts := make([]map[string]string, 0)
	accounts = append(accounts, account)
	return map[string]interface{}{
		"user_name":     user_name,
		"user_accounts": accounts,
	}, nil
}

func GetUserEmailAccounts(user_id string) ([]map[string]string, error) {
	user_profile, err := GetUserProfile(user_id)
	if err != nil {
		return make([]map[string]string, 0), nil
	}
	accounts, ok := user_profile["user_accounts"].([]map[string]string)
	if !ok {
		return make([]map[string]string, 0), errors.New("fail to get user_accounts from profile")
	}
	return accounts, nil
}

func ParseEmailsToEvents(emails []map[string]interface{}, retry int) []map[string]string {
	all_events := make([]map[string]string, 0)
	for _, email := range emails {
		events, _ := clients.ParseEmail(email, "Asia/Shanghai", retry)
		for _, event := range events {
			event["email_id"] = email["email_id"].(string)
			all_events = append(all_events, event)
		}
	}
	return all_events
}

func StoreUserEvents(events []map[string]string, user_id string) error {
	return nil
}

func ValidateUserSecret(user_id string, secret string) bool {
	return true
}

func GetUserEvents(user_id string) []map[string]interface{} {
	events := make([]map[string]interface{}, 0)
	event := map[string]interface{}{
		"event_type": "registration",
		"end_time":   "2023-04-06T12:00:00 Asia/Shanghai",
		"summary":    "2023大学杰出毕业生奖提名者自荐材料征集",
		"venue":      "https://wj.cuhk.edu.cn/vm/YVgulbu.aspx",
	}
	events = append(events, event)
	return events
}
