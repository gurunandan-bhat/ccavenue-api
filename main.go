package main

import (
	"ccavenue/client"
	"ccavenue/config"
	"fmt"
	"log"
)

func main() {

	cfg, err := config.Configuration()
	if err != nil {
		log.Fatalf("Error reading configuration: %s", err)
	}
	ccavClient, err := client.NewClient(cfg, "1.2")
	if err != nil {
		log.Fatal("Error creating client: ", err)
	}

	filter := client.OrderLookupFilter{
		FromDate:       "31-10-2025",
		PageNumber:     1,
		OrderMaxAmount: 201.00,
	}

	jsonStr, err := ccavClient.Post(filter)
	if err != nil {
		log.Fatal("Error from orders request: ", err)
	}

	fmt.Println(string(*jsonStr))
}
