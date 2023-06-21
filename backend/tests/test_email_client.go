package main

import (
	"encoding/json"
	"fmt"
	"os"

	"email-wizard/backend/clients"
)

func test_email() {
	var user_config map[string]string
	body, err := os.ReadFile("outlook.json")
	if err != nil {
		fmt.Println("fail to open file")
		return
	}
	// fmt.Printf("%v\n", string(body))
	json.Unmarshal(body, &user_config)

	var n_mails int32 = 5

	emails, err := clients.GetEmails(user_config, n_mails)
	if err == nil {
		fmt.Printf("%v emails retrieved\n", len(emails.Items))
		fmt.Println("email GetEmail test passed.")
	} else {
		fmt.Println("email: GetEmail test failed.")
	}
}