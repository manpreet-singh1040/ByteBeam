package fileHandling

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)
type Metadata struct{
	Name string
	Size int64
	Type string
	ChunkSize int
	TotalChunks int
}

type File struct{
	Metadata Metadata
	FilePool chan *os.File
}

type Message struct{
	Type string `json:"type"`
	Data string `json:"data"`
	Metadata Metadata `json:"metadata"`
	ChunkIndex int `json:"chunkIndex"`
}

func (f *File) TryAcquire() (*os.File, bool) {
	select{
	case temp:= <- f.FilePool: 
		return temp, true
	default:
		return nil, false
	}
}

func (f *File) Initialize(address string) {
	fmt.Println("hello from file handling package")
	file,err := os.Open(address)
	if err!=nil {
		fmt.Println("error in opening the file")
		// panic(err)
		fmt.Println(err)
	}
	fileInfo,err := file.Stat()
	if err!=nil {
		fmt.Println("error in getting file info")
		panic(err)
	}
	Metadata := Metadata{
		Name: fileInfo.Name(),
		Size: fileInfo.Size(),
		Type: "txt",
		ChunkSize: 1024,
		TotalChunks: int(fileInfo.Size()/1024) +1,
	}
	f.Metadata = Metadata
	fmt.Println("the file meta data is ",Metadata)
	f.FilePool = make( chan *os.File ,4)
	file.Close()
	for i:=0 ; i<4;i++  {
		tempFile, err := os.OpenFile(address, os.O_RDWR, 0644)
		if err!=nil {
			fmt.Println("error in opening the file")
			// panic(err)
		}
		f.FilePool <- tempFile
	}
	// defer file.Close()
}

func (f *File) ReadAndSendChunk(curretChunk int , conn *net.UDPConn, addr *net.UDPAddr) {
		newConn , err := net.DialUDP("udp", nil, addr)
		if err != nil {
			fmt.Println("error in dialing the udp")
			// panic(err)
			log.Print(err)
		}
		defer newConn.Close()
		var offset int64
		var buffer []byte
		if curretChunk == f.Metadata.TotalChunks-1 {
			// Last chunk: calculate offset and buffer size for remaining bytes
			offset = int64(curretChunk) * int64(f.Metadata.ChunkSize)
			remaining := f.Metadata.Size - offset
			buffer = make([]byte, int(remaining))
			} else {
				// Not last chunk: standard chunk size
				offset = int64(curretChunk) * int64(f.Metadata.ChunkSize)
				buffer = make([]byte, f.Metadata.ChunkSize)
			}
			FilePointer, ok := f.TryAcquire()
			if !ok {
				fmt.Println("error in acquiring the file pointer")
				// panic(err)
				return
			}
			_, err = FilePointer.Seek(offset, 0)
			if err != nil {
				fmt.Println("error in seeking to chunk offset")
				// panic(err)
				fmt.Println(err)
				fmt.Println(offset)
			}
			n, err := FilePointer.Read(buffer)
			if err != nil {
				// fmt.Println("error in reading the file")
				// log.Print(err)
				// fmt.Println("sending chunk number ", curretChunk," the total chunks are ", f.Metadata.TotalChunks)
			// panic(err)
		}
		encodedChunk := base64.StdEncoding.EncodeToString(buffer[:n])
		newMsg :=Message{
			Type: "file-chunk",
			ChunkIndex: curretChunk,
			Data: encodedChunk,
			Metadata: f.Metadata,
		}
		jsonData, err := json.Marshal(newMsg)
		if err != nil {
			fmt.Println("error in marshaling json")
			panic(err)
		}
		newConn.Write(jsonData)
		f.FilePool <- FilePointer
		// fmt.Printf("Sent chunk %d/%d\n", curretChunk, f.Metadata.TotalChunks)
}


func FileHandler (file *File) {
	fmt.Println("hello from file handler")

}