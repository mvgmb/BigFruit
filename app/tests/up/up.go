package main

import (
	"log"

	"github.com/mvgmb/BigFruit/client"
	pb "github.com/mvgmb/BigFruit/proto"
)

func main() {
	reqs := make([]*pb.StorageObjectUploadRequest, 2)
	reqs[0] = &pb.StorageObjectUploadRequest{
		FilePath: "/home/mario/Documents/git/BigFruit/data1.txt",
		Start:    0,
		Bytes:    []byte("cacholo1"),
	}
	reqs[1] = &pb.StorageObjectUploadRequest{
		FilePath: "/home/mario/Documents/git/BigFruit/data2.txt",
		Start:    0,
		Bytes:    []byte("cacholo2"),
	}
	req := make(chan *pb.StorageObjectUploadRequest)
	res := make(chan *pb.StorageObjectUploadResponse)

	sop := client.StorageObjectProxy{}

	done := make(chan bool)

	go func() {
		for {
			response, more := <-res
			if more {
				log.Println(response)
			} else {
				done <- true
				return
			}
		}
	}()
	go func() {
		req <- reqs[0]
		req <- reqs[1]
		close(req)
	}()

	err := sop.Upload(req, res)
	if err != nil {
		log.Println(err)
	}
	<-done
}
