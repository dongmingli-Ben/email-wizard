package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"email-wizard/backend/clients"
)

func TestEmail(t *testing.T) {
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
	if err != nil {
		t.Error(err.Error())
	}
	if len(emails.Items) != 5 {
		t.Error(fmt.Sprintf("%v/%v emails retrieved", len(emails.Items), n_mails))
	}
	fmt.Printf("%v emails retrieved\n", len(emails.Items))
	fmt.Println("email GetEmail test passed.")
}