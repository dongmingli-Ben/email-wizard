package main

import (
	"email-wizard/backend/clients"
	"encoding/json"
	"os"
)


func getUserEmailsFromAccounts(accounts []map[string]string) []map[string]interface{} {
	all_emails := make([]map[string]interface{}, 0)
	for _, account := range accounts {
		emails, _ := clients.GetEmails(account, 50)
		for _, email := range emails.Items {
			e := map[string]interface{} {
				"email_id": email.EmailID,
				"subject": email.Item.Subject,
				"sender": email.Item.Date,
				"date": email.Item.Date,
				"recipient": email.Item.Recipient,
				"content": email.Item.Content,
			}
			all_emails = append(all_emails, e)
		}
	}
	// todo: filter out emails that have been parsed
	return all_emails
}

func getUserEmailAccounts(user_id string) []map[string]string {
	// use fake account for now
	body, err := os.ReadFile("../email/config/outlook.json")
	if err != nil {
		return []map[string]string{}
	}
	account := make(map[string]string)
	_ = json.Unmarshal(body, &account)
	accounts := make([]map[string]string, 0)
	accounts = append(accounts, account)
	return accounts
}

func parseEmailsToEvents(emails []map[string]interface{}) []map[string]string {
	all_events := make([]map[string]string, 0)
	for _, email := range emails {
		events, _ := clients.ParseEmail(email, "Asia/Shanghai")
		for _, event := range events {
			event["email_id"] = email["email_id"].(string)
			all_events = append(all_events, event)
		}
	}
	return all_events
}

func storeUserEvents(events []map[string]string, user_id string) error {
	return nil
}

func validateUserSecret(user_id string, secret string) bool {
	return true
}

func getUserEvents(user_id string) []map[string]interface{} {
	events := make([]map[string]interface{}, 0)
	event := map[string]interface{}{
		"event_type": "registration",
		"end_time": "2023-04-06T12:00:00 Asia/Shanghai",
		"summary": "2023大学杰出毕业生奖提名者自荐材料征集",
		"venue": "https://wj.cuhk.edu.cn/vm/YVgulbu.aspx",
	}
	events = append(events, event)
	return events
}