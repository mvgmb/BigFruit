package main

import (
	"log"

	"github.com/mvgmb/BigFruit/proto/naming"
	"github.com/mvgmb/BigFruit/server"
	"github.com/mvgmb/BigFruit/util"
)

var (
	options = &util.Options{
		Host:     "localhost",
		Port:     0,
		Protocol: "tcp",
	}
)

func main() {
	shr, err := server.NewServerRequestHandler(options)
	if err != nil {
		log.Fatal(err)
	}

	storageObjectAOR := &naming.AOR{
		Host: shr.Options.Host,
		Port: uint32(shr.Options.Port),
		Id:   "StorageObject",
	}
	bindRequest := &naming.BindRequest{
		ServiceName: "StorageObject",
		Aor:         storageObjectAOR,
	}
	res, err := server.Bind(bindRequest)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)

	shr.Loop()
}
