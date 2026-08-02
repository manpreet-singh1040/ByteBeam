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

func ReceiverServer(senderIp string){
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
	currentWindow := make(map[int] struct{})
	duplicatePackets :=0
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
			go RequestFile(newFile.Metadata,currentWindow,&mu,senderIp)
		}
		if(msg.Type =="file-chunk"){
			
			decodedData, err := base64.StdEncoding.DecodeString(msg.Data)
			if err != nil {
				fmt.Println("Error decoding base64 data:", err)
				continue
			}
			fileMap[msg.ChunkIndex] = 1
			// fmt.Println("the chunk number recieved is : ",msg.ChunkIndex  )
			// fmt.Println("the total chunks recieved are : ",len(fileMap)  )
			newFile.WriteChunkData(decodedData,msg.ChunkIndex,&mu,currentWindow,&duplicatePackets)
		}
		// fmt.Println(len(fileMap),"the duplicate are ",duplicatePackets)
	}
}