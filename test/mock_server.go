package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/public/quick-setup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)
		log.Printf("Received quick-setup request: %+v", req)

		ip, _ := req["ip"].(string)
		port := int(req["port"].(float64))
		protocol, _ := req["protocol"].(string)
		connectorName, _ := req["connector_name"].(string)

		resp := map[string]interface{}{
			"code":    200,
			"message": "success",
			"data": map[string]interface{}{
				"connector": map[string]interface{}{
					"id":             1,
					"connector_name": connectorName,
					"display_name":   connectorName,
					"api_port":       63042,
					"external_ip":    "127.0.0.1",
					"username":       connectorName,
					"password":       "test-secret-password",
				},
				"app": map[string]interface{}{
					"id":        10,
					"app_id":    "abc123def456-yishield",
					"site_name": fmt.Sprintf("%s-%s-%d", protocol, ip, port),
					"site_url":  fmt.Sprintf("https://abc123def456-yishield.ac-hostname.example.com"),
					"protocol":  protocol,
					"resource": map[string]interface{}{
						"ip":       "192.168.1.100",
						"port":     35021,
						"ac_id":    "ac-node-test",
						"hostname": "abc123def456-yishield.ac-hostname.example.com",
						"maskhost": true,
						"protocol": "tcp",
					},
				},
				"api_key": map[string]interface{}{
					"id":          5,
					"code":        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZXN0IjoibW9jayJ9.mock-signature",
					"nhp_server":  "nhp.test.example.com",
					"key_type":    "api_key",
					"expire_time": "2027-03-16T00:00:00Z",
					"app_id":      "abc123def456-yishield",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		log.Println("Responded with mock quick-setup data")
	})

	addr := "127.0.0.1:18080"
	log.Printf("Mock API server started at http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
