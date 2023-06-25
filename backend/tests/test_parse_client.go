package main

import (
	"encoding/json"
	"fmt"
	"os"

	"email-wizard/backend/clients"
)

func test_parse() {
	var email map[string]interface{}
	body, err := os.ReadFile("example_email.json")
	if err != nil {
		fmt.Println("fail to open file")
		return
	}
	// fmt.Printf("%v\n", string(body))
	json.Unmarshal(body, &email)

	timezone := "Asia/Shanghai"

	events, err := clients.ParseEmail(email, timezone, 5)
	if err == nil {
		fmt.Printf("%v events parsed for email\n", len(events))
		// fmt.Println(events)
		fmt.Println("parse: ParseEmail test passed.")
	} else {
		fmt.Println("parse: ParseEmail test failed.")
	}
}
