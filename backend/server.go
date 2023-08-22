package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"email-wizard/backend/utils"
)

func UpdateUserEventsForAccount(user_id int, account map[string]interface{}) error {
	// read recent emails from user email accounts (and store emails to DB)
	emails, err := utils.GetUserEmailsFromAccount(account)
	if err != nil {
		return err
	}

	// filter for un-parsed emails
	emails, err = utils.GetUserUnparsedEmails(emails, account["username"].(string), user_id)
	if err != nil {
		return err
	}
	for _, email := range emails {
		// todo: atomic!
		err = utils.StoreUserEmails([]map[string]interface{}{email}, account, user_id)
		if err != nil {
			return err
		}
		// parse into events
		events, err := utils.ParseEmailToEvents(email, 5)
		if err != nil {
			return err
		}
		// store back to db
		err = utils.StoreUserEvents(events, user_id,
			email["email_id"].(string),
			account["username"].(string))
		if err != nil {
			return err
		}
	}

	return nil
}

// reads new emails, parses them into events, and store them in DB
func updateAccountEvents(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		fmt.Println(io.ReadAll(c.Request.Body))
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data"})
		return
	}
	var email_address string
	var user_secret string
	var user_id int
	var kwargs map[string]interface{}
	var ok bool
	if email_address, ok = payload["address"].(string); !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": fmt.Sprintf("address not found: %v", payload)})
		return
	}
	if _user_id, ok := payload["user_id"].(float64); !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": fmt.Sprintf("user_id not found: %v", payload)})
		return
	} else {
		user_id = int(_user_id)
	}
	if user_secret, ok = payload["user_secret"].(string); !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": fmt.Sprintf("user_secret not found: %v", payload)})
		return
	}
	if kwargs, ok = payload["kwargs"].(map[string]interface{}); !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": fmt.Sprintf("kwargs not found: %v", payload)})
		return
	}
	if ok, err := utils.ValidateUserSecret(user_id, user_secret); !ok || err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "fail to authenticate user secret"})
		return
	}
	account, err := utils.GetUserEmailAccountFromAddress(user_id, email_address)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"errMsg": err.Error()})
		return
	}
	for key, val := range kwargs {
		account[key] = val
	}
	err = UpdateUserEventsForAccount(user_id, account)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"errMsg": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, gin.H{"errMsg": ""})
}

// getEvents only read from events DB
func getEvents(c *gin.Context) {
	q := c.Request.URL.Query()
	user_id, err := strconv.Atoi(q.Get("user_id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("non-integer user_id: %v", c.Param("user_id"))})
		return
	}
	secret := q.Get("user_secret")
	if ok, err := utils.ValidateUserSecret(user_id, secret); err != nil || !ok {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("wrong secret for user_id %v", user_id)})
		return
	}

	// read events from db
	events, err := utils.GetUserEvents(user_id)
	if err != nil {
		fmt.Println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "fail to load events"})
		return
	}

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
	var account map[string]interface{}
	if email_type == "IMAP" {
		account = map[string]interface{}{
			"protocol":    "IMAP",
			"username":    username,
			"password":    password,
			"imap_server": q.Get("imap_server"),
		}
	} else {
		account = map[string]interface{}{
			"protocol":    "POP3",
			"username":    username,
			"password":    password,
			"imap_server": q.Get("imap_server"),
		}
	}
	emails, err := utils.GetUserEmailsFromAccount(account)
	if err == nil {
		c.IndentedJSON(http.StatusOK, emails)
		return
	}
	c.IndentedJSON(http.StatusInternalServerError, gin.H{"errMsg": err.Error()})
}

func addUserMailbox(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data"})
		return
	}
	var mailbox_type string
	var user_id int
	var user_secret string
	var mailbox_address string
	if _mailbox_type, ok := payload["type"].(string); !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data: type"})
		return
	} else {
		mailbox_type = _mailbox_type
	}
	if _user_id, ok := payload["user_id"].(float64); !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data: user_id"})
		return
	} else {
		user_id = int(_user_id)
	}
	if _user_secret, ok := payload["user_secret"].(string); !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data: user_secret"})
		return
	} else {
		user_secret = _user_secret
	}
	if _mailbox_address, ok := payload["address"].(string); !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data: address"})
		return
	} else {
		mailbox_address = _mailbox_address
	}

	if ok, err := utils.ValidateUserSecret(user_id, user_secret); err != nil || !ok {
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
		if _password, ok := payload["password"].(string); !ok {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data: password"})
			return
		} else {
			password = _password
		}
		if _imap_server, ok := payload["imap_server"].(string); !ok {
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
		if _password, ok := payload["password"].(string); !ok {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data: password"})
			return
		} else {
			password = _password
		}
		if _pop3_server, ok := payload["pop3_server"].(string); !ok {
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
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		fmt.Println(io.ReadAll(c.Request.Body))
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data"})
		return
	}
	username, ok_username := payload["username"]
	password, ok_password := payload["password"]
	if !(ok_username && ok_password) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data"})
		return
	}
	err := utils.AddUserDB(username.(string), password.(string))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, "")
}

func getUserProfile(c *gin.Context) {
	q := c.Request.URL.Query()
	user_id, err := strconv.Atoi(q.Get("user_id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": fmt.Sprintf("non-integer user_id: %v", q.Get("user_id"))})
	}
	user_secret := q.Get("user_secret")
	if ok, err := utils.ValidateUserSecret(user_id, user_secret); err != nil || !ok {
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
	if ok, err := utils.ValidateUserPassword(username, password); err != nil || !ok {
		c.IndentedJSON(http.StatusForbidden, gin.H{"errMsg": fmt.Sprintf("fail to validate user name %v with password %v", username, password)})
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
	router.POST("/events", updateAccountEvents)
	router.GET("/verify_email", getEmails)
	router.GET("/verify_user", authenticateUser)
	router.GET("/user_profile", getUserProfile)
	router.POST("/add_mailbox", addUserMailbox)
	router.POST("/add_user", addUser)

	// router.Run(":8080")
	router.RunTLS(":8080", "cert/www.toymaker-ben.online.pem", "cert/www.toymaker-ben.online.key")
}
