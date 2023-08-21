package utils

import (
	"email-wizard/backend/clients"
)

func AddUserMailboxOutlook(user_id int, mailbox_address string) error {
	res, err := clients.Query([]string{"mailboxes"}, map[string]interface{}{
		"user_id": user_id,
	}, "users")
	if err != nil {
		return err
	}
	// raw_mailboxes := res[0]["mailboxes"].([]interface{})
	// mailboxes := make([]map[string]interface{}, 0)
	// for _, mbox := range raw_mailboxes {
	// 	mailboxes = append(mailboxes, mbox.(map[string]interface{}))
	// }
	mailboxes := res[0]["mailboxes"].([]map[string]interface{})
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
	var mailboxes []map[string]interface{}
	if _, ok := res[0]["mailboxes"].([]map[string]interface{}); !ok {
		mailboxes = make([]map[string]interface{}, 0)
	} else {
		mailboxes = res[0]["mailboxes"].([]map[string]interface{})
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
	mailboxes := res[0]["mailboxes"].([]map[string]interface{})
	mailboxes = append(mailboxes, map[string]interface{}{
		"username":    mailbox_address,
		"protocol":    "IMAP",
		"password":    password,
		"pop3_server": pop3_server,
	})
	err = clients.UpdateValue("mailboxes", mailboxes, map[string]interface{}{
		"user_id": user_id,
	}, "users")
	return err
}
