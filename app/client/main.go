package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mvgmb/BigFruit/app/proto/storage_object"
	"github.com/mvgmb/BigFruit/client"
	"github.com/mvgmb/BigFruit/proto/naming"
	"github.com/mvgmb/BigFruit/util"
)

func main() {
	lookupManyRequest := &naming.LookupManyRequest{
		ServiceName: "StorageObject",
		NumberOfAor: 1,
	}
	res, err := client.LookupMany(lookupManyRequest)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(res.AorList))

	var options []*util.Options

	for i := range res.AorList {
		option := &util.Options{
			Host:     res.AorList[i].Host,
			Port:     uint16(res.AorList[i].Port),
			Protocol: "tcp",
		}
		options = append(options, option)
	}
	noTests := 1
	results := make([]time.Duration, noTests)

	dataSize := int32(645794886)
	total := float64(0)
	for i := 0; i < noTests; i++ {
		results[i] = download("/home/mario/Desktop/data2.zip", "data2.zip", dataSize, 21500, options)
		fmt.Println(results[i])
		total += results[i].Seconds()
	}
	fmt.Println(results)
	fmt.Println("avg", total/float64(len(results)))
}

func download(filePath, outputPath string, noBytes, offset int32, options []*util.Options) time.Duration {
	reqCh := make(chan proto.Message)
	resCh := make(chan proto.Message)

	go func() {
		var i int32
		for i = 0; i < noBytes; i += offset {
			reqCh <- &storage_object.DownloadRequest{
				FilePath: filePath,
				Start:    int64(i),
				Offset:   offset,
			}
		}
		close(reqCh)
	}()

	done := make(chan bool)
	go func() {
		file, err := os.Create(outputPath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		cur := int64(0)
		for {
			res, more := <-resCh
			if more {
				response := res.(*storage_object.DownloadResponse)
				_, err = file.WriteAt(response.Bytes, cur)
				if err != nil {
					log.Println(err)
				}
				cur += int64(binary.Size(response.Bytes))
			} else {
				break
			}
		}
		done <- true

	}()
	bigFruit := client.NewBigFruit()

	t := time.Now()
	err := bigFruit.Invoke("StorageObject", "Download", options, false, reqCh, resCh)
	if err != nil {
		done <- true
		log.Println(err)
	}
	<-done
	return time.Since(t)
}

func upload(filePath, outputPath string, noBytes, offset int32, options []*util.Options) time.Duration {
	reqCh := make(chan proto.Message)
	resCh := make(chan proto.Message)

	go func() {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		noReq := 0
		var cur int32
		for cur = 0; cur < noBytes; cur += offset {
			_, err = file.Seek(int64(cur), 0)
			if err != nil {
				log.Println(err)
			}

			buffer := make([]byte, offset)
			noBytesRead, err := file.Read(buffer)
			if err != nil {
				log.Println(err)
			}

			reqCh <- &storage_object.UploadRequest{
				FilePath: outputPath,
				Start:    int64(cur),
				Bytes:    buffer[:noBytesRead],
			}
			noReq++
		}
		close(reqCh)
	}()

	done := make(chan bool)
	go func() {
		for {
			res, more := <-resCh
			if more {
				response := res.(*storage_object.UploadResponse)
				if response.Error != "" {
					log.Println(response.Error)
				}
			} else {
				break
			}
		}
		done <- true
	}()
	bigFruit := client.NewBigFruit()

	t := time.Now()
	err := bigFruit.Invoke("StorageObject", "Upload", options, true, reqCh, resCh)
	if err != nil {
		done <- true
		log.Println(err)
	}
	<-done
	return time.Since(t)
}
