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
	invoker, err := server.NewInvoker(options)
	if err != nil {
		log.Fatal(err)
	}

	aor := util.AOR{
		Host: invoker.Options.Host,
		Port: invoker.Options.Port,
		ID:   "Object",
	}
	err = server.Bind(&aor)
	if err != nil {
		log.Fatal(err)
	}

	invoker.Invoke()
}
