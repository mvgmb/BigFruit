package main

import (
	"log"

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

// func upload() {
// 	reqs := make([]*storage_object.UploadRequest, 2)
// 	reqs[0] = &storage_object.UploadRequest{
// 		FilePath: "/home/mario/Documents/git/BigFruit/data1.txt",
// 		Start:    0,
// 		Bytes:    []byte("cacholo1"),
// 	}
// 	reqs[1] = &storage_object.UploadRequest{
// 		FilePath: "/home/mario/Documents/git/BigFruit/data2.txt",
// 		Start:    0,
// 		Bytes:    []byte("cacholo2"),
// 	}
// 	req := make(chan *storage_object.UploadRequest)
// 	res := make(chan *storage_object.UploadResponse)

// 	sop := client.StorageObjectProxy{}

// 	done := make(chan bool)

// 	go func() {
// 		for {
// 			response, more := <-res
// 			if more {
// 				log.Println(response)
// 			} else {
// 				done <- true
// 				return
// 			}
// 		}
// 	}()
// 	go func() {
// 		req <- reqs[0]
// 		req <- reqs[1]
// 		close(req)
// 	}()

// 	err := sop.Upload(req, res)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	<-done
// }

// func download() {
// 	req := make(chan *storage_object.DownloadRequest)
// 	res := make(chan *storage_object.DownloadResponse)

// 	sop := client.StorageObjectProxy{}

// 	done := make(chan bool)

// 	go func() {
// 		for {
// 			response, more := <-res
// 			if more {
// 				log.Println(response)
// 			} else {
// 				done <- true
// 				return
// 			}
// 		}
// 	}()
// 	go func() {
// 		request := &storage_object.DownloadRequest{
// 			FilePath: "/home/mario/Documents/git/BigFruit/MiniJavaSkeleton.zip",
// 			Start:    0,
// 			Offset:   100,
// 		}
// 		req <- request
// 		close(req)
// 	}()

// 	err := sop.Download(req, res)
// 	if err != nil {
// 		log.Println(err)
// 		done <- true
// 	}
// 	<-done
// }
