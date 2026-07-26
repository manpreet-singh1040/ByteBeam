package fileHandling

import (
	"fmt"
	"os"
	"time"
	// "sync"
	// "encoding/base64"
)

type Job struct{
	Data []byte
	ChunkIndex int
	StartTime time.Time
}

type FileWritter struct{
	File *File
	WritterQueue chan *Job
}

func (fw *FileWritter) HandleWriteWorker(job *Job,FileWriter *os.File ) (time.Duration) {
	start := time.Now()
	offset := int64(job.ChunkIndex * fw.File.Metadata.ChunkSize)
	// flag := false
	// var FileWriter *os.File
	// for !flag {
	// 	temp, ok := fw.File.TryAcquire()
	// 	FileWriter = temp
	// 	flag=ok
	// }
	_, err := FileWriter.WriteAt(job.Data, offset)
	if err != nil {
		fmt.Println("error writing chunk data:", err)
	}
	// fw.File.FilePool <- FileWriter
	end := time.Since(start)
	// fmt.Printf("Chunk %d written in %s\n", job.ChunkIndex, end)
	return end
}

func (fw *FileWritter) StartWorkerPool(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go func() {
			tempFilePointer, err := os.OpenFile(fw.File.Metadata.Name, os.O_RDWR, 0644)
			ii := 0
			totalTime := time.Duration(0)
			// totalIdealTime := time.Duration(0)
			// cuTime := time.Now()
			for job := range fw.WritterQueue {
				// totalIdealTime += time.Since(cuTime)
				// cuTime = time.Now()
				if err != nil {
					fmt.Println("error in opening the file")
					// panic(err)
				}
				totalTime+=  fw.HandleWriteWorker(job,tempFilePointer)
				ii++
				if ii%1000 == 0 {
					fmt.Printf("%s\n",totalTime)
					totalTime = 0
					// totalIdealTime = 0
				}
			}
		}()
	}
}

func (fw *FileWritter) EnqueueWriteJob(data []byte, chunkIndex int) {
	job := Job{
		Data:       data,
		ChunkIndex: chunkIndex,
		StartTime:  time.Now(),
	}
	fw.WritterQueue <- &job
}

