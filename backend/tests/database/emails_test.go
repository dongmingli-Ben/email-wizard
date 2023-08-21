package data_tests

import (
	"fmt"
	"reflect"
	"time"

	utils "email-wizard/backend/clients"
	"testing"

	_ "github.com/lib/pq"
)

func prepare_db() error {
	if err := utils.Reset(); err != nil {
		return err
	}
	_, err := utils.AddRow(map[string]interface{}{
		// "user_id":       1234323,
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
	results, err := utils.Query([]string{"user_id", "email_address", "event_ids"}, map[string]interface{}{"user_id": 1},
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
	date := time.Date(2023, time.August, 20, 12, 0, 0, 0, time.FixedZone("Asia/Shanghai", 8*60*60))
	pk_values, err := utils.AddRow(map[string]interface{}{
		"user_id":          1,
		"email_id":         "oe2o950jgrnwgr",
		"email_address":    "jake@example.com",
		"mailbox_type":     "outlook",
		"email_subject":    "example subject",
		"email_sender":     "salon@example.com",
		"email_recipients": []string{"jake@example.com", "sully@example.com"},
		"email_date":       date,
		"email_content":    "example content",
		"event_ids":        []int32{0, 1, 2},
	}, "emails")
	if err != nil {
		t.Error(err.Error())
	}
	results, err := utils.Query([]string{"user_id", "email_address", "event_ids", "email_date"},
		pk_values, "emails")
	if err != nil {
		t.Error(err.Error())
	}
	if len(results) != 1 {
		t.Error("mismatched length")
	}
	ret_ts, err := time.Parse(time.RFC3339, results[0]["email_date"].(string))
	if err != nil {
		t.Error(err.Error())
	}
	if results[0]["user_id"].(float64) != 1 ||
		results[0]["email_address"].(string) != "jake@example.com" ||
		!reflect.DeepEqual(results[0]["event_ids"].([]interface{}), []interface{}{0., 1., 2.}) ||
		ret_ts.In(date.Location()) != date {
		fmt.Println(results[0]["user_id"].(float64), results[0]["email_address"].(string), results[0]["event_ids"].([]interface{}))
		fmt.Println(results[0]["user_id"].(float64) != 1)
		fmt.Println(results[0]["email_address"].(string) != "jake@example.com")
		fmt.Println(!reflect.DeepEqual(results[0]["event_ids"].([]interface{}), []interface{}{0., 1., 2.}))
		fmt.Println(ret_ts.In(date.Location()) != date)
		t.Error("mismatched content")
	}
}

func TestAddUpdateQueryEmails(t *testing.T) {
	if err := prepare_db(); err != nil {
		t.Error(err.Error())
	}
	pk_values, err := utils.AddRow(map[string]interface{}{
		"user_id":          1,
		"email_id":         "oe2o950jgrnwgr",
		"email_address":    "jake@example.com",
		"mailbox_type":     "outlook",
		"email_subject":    "example subject",
		"email_sender":     "salon@example.com",
		"email_recipients": []string{"jake@example.com", "sully@example.com"},
		"email_date":       time.Date(2023, time.August, 20, 12, 0, 0, 0, time.FixedZone("Asia/Shanghai", 8*60*60)),
		"email_content":    "example content",
		"event_ids":        []int32{0, 1, 2},
	}, "emails")
	if err != nil {
		t.Error(err.Error())
	}
	err = utils.UpdateValue("event_ids",
		[]int32{4, 5, 6},
		pk_values,
		"emails")
	if err != nil {
		t.Error(err.Error())
	}
	results, err := utils.Query([]string{"event_ids"},
		pk_values, "emails")
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
	pk_values, err := utils.AddRow(map[string]interface{}{
		"user_id":          1,
		"email_id":         "oe2o950jgrnwgr",
		"email_address":    "jake@example.com",
		"mailbox_type":     "outlook",
		"email_subject":    "example subject",
		"email_sender":     "salon@example.com",
		"email_recipients": []string{"jake@example.com", "sully@example.com"},
		"email_date":       time.Date(2023, time.August, 20, 12, 0, 0, 0, time.FixedZone("Asia/Shanghai", 8*60*60)),
		"email_content":    "example content",
		"event_ids":        []int32{0, 1, 2},
	}, "emails")
	if err != nil {
		t.Error(err.Error())
	}
	err = utils.DeleteRows(pk_values,
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
