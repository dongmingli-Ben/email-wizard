package utils

import "email-wizard/backend/clients"

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

func UpdateUserEventsForAccount(user_id int, account map[string]interface{}) error {
	// read recent emails from user email accounts (and store emails to DB)
	emails, err := GetUserEmailsFromAccount(account)
	if err != nil {
		return err
	}

	// filter for un-parsed emails
	emails, err = GetUserUnparsedEmails(emails, account["username"].(string), user_id)
	if err != nil {
		return err
	}
	for _, email := range emails {
		// todo: atomic!
		err = StoreUserEmails([]map[string]interface{}{email}, account, user_id)
		if err != nil {
			return err
		}
		// parse into events
		events, err := ParseEmailToEvents(email, 5)
		if err != nil {
			return err
		}
		// store back to db
		err = StoreUserEvents(events, user_id,
			email["email_id"].(string),
			account["username"].(string))
		if err != nil {
			return err
		}
	}

	return nil
}