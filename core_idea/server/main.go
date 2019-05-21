package main

import (
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
)

const (
	noListeners = 20
	size        = 3
)

func main() {
	fmt.Println("Launching server...")
	for i := 0; i < noListeners; i++ {
		go server(8081 + i)
	}
	server(8080)
}

func server(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatal(err)
	}

	netConn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer netConn.Close()

	for {
		buffer := make([]byte, math.MaxInt16)

		n, err := netConn.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}

		from, err := strconv.Atoi(string(buffer[:n]))
		if err != nil {
			log.Fatal(err)
		}

		readBuffer := make([]byte, size)

		file, err := os.Open("/home/mario/Documents/git/BigFruit/app/client/fileName")
		if err != nil {
			log.Fatal(err)
		}

		_, err = file.Seek(int64(from), 0)
		if err != nil {
			log.Fatal(err)
		}

		var bytes []byte

		noBytesRead, err := file.Read(readBuffer)
		if err != nil {
			log.Println("read", err)
			bytes = []byte("EOF")
		} else {
			bytes = readBuffer[:noBytesRead]
		}

		_, err = netConn.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		file.Close()
	}

}
