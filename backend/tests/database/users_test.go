package data_tests

import (
	"fmt"

	utils "email-wizard/backend/clients"
	"testing"

	_ "github.com/lib/pq"
)

func TestQueryUsers(t *testing.T) {
	if err := utils.Reset(); err != nil {
		t.Error(err.Error())
	}
	results, err := utils.Query([]string{"user_name", "user_secret"}, map[string]interface{}{"user_id": 1234323},
		"users")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(results)
}

func TestAddQueryUsers(t *testing.T) {
	if err := utils.Reset(); err != nil {
		t.Error(err.Error())
	}
	pk_values, err := utils.AddRow(map[string]interface{}{
		// "user_id":       1234323,
		"user_secret":   "oe2o950jgrnwgr",
		"user_name":     "jake",
		"user_password": "sjgn",
		// "mailboxes":     []map[string]interface{}{},
	}, "users")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(pk_values)
	results, err := utils.Query([]string{"user_id", "user_secret", "mailboxes"}, pk_values,
		"users")
	if err != nil {
		t.Error(err.Error())
	}
	if len(results) != 1 {
		t.Error("mismatched length")
	}
	if results[0]["user_id"].(float64) != 1 || results[0]["user_secret"].(string) != "oe2o950jgrnwgr" {
		t.Error("mismatched content")
	}
}

func TestAddUpdateQueryUsers(t *testing.T) {
	if err := utils.Reset(); err != nil {
		t.Error(err.Error())
	}
	mailbox1 := map[string]interface{}{
		"address": "asdjr",
	}
	pk_values, err := utils.AddRow(map[string]interface{}{
		// "user_id":       1234323,
		"user_secret":   "oe2o950jgrnwgr",
		"user_name":     "jake",
		"user_password": "sjgn",
		"mailboxes":     []map[string]interface{}{},
	}, "users")
	if err != nil {
		t.Error(err.Error())
	}
	err = utils.UpdateValue("mailboxes",
		[]map[string]interface{}{mailbox1},
		pk_values,
		"users")
	if err != nil {
		t.Error(err.Error())
	}
	results, err := utils.Query([]string{"mailboxes"}, map[string]interface{}{},
		"users")
	if err != nil {
		t.Error(err.Error())
	}
	if len(results) != 1 {
		t.Error("mismatched length")
	}
	if results[0]["mailboxes"].([]interface{})[0].(map[string]interface{})["address"] != "asdjr" {
		t.Error("mismatched content")
	}
}

func TestAddDeleteQueryUsers(t *testing.T) {
	if err := utils.Reset(); err != nil {
		t.Error(err.Error())
	}
	pk_values, err := utils.AddRow(map[string]interface{}{
		// "user_id":       1234323,
		"user_secret":   "oe2o950jgrnwgr",
		"user_name":     "jake",
		"user_password": "sjgn",
		"mailboxes":     []map[string]interface{}{},
	}, "users")
	if err != nil {
		t.Error(err.Error())
	}
	err = utils.DeleteRows(pk_values, "users")
	if err != nil {
		t.Error(err.Error())
	}
	results, err := utils.Query([]string{"mailboxes"}, map[string]interface{}{},
		"users")
	if err != nil {
		t.Error(err.Error())
	}
	if len(results) != 0 {
		t.Error("mismatched length")
	}
}
