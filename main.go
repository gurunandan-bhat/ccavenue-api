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

	filter := client.OrderFilter{
		FromDate:   "22-10-2023",
		OrderEmail: "veena.cpsi@gmail.com",
	}

	orders, err := ccavClient.OrderLookup(filter)
	if err != nil {
		log.Fatal("Error from orders request: ", err)
	}

	orderBytes, err := json.MarshalIndent(orders, "", "\t")
	if err != nil {
		log.Fatal("Error marshalling order status to JSON: ", err)
	}

	fmt.Println(string(orderBytes))
}
