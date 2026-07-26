package fileHandling

import (
	"fmt"
	"os"
	"sync"
	"log"
	"strconv"
	// "encoding/base64"
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
	FileWritter *FileWritter
}


func (p *File) TryAcquire() (*os.File, bool) {
    select {
    case f := <-p.FilePool:
        return f, true
    default:
        return nil, false // no file available
    }
}

func (f *File) Initialize()  {
	fmt.Println("hello from file handling package ")
	fileSize:= f.Metadata.Size
	newFile, err :=os.Create(f.Metadata.Name)
	if err!=nil {
		fmt.Println("error in creating the file")
		panic(err)
	}
	// defer newFile.Close()
	newFile.Seek(fileSize-1,0)
	newFile.Write([]byte{0})
	fmt.Println("file created with size:", fileSize);
	value := os.Getenv("FILE_POINTER_POOL_SIZE")

     filePointerPoolSize, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Invalid FILE_POINTER_POOL_SIZE: %v", err)
	}
	f.FilePool = make(chan *os.File,filePointerPoolSize)
	for i:=0 ; i<filePointerPoolSize;i++ {
		tempFile, err := os.OpenFile(f.Metadata.Name, os.O_RDWR, 0644)
		
		if err!=nil {
			fmt.Println("error in opening the file")
			// panic(err)
		}
		f.FilePool <- tempFile
	}
	value2 := os.Getenv("FILE_WRITE_JOB_QUEUE_SIZE")
	 fileWriteJobQueueSize, err := strconv.Atoi(value2)
	if err != nil {
		log.Fatalf("Invalid FILE_WRITE_JOB_QUEUE_SIZE: %v", err)
	}
	newFileWritter := &FileWritter{
		File: f,
		WritterQueue: make(chan *Job, fileWriteJobQueueSize),
	}
	f.FileWritter = newFileWritter
	value3 := os.Getenv("FILE_WRITER_WORKER_POOL_SIZE")
	 fileWriterWorkerPoolSize, err := strconv.Atoi(value3)
	if err != nil {
		log.Fatalf("Invalid FILE_WRITER_WORKER_POOL_SIZE: %v", err)
	}
	go f.FileWritter.StartWorkerPool(fileWriterWorkerPoolSize) // Start worker goroutines
	// f.File = newFile
}

func (f *File)  WriteChunkData(data []byte , chunkIndex int, mu *sync.Mutex , curretWindow map[int]struct{}, duplicatePackets *int)  {

	// offset := int64(chunkIndex * f.Metadata.ChunkSize)
	// fileWriteMutex.Lock()
	// f.File.Seek(offset,0)
	// fmt.Println("the data is:", data)
	// fmt.Println("the offset is:",offset)
	// _, err := f.File.Write(data)

	//this needs to be uncommented
	// FileWriter, ok := f.TryAcquire()
	// if !ok {
	// 	// fmt.Println("no file available for writing")
	// 	return
	// }
	// // FileWriter := <- f.FilePool
	// // fmt.Println("i am printing chunk index",chunkIndex)
	// _, err := FileWriter.WriteAt(data, offset)
	// if err != nil {
	// 	fmt.Println("error writing chunk data:", err)
	// 	// fileWriteMutex.Unlock()
	// 	f.FilePool <- FileWriter
	// }
	// f.FilePool <- FileWriter
	// fmt.Println("finished chunk index",chunkIndex)
	// fileWriteMutex.Unlock()
	// fmt.Println("written", n, "bytes of chunk at index:", chunkIndex)
	mu.Lock()
			// size = len(*fileMap)
		// (*fileMap)[chunkIndex] = 1
			// delete(currentWindow, chunkIndex)
		_, exists := curretWindow[chunkIndex]
		if exists {
			delete(curretWindow, chunkIndex)
			f.FileWritter.EnqueueWriteJob(data, chunkIndex)
			} else {
			*duplicatePackets++
		}
    mu.Unlock()
}
