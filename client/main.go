package main

import (
	_"encoding/json"
	"fmt"
	_"net"
	"bufio"
	"os"
	"strings"
	// _"github.com/manpreet-singh1040/ByteBeam/client/encoding"
	// _"github.com/manpreet-singh1040/ByteBeam/client/fileHandling"
	"github.com/manpreet-singh1040/ByteBeam/client/fileShare"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	reader := bufio.NewReader(os.Stdin)

	// var receiverIp string
	// var receiverPort int
	var senderIp string
	// var senderPort int
	// var filePathToSend string

	// fmt.Print("Receiver IP: ")
	// if _, err := fmt.Scan(&receiverIp); err != nil {
	// 	log.Fatal("Error reading Receiver IP:", err)
	// }

	// fmt.Print("Receiver Port: ")
	// if _, err := fmt.Scan(&receiverPort); err != nil {
	// 	log.Fatal("Error reading Receiver Port:", err)
	// }

	fmt.Print("Sender IP: ")
	senderIp, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading Sender IP:", err)
	}
	senderIp = strings.TrimSpace(senderIp)

	// fmt.Print("Sender Port: ")
	// if _, err := fmt.Scan(&senderPort); err != nil {
	// 	log.Fatal("Error reading Sender Port:", err)
	// }

	fmt.Print("File Path: ")
	filePathToSend, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading File Path:", err)
	}
	filePathToSend = strings.TrimSpace(filePathToSend)

	fmt.Println("\n----- Input Received -----")
	// fmt.Println("Receiver IP   :", receiverIp)
	// fmt.Println("Receiver Port :", receiverPort)
	fmt.Println("Sender IP     :", senderIp)
	// fmt.Println("Sender Port   :", senderPort)
	fmt.Println("File Path     :", filePathToSend)

	fileShare.FileShare(senderIp, filePathToSend)
}