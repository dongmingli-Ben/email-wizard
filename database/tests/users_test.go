package data_tests

import (
	"fmt"
	"os/exec"

	"email-wizard/data/utils"
	"testing"

	_ "github.com/lib/pq"
)

func reset() error {
	cmd := exec.Command("sh", "../reset_db.sh")
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Println(string(output))
		return err
	}
	cmd = exec.Command("sh", "../init_db.sh")
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Println(string(output))
		return err
	}
	return nil
}

func TestQueryUsers(t *testing.T) {
	if err := reset(); err != nil {
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
	if err := reset(); err != nil {
		t.Error(err.Error())
	}
	err := utils.AddRow(map[string]interface{}{
		"user_id":       1234323,
		"user_secret":   "oe2o950jgrnwgr",
		"user_name":     "jake",
		"user_password": "sjgn",
		"mailboxes":     []map[string]interface{}{},
	}, "users")
	if err != nil {
		t.Error(err.Error())
	}
	results, err := utils.Query([]string{"user_id", "user_secret"}, map[string]interface{}{"user_id": 1234323},
		"users")
	if err != nil {
		t.Error(err.Error())
	}
	if len(results) != 1 {
		t.Error("mismatched length")
	}
	if results[0]["user_id"].(int64) != 1234323 || results[0]["user_secret"].(string) != "oe2o950jgrnwgr" {
		t.Error("mismatched content")
	}
}

func TestAddUpdateQueryUsers(t *testing.T) {
	if err := reset(); err != nil {
		t.Error(err.Error())
	}
	mailbox1 := map[string]interface{}{
		"address": "asdjr",
	}
	err := utils.AddRow(map[string]interface{}{
		"user_id":       1234323,
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
		map[string]interface{}{"user_id": 1234323},
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
	if results[0]["mailboxes"].([]map[string]interface{})[0]["address"] != "asdjr" {
		t.Error("mismatched content")
	}
}

func TestAddDeleteQueryUsers(t *testing.T) {
	if err := reset(); err != nil {
		t.Error(err.Error())
	}
	err := utils.AddRow(map[string]interface{}{
		"user_id":       1234323,
		"user_secret":   "oe2o950jgrnwgr",
		"user_name":     "jake",
		"user_password": "sjgn",
		"mailboxes":     []map[string]interface{}{},
	}, "users")
	if err != nil {
		t.Error(err.Error())
	}
	err = utils.DeleteRows(map[string]interface{}{"user_id": 1234323},
		"users")
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
