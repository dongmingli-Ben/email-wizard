package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"email-wizard/backend/logger"
	"email-wizard/backend/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

/*
TODO: update client auth flow to auth code flow and eliminate the need of kwargs in the request
*/
func updateAccountEventsAsync(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		fmt.Println(io.ReadAll(c.Request.Body))
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data"})
		return
	}
	user_id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": fmt.Sprintf("bad user_id: %v", c.Param("user_id"))})
		return
	}
	user_secret := c.Request.Header.Get("X-User-Secret")
	var email_address string
	var kwargs map[string]interface{}
	var ok bool
	if email_address, ok = payload["address"].(string); !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": fmt.Sprintf("address not found: %v", payload)})
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
	creds, err := utils.PrepareAndRefreshEmailAccountCredentials(user_id, account)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"errMsg": err.Error()})
		return
	}
	for key, val := range kwargs {
		creds[key] = val
	}
	account["credentials"] = creds
	err = utils.UpdateUserEventsForAccountAsync(user_id, account)
	// err = utils.UpdateUserEventsForAccount(user_id, account)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"errMsg": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, gin.H{"errMsg": ""})
}

// getEvents only read from events DB
func searchEvents(c *gin.Context) {
	user_id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("non-integer user_id: %v", c.Param("user_id"))})
		return
	}
	secret := c.Request.Header.Get("X-User-Secret")
	if ok, err := utils.ValidateUserSecret(user_id, secret); err != nil || !ok {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("wrong secret for user_id %v", user_id)})
		return
	}

	query := c.Request.URL.Query().Get("query")
	var events []map[string]interface{}
	if query == "" {
		// read events from db
		events, err = utils.GetUserEvents(user_id)
		if err != nil {
			fmt.Println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "fail to load events"})
			return
		}
	} else {
		// search from elastic
		events, err = utils.SearchUserEvents(user_id, query)
		if err != nil {
			fmt.Println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("fail to search for events with query: %s", query)})
			return
		}
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
			"protocol": "IMAP",
			"username": username,
			"credentials": map[string]interface{}{
				"password":    password,
				"imap_server": q.Get("imap_server"),
			},
		}
	} else {
		account = map[string]interface{}{
			"protocol": "POP3",
			"username": username,
			"credentials": map[string]interface{}{
				"password":    password,
				"imap_server": q.Get("imap_server"),
			},
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
	user_id, err := strconv.Atoi(c.Param(("user_id")))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": fmt.Sprintf("bad user_id: %v", c.Param("user_id"))})
		return
	}
	user_secret := c.Request.Header.Get("X-User-Secret")
	var mailbox_type string
	var mailbox_address string
	var credentials map[string]interface{}
	if _mailbox_type, ok := payload["type"].(string); !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data: type"})
		return
	} else {
		mailbox_type = _mailbox_type
	}
	if _mailbox_address, ok := payload["address"].(string); !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data: address"})
		return
	} else {
		mailbox_address = _mailbox_address
	}
	if _credentials, ok := payload["credentials"].(map[string]interface{}); !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data: address"})
		return
	} else {
		credentials = _credentials
	}

	if ok, err := utils.ValidateUserSecret(user_id, user_secret); err != nil || !ok {
		c.IndentedJSON(http.StatusForbidden, gin.H{"errMsg": fmt.Sprintf("wrong secret for user_id %v", user_id)})
		return
	}

	if err := utils.AddUserMailbox(user_id, mailbox_type, mailbox_address, credentials); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"errMsg": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, "")
}

func updateUserMailbox(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data"})
		return
	}
	user_id, err := strconv.Atoi(c.Param(("user_id")))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": fmt.Sprintf("bad user_id: %v", c.Param("user_id"))})
		return
	}
	user_secret := c.Request.Header.Get("X-User-Secret")
	mailbox_address := c.Param("address")
	var credentials map[string]interface{}
	if _credentials, ok := payload["credentials"].(map[string]interface{}); !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data: address"})
		return
	} else {
		credentials = _credentials
	}

	if ok, err := utils.ValidateUserSecret(user_id, user_secret); err != nil || !ok {
		c.IndentedJSON(http.StatusForbidden, gin.H{"errMsg": fmt.Sprintf("wrong secret for user_id %v", user_id)})
		return
	}

	if err := utils.UpdateUserMailbox(user_id, mailbox_address, credentials); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"errMsg": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, "")
}

func removeUserMailbox(c *gin.Context) {
	user_id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": fmt.Sprintf("bad user_id: %v", c.Param("user_id"))})
		return
	}
	address := c.Param("address")
	user_secret := c.Request.Header.Get("X-User-Secret")
	if ok, err := utils.ValidateUserSecret(user_id, user_secret); err != nil || !ok {
		c.IndentedJSON(http.StatusForbidden, gin.H{"errMsg": fmt.Sprintf("wrong secret for user_id %v", user_id)})
		return
	}

	if err := utils.RemoveUserMailbox(user_id, address); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"errMsg": err.Error()})
		return
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
	// defer logger.LogErrorStackTrace()
	// logger.Info("receives req for getUserProfile", zap.Any("params", c.Request.Body))
	user_id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": fmt.Sprintf("non-integer user_id: %v", c.Param("user_id"))})
	}
	user_secret := c.Request.Header.Get("X-User-Secret")
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
	// logger.Info("getUserProfile success")
}

func authenticateUser(c *gin.Context) {
	var payload map[string]interface{}
	// defer logger.LogErrorStackTrace()
	// req, _ := io.ReadAll(c.Request.Body)
	// logger.Info("receives req for authenticateUser", zap.String("req", string(req)))
	if err := c.BindJSON(&payload); err != nil {
		fmt.Println(io.ReadAll(c.Request.Body))
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data"})
		return
	}
	username, ok_username := payload["username"].(string)
	password, ok_password := payload["password"].(string)
	if !(ok_username && ok_password) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": "Invalid JSON data"})
		return
	}
	if ok, err := utils.ValidateUserPassword(username, password); err != nil || !ok {
		c.IndentedJSON(http.StatusForbidden, gin.H{"errMsg": fmt.Sprintf("fail to validate user name %v with password %v", username, password)})
		return
	}
	user_id, user_secret, err := utils.GetUserIdSecret(username)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"errMsg": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"user_id": user_id, "user_secret": user_secret})
	// logger.Info("authenticateUser success")
}

func main() {
	logger.InitLogger("log", "backend", 1, 7, "INFO")
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
		"X-User-Secret",
	}
	router.Use(cors.New(config))
	router.Use(logger.RequestLogger())
	router.GET("/users/:user_id/events", searchEvents)
	router.POST("/users/:user_id/events", updateAccountEventsAsync)
	router.GET("/verify_email", getEmails)
	router.POST("/authenticate", authenticateUser)
	router.GET("/users/:user_id/profile", getUserProfile)
	router.POST("/users/:user_id/mailboxes", addUserMailbox)
	router.DELETE("/users/:user_id/mailboxes/:address", removeUserMailbox)
	router.PUT("/users/:user_id/mailboxes/:address", updateUserMailbox)
	router.POST("/users", addUser)

	// websocket
	hub := utils.NewHub()
	router.GET("/ws/:user_id", func(c *gin.Context) {
		defer logger.LogErrorStackTrace()
		user_id, err := strconv.Atoi(c.Param("user_id"))
		user_secret := c.Request.Header.Get("X-User-Secret")

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.Error("fail to upgrade to websocket",
				zap.String("error", err.Error()))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"errMsg": "cannot upgrade to websocket"})
			return
		}

		if user_secret == "" {
			// expect the first message from client to be user_secret
			_, msg, err := conn.ReadMessage()
			if err != nil {
				logger.Error("fail to read message from client",
					zap.String("error", err.Error()))
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"errMsg": "cannot read message from client"})
				return
			}
			user_secret = string(msg)
			logger.Info("receives user_secret from client", zap.String("user_secret", user_secret))
		}

		if ok, err := utils.ValidateUserSecret(user_id, user_secret); err != nil || !ok {
			c.IndentedJSON(http.StatusForbidden, gin.H{"errMsg": fmt.Sprintf("wrong secret for user_id %v", user_id)})
			return
		}
		if err != nil {
			logger.Error("bad user_id",
				zap.String("error", err.Error()))
			c.IndentedJSON(http.StatusBadRequest, gin.H{"errMsg": fmt.Sprintf("bad user_id: %v", c.Param("user_id"))})
			return
		}
		utils.HandleWebSocketConnection(conn, user_id, hub)
	})

	// router.Run(":8080")
	router.RunTLS(":8080", "cert/www.toymaker-ben.online.pem", "cert/www.toymaker-ben.online.key")
}
