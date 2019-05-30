package main

import (
	"log"

	"github.com/mvgmb/BigFruit/naming"
	"github.com/mvgmb/BigFruit/util"
)

var (
	options = util.Options{
		Host:     "localhost",
		Port:     1337,
		Protocol: "tcp",
	}
)

func main() {
	invoker, err := naming.NewInvoker(&options)
	if err != nil {
		log.Fatal(err)
	}
	invoker.Invoke()
}
