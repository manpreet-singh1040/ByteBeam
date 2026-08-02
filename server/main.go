package main

import (
	"fmt"
	// "os"
	// "net"
	// "encoding/json"
	"github.com/manpreet-singh1040/ByteBeam/server/receive"
	"github.com/joho/godotenv"
	"log"
)


func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Enter the details !!")
	var receiverIp string
	var receiverPort int
	var senderIp string
	var senderPort int

	fmt.Print("Receiver IP: ")
	if _, err := fmt.Scan(&receiverIp); err != nil {
		log.Fatal("Error reading Receiver IP:", err)
	}

	fmt.Print("Receiver Port: ")
	if _, err := fmt.Scan(&receiverPort); err != nil {
		log.Fatal("Error reading Receiver Port:", err)
	}

	fmt.Print("Sender IP: ")
	if _, err := fmt.Scan(&senderIp); err != nil {
		log.Fatal("Error reading Sender IP:", err)
	}

	fmt.Print("Sender Port: ")
	if _, err := fmt.Scan(&senderPort); err != nil {
		log.Fatal("Error reading Sender Port:", err)
	}

	fmt.Println("\n----- Input Received -----")
	fmt.Println("Receiver IP   :", receiverIp)
	fmt.Println("Receiver Port :", receiverPort)
	fmt.Println("Sender IP     :", senderIp)
	fmt.Println("Sender Port   :", senderPort)
	receive.ReceiverServer(senderIp)
}