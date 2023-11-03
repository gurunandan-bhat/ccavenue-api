package main

import (
	"ccavenue/client"
	"ccavenue/config"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func main() {

	cfg, err := config.Configuration()
	if err != nil {
		log.Fatalf("Error reading configuration: %s", err)
	}
	ccavClient, err := client.NewClient(cfg, 15*time.Second)
	if err != nil {
		log.Fatal("Error creating client: ", err)
	}

	orderStatus, err := ccavClient.Orders("30-10-2023", "01-11-2023")
	if err != nil {
		log.Fatal("Error from orders request: ", err)
	}

	statusBytes, err := json.MarshalIndent(orderStatus, "", "\t")
	if err != nil {
		log.Fatal("Error marshalling order status to JSON: ", err)
	}

	fmt.Println(string(statusBytes))
}
