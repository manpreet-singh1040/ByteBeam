package fileShare
import(
	"fmt"
	"net"
	_"time"
	"os"
	"encoding/json"
	"github.com/manpreet-singh1040/ByteBeam/client/encoding"
	"github.com/manpreet-singh1040/ByteBeam/client/fileHandling"
)

func FileShare(){
	fmt.Println("hello from file share package")
	fmt.Println("hey bud nice work")
	conn, _ := net.Dial("udp", "0.0.0.0:8080")
	encoding.Encoding()
	defer conn.Close()
	// address:= "D:\\phone camera\\freasher vid\\VID_20230116_073608.mp4"
	address:= "D:\\series\\friends\\Season 4\\F.R.I.E.N.D.S S4-E8 The One with Chandler In a Box 720p.mkv"
	// address := "D:\\phone camera\\freasher vid\\VID_20221220_093101.mp4"
	File := fileHandling.File{}
	File.Initialize(address)
	mess := fileHandling.Message{
		Type: "metadata",
		Metadata: File.Metadata,
	}
	jsonData, err := json.Marshal(mess)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	conn.Write(jsonData)
	ListeningServer(File)
}	

func ListeningServer(File fileHandling.File){
	fmt.Println("hello from listening server")
	addr := net.UDPAddr{
		Port: 3000,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error starting UDP server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("UDP server listening on", addr.String())
	toalReq := make(map[int] struct{})
	buffer := make([]byte, 65535)
	fmt.Println("the total chunks are ",File.Metadata.TotalChunks)
	tempFile, err := os.OpenFile("D:\\series\\friends\\Season 4\\F.R.I.E.N.D.S S4-E8 The One with Chandler In a Box 720p.mkv", os.O_RDWR, 0644)
	if err!=nil {
		fmt.Println("error in opening the file")
		// panic(err)
	}
	newAddr,_ :=net.ResolveUDPAddr("udp", "0.0.0.0:8080")
	newConn , err := net.DialUDP("udp", nil, newAddr)
	if err != nil {
		fmt.Println("error in dialing the udp")
		// panic(err)
		// log.Print(err)
	}
	defer newConn.Close()
	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}
		// fmt.Printf("Received %d bytes from %s: %s\n", n, clientAddr, string(buffer[:n]))
		message := string(buffer[:n])
		var msg fileHandling.Message
		err = json.Unmarshal([]byte(message), &msg)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			continue
		}
		if(msg.Type == "chunk-request"){
			// fmt.Println("request for chunk received no." , msg.ChunkIndex)
			// toalReq = append(toalReq, msg.ChunkIndex)
			// if(msg.ChunkIndex== File.Metadata.TotalChunks-1) {
			// 	continue
			// }
			// start := time.Now()
			if msg.ChunkIndex != -1{
				toalReq[msg.ChunkIndex] = struct{}{}
				go File.ReadAndSendChunk(msg.ChunkIndex, newConn, newAddr,tempFile)
			}
			// end := time.Since(start)
			// fmt.Println("time taken to send the chunks is ", end)
			// fmt.Println("total req is:  ", len(toalReq))
			// File.ReadAndSendChunk(msg.ChunkIndex, conn, addr)
		}
	}
}