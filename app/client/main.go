package main

import (
	"log"

	"github.com/mvgmb/BigFruit/client"
	pb "github.com/mvgmb/BigFruit/proto"
)

func main() {
	lookupManyRequest := &pb.NamingServiceLookupManyRequest{
		ServiceName: "StorageObject",
		NumberOfAor: 1,
	}
	res, err := client.LookupMany(lookupManyRequest)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)
	// requestor, err := client.NewRequestor()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // options, err := requestor.Lookup("StorageObject")
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }
	// // log.Println(options)

	// options := &util.Options{
	// 	Host:     "localhost",
	// 	Port:     8080,
	// 	Protocol: "tcp",
	// }
	// proxy := client.NewStorageObjectProxy(requestor, options)

	// // req := util.NewMessage(200, "OK", "StorageObject.Download", []byte("/home/mario/Documents/git/BigFruit/data.txt"), []byte("0"), []byte("100"))
	// // req := util.NewMessage(200, "OK", "StorageObject.Upload", []byte("/home/mario/Documents/git/BigFruit/data.txt"), []byte("0"), []byte("cafeeee"))

	// requestor.Open(options)

	// noRequests := 1
	// uploadRequest := &pb.StorageObjectUploadRequest{
	// 	FilePath: "/home/mario/Documents/git/BigFruit/data.txt",
	// 	Start:    0,
	// 	Bytes:    []byte("cafepadas"),
	// }

	// resChan, errChan := proxy.Upload(uploadRequest)
	// for i := 0; i < noRequests; i++ {
	// 	err = <-errChan
	// 	if err != nil {
	// 		log.Println(err)
	// 	} else {
	// 		log.Println(<-resChan)
	// 	}
	// }
	// // bytes, err = proxy.Download("/home/mario/Documents/git/BigFruit/data.txt", 0, 10)
	// // if err != nil {
	// // 	log.Println(err)
	// // }

	// requestor.Close()

}
