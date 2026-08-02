package receive

import (
	_"maps"
	"encoding/json"
	"fmt"
	"net"
	"time"
	"sync"
	"github.com/manpreet-singh1040/ByteBeam/server/fileHandling"
)

func RequestFile(newFileMetadata fileHandling.Metadata , currentWindow map[int] struct{}, mu *sync.Mutex,senderIp string) {
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
	// for i:=0 ; i< newFileMetadata.TotalChunks ; i+=10 {
	// 	endIdx := i + 9
	// 	if endIdx >= newFileMetadata.TotalChunks {
	// 		endIdx = newFileMetadata.TotalChunks -1
	// 	}
	// 	RequestFilePart(newFileMetadata , i , endIdx , currentWindow , mu)
	// }
	start := time.Now()
	RequestFilePart(newFileMetadata , 0 , newFileMetadata.TotalChunks-1 , currentWindow , mu, senderIp)
	end := time.Since(start)
	fmt.Printf("File transfer completed in %s. All chunks received. \n", end)
	fmt.Println("File transfer completed. All chunks received. ", len(currentWindow), "chunks received  total chunks ", newFileMetadata.TotalChunks)

}

func RequestFilePart(newFileMetadata fileHandling.Metadata , startIdx int , endIdx int,currentWindow map[int]struct{}, mu *sync.Mutex,senderIp string) {
	// currentWindow := make(map[int] int)
	// for n < Size {
	// 	for n < Size && len(currentWindow) < windowSize {
	// 		if(tillNow[n] == 0){
	// 			(currentWindow)[startIdx + n] = struct{}{}
	// 		}
	// 	}
	// 	for i:=startIdx ; i<= endIdx ; i++ {
	// 		mu.Lock()
	// 		_, exists := (*recievedChunks)[i]
	// 		mu.Unlock()
	// 		if exists{
	// 			tillNow[i - startIdx]++
	// 			if(tillNow[i - startIdx] ==1){
	// 				n++;
	// 			}
	// 			continue
	// 		}
	// 		// n++;
	// 		// tillNow[i - startIdx]++
	// 		newMessage := Message{
	// 				Type: "chunk-request",
	// 				ChunkIndex: i,
	// 				Metadata: newFileMetadata,
	// 			}
	// 		jsonMessage , err := json.Marshal(newMessage)
	// 		if err != nil {
	// 			fmt.Println("error in marshaling the json")
	// 			panic(err)
	// 		}
	// 		_, err = conn.Write(jsonMessage)
	// 		if err != nil {
		// 			fmt.Println("error in sending the json")
		// 				// panic(err)
	// 		}
	
	// 	}
	// 	// fmt.Println("done from ",startIdx," to ",endIdx, "the n is ", n)
	// }
	
	addr := net.UDPAddr{
		Port: 3000,
		IP:   net.ParseIP(senderIp),
	}
	conn , err := net.DialUDP("udp", nil, &addr)
	if err != nil {
		fmt.Println("error in dialing the udp")
		panic(err)
	}
	defer conn.Close()
	Size := endIdx - startIdx + 1
	// tillNow := make([]int64, Size)
	n:=0
	for n<Size {
		windowSize := 1000
		tempArray := make([]int, windowSize)
		tempi := 0
		mu.Lock()
		for key := range currentWindow {
			tempArray[tempi] = key
			tempi++
		}
		// flag := 0
		for len(currentWindow) < windowSize && n < Size {
			currentWindow[startIdx + n] = struct{}{}
			tempArray[tempi] = startIdx + n
			tempi++
			n++
		}

		for tempi < windowSize && n == Size {
			tempArray[tempi] = -1
			tempi++
		}
		// fmt.Println("the temp array is -->",tempArray)
		mu.Unlock()
		for i:=0 ; i< len(tempArray) ; i++ {
			if tempArray[i] == -1 {
				continue
			}
			newMessage := Message{
				Type: "chunk-request",
				ChunkIndex: tempArray[i],
				Metadata: newFileMetadata,
			}
			jsonMessage , err := json.Marshal(newMessage)
			if err != nil {
				// mu.Unlock()
				fmt.Println("error in marshaling the json")
				panic(err)
			}
			_, err = conn.Write(jsonMessage)
			if err != nil {
				fmt.Println("error in sending the json")
				panic(err)
			}
		}
		// println("the n is --->",n," the size is --->",Size)
		time.Sleep(40 * time.Millisecond)
	}
	flag := true
	for flag {
		tempArray := make([]int, len(currentWindow))
		mu.Lock()
		if len(currentWindow) == 0 {
			flag = false
			mu.Unlock()
			continue
		}
		i := 0
		for key := range currentWindow {
			if i >= len(tempArray) {
				break
			}
			tempArray[i] = key
			i++
		}
		mu.Unlock()
		for i:=0 ; i< len(tempArray) ; i++ {
			if tempArray[i] == -1 {
				continue
			}
			newMessage := Message{
				Type: "chunk-request",
				ChunkIndex: tempArray[i],
				Metadata: newFileMetadata,
			}
			jsonMessage , err := json.Marshal(newMessage)
			if err != nil {
				// mu.Unlock()
				fmt.Println("error in marshaling the json")
				panic(err)
			}
			_, err = conn.Write(jsonMessage)
			if err != nil {
				fmt.Println("error in sending the json")
				panic(err)
			}
		}
	}
	// fmt.Println("DONE_________________________________________________________")
}