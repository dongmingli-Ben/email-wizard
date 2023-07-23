package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func updateUserEvents(user_id string) error {

	var accounts []map[string]string = getUserEmailAccounts(user_id)

	// read recent emails from user email accounts
	emails := getUserEmailsFromAccounts(accounts)

	// filter for un-parsed emails
	emails = getUserUnparsedEmails(emails, user_id)

	// parse into events
	events := parseEmailsToEvents(emails, 5)

	// store back to db
	err := storeUserEvents(events, user_id)
	return err
}

// getEvents reads user's email, parse them, read from db and return them all in one.
func getEvents(c *gin.Context) {
	user_id := c.Param("user_id")
	secret := c.Param("secret")
	if !validateUserSecret(user_id, secret) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("wrong secret for user_id %v", user_id)})
		return
	}

	err := updateUserEvents(user_id)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}

	// read events from db
	events := getUserEvents(user_id)

	c.IndentedJSON(http.StatusOK, events)
}

// try to get user's email based on provided credentials, only support IMAP (and POP3)
func getEmails(c *gin.Context) {
	q := c.Request.URL.Query()
	username := q.Get("username")
	password := q.Get("password")
	email_type := q.Get("type")
	if email_type != "IMAP" && email_type != "POP3" {
		fmt.Println(email_type, "not IMAP or POP3")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "type only accepts IMAP and POP3"})
		return
	}
	accounts := make([]map[string]string, 0)
	if email_type == "IMAP" {
		account := map[string]string{
			"protocol": "IMAP",
			"username": username,
			"password": password,
			"imap_server": q.Get("imap_server"),
		}
		accounts = append(accounts, account)
	} else {
		account := map[string]string{
			"protocol": "POP3",
			"username": username,
			"password": password,
			"imap_server": q.Get("imap_server"),
		}
		accounts = append(accounts, account)	
	}
	emails := getUserEmailsFromAccounts(accounts)
	c.IndentedJSON(http.StatusOK, emails)
}

func main() {
	router := gin.Default()
	router.GET("/events", getEvents)
	router.GET("/verify_email", getEmails)

	router.Run("localhost:8080")
}
