package receive
import (
	"fmt"
	"net"
	"encoding/json"
	"encoding/base64"
	// "os"
	"github.com/manpreet-singh1040/ByteBeam/server/fileHandling"
	"sync"
)


type Message struct{
	Type string `json:"type"`
	Data string `json:"data"`
	Metadata fileHandling.Metadata `json:"metadata"`
	ChunkIndex int `json:"chunkIndex"`
}

func ReceiverServer(){
	var mu sync.Mutex
	// var fileWriteMutex sync.Mutex
	addr := net.UDPAddr{
		Port: 8080,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error starting UDP server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("UDP server listening on", addr.String())
	buffer := make([]byte, 65535)
	var newFile fileHandling.File
	fileMap := make(map[int] int)
	remainChunks := make(map[int] int)
	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}
		message := string(buffer[:n])
		var msg Message
		err = json.Unmarshal([]byte(message), &msg)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			continue
		}
		if(msg.Type == "metadata"){
			newFile = fileHandling.File{Metadata: msg.Metadata}
			newFile.Initialize()
			for i:=0 ; i< newFile.Metadata.TotalChunks ; i++ {
				remainChunks[i] = 1
			}
			go RequestFile(newFile.Metadata,&fileMap,&remainChunks,&mu)
		}
		if(msg.Type =="file-chunk"){
			
			decodedData, err := base64.StdEncoding.DecodeString(msg.Data)
			if err != nil {
				fmt.Println("Error decoding base64 data:", err)
				continue
			}
			go newFile.WriteChunkData(decodedData,msg.ChunkIndex,&mu,&fileMap,&remainChunks)
		}
		fmt.Println(len(fileMap))
	}
}