package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func updateUserEvents(user_id string) error {

	var accounts []map[string]string = getUserEmailAccounts(user_id)

	// read recent emails from user email accounts
	emails, err := getUserEmailsFromAccounts(accounts)
	if err != nil {
		return err
	}

	// filter for un-parsed emails
	emails = getUserUnparsedEmails(emails, user_id)

	// parse into events
	events := parseEmailsToEvents(emails, 5)

	// store back to db
	err = storeUserEvents(events, user_id)
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
	emails, err := getUserEmailsFromAccounts(accounts)
	if err == nil {
		c.IndentedJSON(http.StatusOK, emails)
		return
	}
	c.IndentedJSON(http.StatusInternalServerError, gin.H{"errMsg": err.Error()})
}

// func CORSMiddleware() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
//         c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
//         c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
//         c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

//         if c.Request.Method == "OPTIONS" {
//             c.AbortWithStatus(204)
//             return
//         }

//         c.Next()
//     }
// }

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

	router.Run(":8080")
}
