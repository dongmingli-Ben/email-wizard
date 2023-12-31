// Package main implements a client for GetEmail service.
package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	pb "email-wizard/backend/clients/parse_grpc_client"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func parse_email_no_retry(email map[string]interface{}, timezone string) ([]map[string]string, error) {
	// fmt.Println("try requesting...")
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewParserClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()
	email_json, err := json.Marshal(email)
	// fmt.Printf("config: %v\n", user_config)
	if err != nil {
		log.Fatalf("%v cannot be jsonfied into json string", email)
		return []map[string]string{}, err
	}
	r, err := c.ParseEmail(ctx, &pb.EmailContentRequest{
		Email: string(email_json),
		AdditionalInfo: fmt.Sprintf(
			"{\"timezone\": \"%v\"}", timezone),
	})
	if err != nil {
		log.Printf("Fail to call ParseEmail over gRPC: %v", err)
		return []map[string]string{}, err
	}
	// log.Printf("response: %s", r.GetMessage())

	var response map[string]interface{}
	err = json.Unmarshal([]byte(r.GetMessage()), &response)
	if err != nil {
		log.Fatalf("%v cannot be de-serialized into emails, error: %v", r.GetMessage(), err)
		return []map[string]string{}, err
	}
	// collect as events
	events := make([]map[string]string, 0)
	if response["events"] == nil {
		fmt.Println("get empty events", response)
		return events, nil
	}
	raw_events := response["events"].([]interface{})
	// fmt.Println(events)
	for i := 0; i < len(raw_events); i = i + 1 {
		raw_event := raw_events[i].(map[string]interface{})
		event := make(map[string]string)
		for key, value := range raw_event {
			event[key] = value.(string)
		}
		events = append(events, event)
	}

	return events, nil
}

func ParseEmail(email map[string]interface{}, timezone string, max_retry int) ([]map[string]string, error) {
	events := make([]map[string]string, 0)
	success := false
	for i := 0; i < max_retry; i++ {
		es, err := parse_email_no_retry(email, timezone)
		if err != nil {
			log.Printf("parse gRPC client failed, retry with remaining count of %v", max_retry-i-1)
			continue
		}
		events = es
		success = true
		break
	}
	if success {
		return events, nil
	} else {
		return events, errors.New("parse gRPC client failed after retrying")
	}
}
