package data_tests

import (
	"fmt"
	"reflect"

	"email-wizard/data/utils"
	"testing"

	_ "github.com/lib/pq"
)

func prepare_db() error {
	if err := reset(); err != nil {
		return err
	}
	err := utils.AddRow(map[string]interface{}{
		"user_id":       1234323,
		"user_secret":   "oe2o950jgrnwgr",
		"user_name":     "jake",
		"user_password": "sjgn",
		"mailboxes":     []map[string]interface{}{},
	}, "users")
	if err != nil {
		return err
	}
	return nil
}

func TestQueryEmails(t *testing.T) {
	if err := prepare_db(); err != nil {
		t.Error(err.Error())
	}
	results, err := utils.Query([]string{"user_id", "email_address", "event_ids"}, map[string]interface{}{"user_id": 1234323},
		"emails")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(results)
}

func TestAddQueryEmails(t *testing.T) {
	if err := prepare_db(); err != nil {
		t.Error(err.Error())
	}
	err := utils.AddRow(map[string]interface{}{
		"user_id":       1234323,
		"email_id":      "oe2o950jgrnwgr",
		"email_address": "jake@example.com",
		"mailbox_type":  "outlook",
		"email_content": "example content",
		"event_ids":     []int32{0, 1, 2},
	}, "emails")
	if err != nil {
		t.Error(err.Error())
	}
	results, err := utils.Query([]string{"user_id", "email_address", "event_ids"},
		map[string]interface{}{"user_id": 1234323}, "emails")
	if err != nil {
		t.Error(err.Error())
	}
	if len(results) != 1 {
		t.Error("mismatched length")
	}
	if results[0]["user_id"].(int64) != 1234323 ||
		results[0]["email_address"].(string) != "jake@example.com" ||
		!reflect.DeepEqual(results[0]["event_ids"].([]interface{}), []interface{}{0., 1., 2.}) {
		fmt.Println(results[0]["user_id"].(int64), results[0]["email_address"].(string), results[0]["event_ids"].([]interface{}))
		fmt.Println(results[0]["user_id"].(int64) != 1234323)
		fmt.Println(results[0]["email_address"].(string) != "jake@example.com")
		fmt.Println(!reflect.DeepEqual(results[0]["event_ids"].([]interface{}), []interface{}{0., 1., 2.}))
		t.Error("mismatched content")
	}
}

func TestAddUpdateQueryEmails(t *testing.T) {
	if err := prepare_db(); err != nil {
		t.Error(err.Error())
	}
	err := utils.AddRow(map[string]interface{}{
		"user_id":       1234323,
		"email_id":      "oe2o950jgrnwgr",
		"email_address": "jake@example.com",
		"mailbox_type":  "outlook",
		"email_content": "example content",
		"event_ids":     []int32{0, 1, 2},
	}, "emails")
	if err != nil {
		t.Error(err.Error())
	}
	err = utils.UpdateValue("event_ids",
		[]int32{4, 5, 6},
		map[string]interface{}{"user_id": 1234323},
		"emails")
	if err != nil {
		t.Error(err.Error())
	}
	results, err := utils.Query([]string{"event_ids"},
		map[string]interface{}{"user_id": 1234323}, "emails")
	if err != nil {
		t.Error(err.Error())
	}
	if len(results) != 1 {
		t.Error("mismatched length")
	}
	if !reflect.DeepEqual(results[0]["event_ids"].([]interface{}), []interface{}{4., 5., 6.}) {
		t.Error("mismatched content")
	}
}

func TestAddDeleteQueryEmails(t *testing.T) {
	if err := prepare_db(); err != nil {
		t.Error(err.Error())
	}
	err := utils.AddRow(map[string]interface{}{
		"user_id":       1234323,
		"email_id":      "oe2o950jgrnwgr",
		"email_address": "jake@example.com",
		"mailbox_type":  "outlook",
		"email_content": "example content",
		"event_ids":     []int32{0, 1, 2},
	}, "emails")
	if err != nil {
		t.Error(err.Error())
	}
	err = utils.DeleteRows(map[string]interface{}{"user_id": 1234323},
		"emails")
	if err != nil {
		t.Error(err.Error())
	}
	results, err := utils.Query([]string{"event_ids"}, map[string]interface{}{},
		"emails")
	if err != nil {
		t.Error(err.Error())
	}
	if len(results) != 0 {
		t.Error("mismatched length")
	}
}
