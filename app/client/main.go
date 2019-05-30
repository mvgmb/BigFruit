package main

import (
	"fmt"
	"log"

	"github.com/mvgmb/BigFruit/client"
	"github.com/mvgmb/BigFruit/util"
)

func main() {
	requestor, err := client.NewRequestor()
	if err != nil {
		log.Fatal(err)
	}

	options, err := requestor.Lookup("StorageObject")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(options)

	// options := &util.Options{
	// 	Host:     "localhost",
	// 	Port:     8080,
	// 	Protocol: "tcp",
	// }

	req := util.NewMessage(200, "OK", "StorageObject.Download", []byte("/home/mario/Documents/git/BigFruit/data.txt"), []byte("0"), []byte("100"))
	// req := util.NewMessage(200, "OK", "StorageObject.Upload", []byte("/home/mario/Documents/git/BigFruit/data.txt"), []byte("0"), []byte("cafeeee"))

	requestor.Open(options)
	res, err := requestor.Invoke(&req, options)
	if err != nil {
		log.Fatal(err)
	}
	requestor.Close()

	fmt.Println(res)
}
