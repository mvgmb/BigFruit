package main

import (
	"log"

	pb "github.com/mvgmb/BigFruit/proto"
	"github.com/mvgmb/BigFruit/server"
	"github.com/mvgmb/BigFruit/util"
)

var (
	options = &util.Options{
		Host:     "localhost",
		Port:     8080,
		Protocol: "tcp",
	}
)

func main() {
	storageObjectAOR := &pb.AOR{
		Host: "localhost",
		Port: 8080,
		Id:   "StorageObject",
	}
	bindRequest := &pb.NamingServiceBindRequest{
		ServiceName: "StorageObject",
		Aor:         storageObjectAOR,
	}
	res, err := server.Bind(bindRequest)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)

	shr, err := server.NewServerRequestHandler(options)
	if err != nil {
		log.Fatal(err)
	}
	shr.Loop()
}
