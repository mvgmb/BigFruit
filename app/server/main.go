package main

import (
	"log"

	"github.com/mvgmb/BigFruit/server"
	"github.com/mvgmb/BigFruit/util"
)

var (
	options = util.Options{
		Host:     "localhost",
		Port:     0,
		Protocol: "tcp",
	}
)

func main() {
	invoker, err := server.NewInvoker(&options)
	if err != nil {
		log.Fatal(err)
	}
	invoker.Invoke()
}
