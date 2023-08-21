package utils

import "fmt"

func AddUserMailboxOutlook(user_id int, mailbox_address string) error {
	return nil
}

func AddUserMailboxIMAP(user_id int, mailbox_address string, password string, imap_server string) error {
	fmt.Printf("%v saved to DB for user %v", mailbox_address, user_id)
	return nil
}

func AddUserMailboxPOP3(user_id int, mailbox_address string, password string, pop3_server string) error {
	fmt.Printf("%v saved to DB for user %v", mailbox_address, user_id)
	return nil
}
