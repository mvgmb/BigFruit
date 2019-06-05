package main

import (
	"log"

	"github.com/mvgmb/BigFruit/client"
	pb "github.com/mvgmb/BigFruit/proto"
)

func main() {
	req := make(chan *pb.StorageObjectDownloadRequest)
	res := make(chan *pb.StorageObjectDownloadResponse)

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
		request := &pb.StorageObjectDownloadRequest{
			FilePath: "/home/mario/Documents/git/BigFruit/MiniJavaSkeleton.zip",
			Start:    0,
			Offset:   100,
		}
		req <- request
		close(req)
	}()

	err := sop.Download(req, res)
	if err != nil {
		log.Println(err)
		done <- true
	}
	<-done
}
