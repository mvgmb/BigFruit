package main

import (
	"log"

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
	shr, err := server.NewServerRequestHandler(options)
	if err != nil {
		log.Fatal(err)
	}

	aor := util.AOR{
		Host: shr.Options.Host,
		Port: shr.Options.Port,
		ID:   "StorageObject",
	}
	err = server.Bind(&aor)
	if err != nil {
		log.Fatal(err)
	}

	shr.Loop()
}
