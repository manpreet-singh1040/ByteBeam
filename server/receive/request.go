package receive

import (
	_"maps"
	"encoding/json"
	"fmt"
	"net"
	_"time"
	"sync"
	"github.com/manpreet-singh1040/ByteBeam/server/fileHandling"
)

func RequestFile(newFileMetadata fileHandling.Metadata , recievedChunks *map[int]int , remainChunks *map[int]int , mu *sync.Mutex) {
	// fmt.Println("Requesting file from server... total chunks:", newFileMetadata.TotalChunks)
	// addr := net.UDPAddr{
	// 	Port: 3000,
	// 	IP:   net.ParseIP("0.0.0.0"),
	// }
	// conn , err := net.DialUDP("udp", nil, &addr)
	// if err != nil {
	// 	fmt.Println("error in dialing the udp")
	// 	panic(err)
	// }
	// defer conn.Close()
	// for len(*recievedChunks) < newFileMetadata.TotalChunks {
	// 	for i:=0 ; i< newFileMetadata.TotalChunks ; i++ {
	// 		mu.Lock()
	// 		_, exists := (*recievedChunks)[i]
	// 		if exists{
	// 			mu.Unlock()
	// 			continue
	// 		}
	// 		mu.Unlock()
	// 		newMessage := Message{
	// 			Type: "chunk-request",
	// 			ChunkIndex: i,
	// 			Metadata: newFileMetadata,
	// 		}
	// 		jsonMessage , err := json.Marshal(newMessage)
	// 		if err != nil {
	// 			fmt.Println("error in marshaling the json")
	// 			panic(err)
	// 		}
	// 		_, err = conn.Write(jsonMessage)
	// 		if err != nil {
	// 			fmt.Println("error in sending the json")
	// 			// panic(err)
	// 		}
	// 		// time.Sleep(2 * time.Millisecond)
	// 	}
	// 	// time.Sleep(100 * time.Millisecond)
	// }
	for i:=0 ; i< newFileMetadata.TotalChunks ; i+=10000 {
		endIdx := i + 9999
		if endIdx >= newFileMetadata.TotalChunks {
			endIdx = newFileMetadata.TotalChunks -1
		}
		go RequestFilePart(newFileMetadata , i , endIdx , recievedChunks , remainChunks , mu)
	}
	fmt.Println("File transfer completed. All chunks received. ", len(*recievedChunks), "chunks received  total chunks ", newFileMetadata.TotalChunks)

}

func RequestFilePart(newFileMetadata fileHandling.Metadata , startIdx int , endIdx int,recievedChunks *map[int]int , remainChunks *map[int]int , mu *sync.Mutex) {
	addr := net.UDPAddr{
		Port: 3000,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn , err := net.DialUDP("udp", nil, &addr)
	if err != nil {
		fmt.Println("error in dialing the udp")
		panic(err)
	}
	defer conn.Close()
	Size := endIdx - startIdx + 1
	tillNow := make([]int64, Size)
	n:=0
	for n < Size {

		for i:=startIdx ; i<= endIdx ; i++ {
			mu.Lock()
			_, exists := (*recievedChunks)[i]
			if exists{
				mu.Unlock()
				tillNow[i - startIdx]++
				if(tillNow[i - startIdx] ==1){
					n++;
				}
				continue
			}
			// n++;
			// tillNow[i - startIdx]++
			mu.Unlock()
			newMessage := Message{
					Type: "chunk-request",
					ChunkIndex: i,
					Metadata: newFileMetadata,
				}
				jsonMessage , err := json.Marshal(newMessage)
				if err != nil {
					fmt.Println("error in marshaling the json")
					panic(err)
				}
				_, err = conn.Write(jsonMessage)
				if err != nil {
					fmt.Println("error in sending the json")
					// panic(err)
				}
	
		}
		// fmt.Println("done from ",startIdx," to ",endIdx, "the n is ", n)
	}
}