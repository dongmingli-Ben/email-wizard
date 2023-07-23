// Package main implements a client for GetEmail service.
package clients

import (
	"context"
	"encoding/json"
	"log"
	"time"

	pb "email-wizard/backend/clients/email_grpc_client"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EmailItem struct {
	Subject    string   `json:"subject"`
	Sender     string   `json:"sender"`
	Date       string   `json:"date"`
	Recipient  []string `json:"recipient"`
	Content    string   `json:"content"`
}

type Email struct {
	EmailID string    `json:"email_id"`
	Item    EmailItem `json:"item"`
}

type EmailCollections struct {
	Items []Email `json:"items"`
}

func GetEmails(user_config map[string]string, n_mails int32) (EmailCollections, error) {
	// fmt.Println("try requesting...")
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewEmailHelperClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	config_json, err := json.Marshal(user_config)
	// fmt.Printf("config: %v\n", user_config)
	if err != nil {
		log.Printf("%v cannot be jsonfied into json string", user_config)
		return EmailCollections{}, err
	}
	r, err := c.GetEmails(ctx, &pb.EmailRequest{Config: string(config_json), NMails: n_mails})
	if err != nil {
		log.Printf("Fail to call GetEmail over gRPC: %v", err)
		return EmailCollections{}, err
	}
	// log.Printf("response: %s", r.GetMessage())

	var emails EmailCollections
	err = json.Unmarshal([]byte(r.GetMessage()), &emails)
	if err != nil {
		log.Fatalf("%v cannot be de-serialized into emails, error: %v", r.GetMessage(), err)
		return EmailCollections{}, err
	}
	return emails, nil
}
