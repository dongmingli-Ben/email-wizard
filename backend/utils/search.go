package utils

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func RequestElasticSearch(req *http.Request) (*http.Response, error) {
    // Username and password for HTTP basic authentication
	username := "elastic"
	password := "c7-FR3Fr2rWbbS*tDckZ"

	// Path to the CA certificate for SSL/TLS certificate verification
	caCertPath := "../search/config/http_ca.crt"

	// Load the CA certificate
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("error loading CA certificate: %v", err)
	}

	// Create a CA certificate pool and add the CA certificate to it
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a custom HTTP client with TLS configuration
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{Transport: tr}

    req.SetBasicAuth(username, password) // Add HTTP basic authentication headers

	// Set headers (optional)
	req.Header.Set("Content-Type", "application/json")

    // Send the request
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    return resp, nil
}


func AddEventToElasticSearch(event map[string]string, event_id int, user_id int) error {
	ready_event := make(map[string]interface{})
	for key, val := range event {
		ready_event[key] = val
	}
	ready_event["user_id"] = user_id
	ready_event["event_id"] = event_id

	url := fmt.Sprintf("https://localhost:9200/email-wizard-events/_doc/%d", event_id)

    // JSON data to send in the request body
    event_json, err := json.Marshal(ready_event)
    if err != nil {
        return err
    }

    // Create a POST request
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(event_json))
    if err != nil {
        return err
    }

    resp, err := RequestElasticSearch(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Read the response body
    _, err = io.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    // Print the response status code and body
    // fmt.Println("Response Status Code:", resp.Status)
    // fmt.Println("Response Body:", string(body))
    return nil
}

func SearchUserEvents(user_id int, query string) ([]map[string]interface{}, error) {
    // Elasticsearch endpoint URL
	url := "https://localhost:9200/email-wizard-events/_search"

	// Construct the Elasticsearch query
	elastic_query := `{
	  "query": {
	    "bool": {
	      "must": [
	        {
	          "term": {
	            "user_id": "` + fmt.Sprint(user_id) + `"
	          }
	        },
	        {
	          "wildcard": {
	            "summary": {
                    "case_insensitive": true,
                    "value": "*` + query + `*"
                }
	          }
	        }
	      ]
	    }
	  }
	}`

	req, err := http.NewRequest("POST", url, strings.NewReader(elastic_query))
	if err != nil {
		return nil, err
	}
    req.Method = "GET"
    resp, err := RequestElasticSearch(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("elasticsearch responded with status code %d", resp.StatusCode)
	}

	// Read and return the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
    var result map[string]interface{}
    err = json.Unmarshal(body, &result)
    if err != nil {
        return nil, err
    }

    events := make([]map[string]interface{}, 0)
    res := result["hits"].(map[string]interface{})["hits"].([]interface{})
    for _, content := range res {
        event := make(map[string]interface{})
        for key, val := range content.(map[string]interface{})["_source"].(map[string]interface{}) {
            if key == "user_id" || key == "event_id" {
                continue
            }
            event[key] = val
        }
        events = append(events, event)
    }

	return events, nil
}