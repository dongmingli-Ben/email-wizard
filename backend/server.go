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

func main() {
	router := gin.Default()
	router.GET("/events", getEvents)

	router.Run("localhost:8080")
}
