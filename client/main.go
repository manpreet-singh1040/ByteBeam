package main

import (
	_"encoding/json"
	_"fmt"
	_"net"

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
	fileShare.FileShare()
}
