package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"email-wizard/backend/utils"
)

func updateUserEvents(user_id string) error {

	accounts, err := utils.GetUserEmailAccounts(user_id)
	if err != nil {
		return err
	}

	// read recent emails from user email accounts
	emails, err := utils.GetUserEmailsFromAccounts(accounts)
	if err != nil {
		return err
	}

	// filter for un-parsed emails
	emails = utils.GetUserUnparsedEmails(emails, user_id)

	// parse into events
	events := utils.ParseEmailsToEvents(emails, 5)

	// store back to db
	err = utils.StoreUserEvents(events, user_id)
	return err
}

// getEvents reads user's email, parse them, read from db and return them all in one.
func getEvents(c *gin.Context) {
	user_id := c.Param("user_id")
	secret := c.Param("secret")
	if !utils.ValidateUserSecret(user_id, secret) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("wrong secret for user_id %v", user_id)})
		return
	}

	err := updateUserEvents(user_id)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}

	// read events from db
	events := utils.GetUserEvents(user_id)

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
			"protocol":    "IMAP",
			"username":    username,
			"password":    password,
			"imap_server": q.Get("imap_server"),
		}
		accounts = append(accounts, account)
	} else {
		account := map[string]string{
			"protocol":    "POP3",
			"username":    username,
			"password":    password,
			"imap_server": q.Get("imap_server"),
		}
		accounts = append(accounts, account)
	}
	emails, err := utils.GetUserEmailsFromAccounts(accounts)
	if err == nil {
		c.IndentedJSON(http.StatusOK, emails)
		return
	}
	c.IndentedJSON(http.StatusInternalServerError, gin.H{"errMsg": err.Error()})
}

func addUserMailbox(c *gin.Context) {
	var payload map[string]string
	if err := c.BindJSON(&payload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data"})
	}
	var mailbox_type string
	var user_id string
	var user_secret string
	var mailbox_address string
	if _mailbox_type, ok := payload["type"]; !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data: type"})
		return
	} else {
		mailbox_type = _mailbox_type
	}
	if _user_id, ok := payload["userId"]; !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data: userId"})
		return
	} else {
		user_id = _user_id
	}
	if _user_secret, ok := payload["userSecret"]; !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data: userSecret"})
		return
	} else {
		user_secret = _user_secret
	}
	if _mailbox_address, ok := payload["address"]; !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data: address"})
		return
	} else {
		mailbox_address = _mailbox_address
	}

	if !utils.ValidateUserSecret(user_id, user_secret) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"errMsg": fmt.Sprintf("wrong secret for user_id %v", user_id)})
		return
	}

	if mailbox_type == "outlook" {
		err := utils.AddUserMailboxOutlook(user_id, mailbox_address)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"errMsg": err.Error()})
			return
		}
	} else if mailbox_type == "IMAP" {
		var password string
		var imap_server string
		if _password, ok := payload["password"]; !ok {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data: password"})
			return
		} else {
			password = _password
		}
		if _imap_server, ok := payload["imap_server"]; !ok {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data: imap_server"})
			return
		} else {
			imap_server = _imap_server
		}
		err := utils.AddUserMailboxIMAP(user_id, mailbox_address, password, imap_server)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"errMsg": err.Error()})
			return
		}
	} else if mailbox_type == "POP3" {
		var password string
		var pop3_server string
		if _password, ok := payload["password"]; !ok {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data: password"})
			return
		} else {
			password = _password
		}
		if _pop3_server, ok := payload["pop3_server"]; !ok {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data: pop3_server"})
			return
		} else {
			pop3_server = _pop3_server
		}
		err := utils.AddUserMailboxPOP3(user_id, mailbox_address, password, pop3_server)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"errMsg": err.Error()})
			return
		}
	}
	c.IndentedJSON(http.StatusCreated, "")
}

func addUser(c *gin.Context) {
	var payload map[string]string
	if err := c.BindJSON(&payload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data"})
		return
	}
	username, ok_username := payload["username"]
	password, ok_password := payload["password"]
	if !(ok_username && ok_password) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data"})
		return
	}
	err := utils.AddUserDB(username, password)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, "")
}

func getUserProfile(c *gin.Context) {
	q := c.Request.URL.Query()
	user_id := q.Get("userId")
	user_secret := q.Get("userSecret")
	if !utils.ValidateUserSecret(user_id, user_secret) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"errMsg": fmt.Sprintf("wrong secret for user_id %v", user_id)})
		return
	}
	profile, err := utils.GetUserProfile(user_id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"errMsg": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, profile)
}

func authenticateUser(c *gin.Context) {
	q := c.Request.URL.Query()
	username := q.Get("username")
	password := q.Get("password")
	if !utils.ValidateUserPassword(username, password) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"errMsg": 
			fmt.Sprintf("wrong user name %v with password %v", username, password)})
		return
	}
	user_id, user_secret, err := utils.GetUserIdSecret(username)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"errMsg": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"user_id": user_id, "user_secret": user_secret})
}

func main() {
	router := gin.Default()

	// Add CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS", "DELETE", "PUT"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Methods",
	}
	router.Use(cors.New(config))
	router.GET("/events", getEvents)
	router.GET("/verify_email", getEmails)
	router.GET("/verify_user", authenticateUser)
	router.GET("/user_profile", getUserProfile)
	router.POST("/add_mailbox", addUserMailbox)
	router.POST("/add_user", addUser)

	// router.Run(":8080")
	router.RunTLS(":8080", "cert/www.toymaker-ben.online.pem", "cert/www.toymaker-ben.online.key")
}
