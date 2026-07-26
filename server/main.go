package main

import (
	// "fmt"
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
	receive.ReceiverServer()
}