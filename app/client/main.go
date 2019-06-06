package main

import (
	"encoding/binary"
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
		NumberOfAor: 3,
	}
	res, err := client.LookupMany(lookupManyRequest)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)

	var options []*util.Options

	for i := range res.AorList {
		option := &util.Options{
			Host:     res.AorList[i].Host,
			Port:     uint16(res.AorList[i].Port),
			Protocol: "tcp",
		}
		options = append(options, option)
	}
	noTests := 10
	results := make([]time.Duration, noTests)

	// miniJavaSkeletonSize := int32(101449)]
	dataZipSize := int32(434402557)
	for i := 0; i < 10; i++ {
		results[i] = download("/home/mario/Desktop/data.zip", "data.zip", dataZipSize, 20000, options)
	}
	log.Println(results)

}

func download(filePath, outputPath string, noBytes, offset int32, options []*util.Options) time.Duration {
	reqCh := make(chan proto.Message)
	resCh := make(chan proto.Message)

	bf := client.NewBigFruit()

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
	t := time.Now()
	err := bf.Call("StorageObject", "Download", options, false, reqCh, resCh)
	if err != nil {
		done <- true
		log.Println(err)
	}
	<-done
	return time.Since(t)
}

func upload(options []*util.Options) {
	reqCh := make(chan proto.Message)
	resCh := make(chan proto.Message)

	reqs := make([]*storage_object.UploadRequest, 2)
	reqs[0] = &storage_object.UploadRequest{
		FilePath: "/home/mario/Documents/git/BigFruit/data1.txt",
		Start:    0,
		Bytes:    []byte("cachsssssssssssolo1"),
	}
	reqs[1] = &storage_object.UploadRequest{
		FilePath: "/home/mario/Documents/git/BigFruit/data2.txt",
		Start:    0,
		Bytes:    []byte("cacholo2"),
	}

	bf := client.NewBigFruit()

	go bf.Call("StorageObject", "Upload", options, true, reqCh, resCh)

	go func() {
		for i := range reqs {
			reqCh <- reqs[i]
		}
		close(reqCh)
	}()

	for {
		res, more := <-resCh
		if more {
			log.Println(res)
		} else {
			break
		}
	}
}
