package utils

import (
	"email-wizard/backend/clients"
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

func AddUserMailboxOutlook(user_id int, mailbox_address string) error {
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
		"protocol": "outlook",
	})
	err = clients.UpdateValue("mailboxes", mailboxes, map[string]interface{}{
		"user_id": user_id,
	}, "users")
	return err
}

func AddUserMailboxIMAP(user_id int, mailbox_address string, password string, imap_server string) error {
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
		"username":    mailbox_address,
		"protocol":    "IMAP",
		"password":    password,
		"imap_server": imap_server,
	})
	err = clients.UpdateValue("mailboxes", mailboxes, map[string]interface{}{
		"user_id": user_id,
	}, "users")
	return err
}

func AddUserMailboxPOP3(user_id int, mailbox_address string, password string, pop3_server string) error {
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
		"username":    mailbox_address,
		"protocol":    "POP3",
		"password":    password,
		"pop3_server": pop3_server,
	})
	err = clients.UpdateValue("mailboxes", mailboxes, map[string]interface{}{
		"user_id": user_id,
	}, "users")
	return err
}
