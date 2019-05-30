package main

import (
	"log"

	"github.com/mvgmb/BigFruit/client"
)

func main() {
	requestor, err := client.NewRequestor()
	if err != nil {
		log.Fatal(err)
	}

	opt, err := requestor.Lookup("Object")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(opt)

	// options := &util.Options{
	// 	Host:     "localhost",
	// 	Port:     8080,
	// 	Protocol: "tcp",
	// }

	// // req := util.NewMessage(200, "OK", "Object.Download", []byte("/home/mario/Documents/git/BigFruit/data.txt"), []byte("0"), []byte("10"))
	// req := util.NewMessage(200, "OK", "Object.Upload", []byte("/home/mario/Documents/git/BigFruit/data.txt"), []byte("0"), []byte("10"))

	// requestor.Open(options)
	// res, err := requestor.Invoke(&req, options)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// requestor.Close(options)

	// fmt.Println(res)
}
