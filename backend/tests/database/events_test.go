package data_tests

import (
	"fmt"
	"reflect"

	utils "email-wizard/backend/clients"
	"testing"

	_ "github.com/lib/pq"
)

func prepare_db_for_events() error {
	if err := utils.Reset(); err != nil {
		return err
	}
	pk_values, err := utils.AddRow(map[string]interface{}{
		// "user_id":       1234323,
		"user_secret":   "oe2o950jgrnwgr",
		"user_name":     "jake",
		"user_password": "sjgn",
		"mailboxes":     []map[string]interface{}{},
	}, "users")
	if err != nil {
		return err
	}
	_, err = utils.AddRow(map[string]interface{}{
		"user_id":       pk_values["user_id"],
		"email_id":      "oe2o950jgrnwgr",
		"email_address": "jake@example.com",
		"mailbox_type":  "outlook",
		"email_content": "example content",
		"event_ids":     []int32{0, 1, 2},
	}, "emails")
	if err != nil {
		return err
	}
	return nil
}

func TestQueryEvents(t *testing.T) {
	if err := prepare_db_for_events(); err != nil {
		t.Error(err.Error())
	}
	results, err := utils.Query([]string{"event_id", "email_id", "email_address", "event_content"},
		map[string]interface{}{"email_id": "oe2o950jgrnwgr"},
		"events")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(results)
}

func TestAddQueryEvents(t *testing.T) {
	if err := prepare_db_for_events(); err != nil {
		t.Error(err.Error())
	}
	content := map[string]interface{}{
		"type": "notification", "content": "test notification",
	}
	pk_values, err := utils.AddRow(map[string]interface{}{
		"email_id":      "oe2o950jgrnwgr",
		"email_address": "jake@example.com",
		"event_content": content,
	}, "events")
	if err != nil {
		t.Error(err.Error())
	}
	results, err := utils.Query([]string{"event_id", "email_id", "email_address", "event_content"},
		pk_values, "events")
	if err != nil {
		t.Error(err.Error())
	}
	if len(results) != 1 {
		t.Error("mismatched length")
	}
	if results[0]["email_id"].(string) != "oe2o950jgrnwgr" ||
		results[0]["email_address"].(string) != "jake@example.com" ||
		!reflect.DeepEqual(results[0]["event_content"].(map[string]interface{}), content) {
		// fmt.Println(results[0]["user_id"].(int64), results[0]["email_address"].(string), results[0]["event_ids"].([]interface{}))
		// fmt.Println(results[0]["user_id"].(int64) != 1234323)
		// fmt.Println(results[0]["email_address"].(string) != "jake@example.com")
		// fmt.Println(!reflect.DeepEqual(results[0]["event_ids"].([]interface{}), []interface{}{0., 1., 2.}))
		t.Error("mismatched content")
	}
}

func TestAddUpdateQueryEvents(t *testing.T) {
	if err := prepare_db_for_events(); err != nil {
		t.Error(err.Error())
	}
	content := map[string]interface{}{
		"type": "notification", "content": "test notification",
	}
	pk_values, err := utils.AddRow(map[string]interface{}{
		"email_id":      "oe2o950jgrnwgr",
		"email_address": "jake@example.com",
		"event_content": content,
	}, "events")
	if err != nil {
		t.Error(err.Error())
	}
	new_content := map[string]interface{}{
		"type": "registration", "content": "test notification",
	}
	err = utils.UpdateValue("event_content", new_content,
		pk_values, "events")
	if err != nil {
		t.Error(err.Error())
	}
	results, err := utils.Query([]string{"event_content"},
		pk_values, "events")
	if err != nil {
		t.Error(err.Error())
	}
	if len(results) != 1 {
		t.Error("mismatched length")
	}
	if !reflect.DeepEqual(results[0]["event_content"].(map[string]interface{}), new_content) {
		t.Error("mismatched content")
	}
}

func TestAddDeleteQueryEvents(t *testing.T) {
	if err := prepare_db_for_events(); err != nil {
		t.Error(err.Error())
	}
	content := map[string]interface{}{
		"type": "notification", "content": "test notification",
	}
	pk_values, err := utils.AddRow(map[string]interface{}{
		"email_id":      "oe2o950jgrnwgr",
		"email_address": "jake@example.com",
		"event_content": content,
	}, "events")
	if err != nil {
		t.Error(err.Error())
	}
	err = utils.DeleteRows(pk_values, "events")
	if err != nil {
		t.Error(err.Error())
	}
	results, err := utils.Query([]string{"event_content"}, map[string]interface{}{},
		"events")
	if err != nil {
		t.Error(err.Error())
	}
	if len(results) != 0 {
		t.Error("mismatched length")
	}
}
