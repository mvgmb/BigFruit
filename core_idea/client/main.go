package main

import (
	"os"
	"net"
	"fmt"
	"log"
	"math"
	"strconv"
	"container/list"
)

const (   
	noThreads = 20
	size = 3
)

var (
	done [noThreads]chan []byte
	errors [noThreads]chan error
	offset = noThreads * size
)

func main() {
	list := list.New()
	
	for i := 0; i < noThreads; i++ {
		done[i] = make(chan []byte)
		errors[i] = make(chan error)
		list.PushBack(i)

		go work(i, 8080 + i)
	}
		
	file, err := os.Create("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	
	current := 0
	it := list.Front()

	for it != nil {
		cur := it
		it = it.Next()
		
		err = <-errors[cur.Value.(int)]
		
		if err != nil {
			log.Println(err)
			list.Remove(cur)
			continue
		}
		
		bytes := <-done[cur.Value.(int)]
		
		_, err := file.WriteAt(bytes, int64(current))
		if err != nil {
			log.Println(err)
		}
		
		current += size
		
		if it == nil {
			it = list.Front()
		}
	}

}

func work(index, port int) {
	netConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", "localhost", port))
	if err != nil {
		log.Fatal(err)
	}
	defer netConn.Close()

	current := index * size

	for {
		_, err = netConn.Write([]byte(strconv.Itoa(current)))
		if err != nil {
			errors[index]<- err
			break
		}
	
		buffer := make([]byte, math.MaxInt16)
		
		n, err := netConn.Read(buffer)
		if err != nil {
			errors[index]<- err
			break
		}
		
		if string(buffer[:n]) == "EOF" {
			errors[index]<- fmt.Errorf(fmt.Sprintf("%d EOF", index))
			break
		}

		errors[index] <- nil
		done[index] <- buffer[:n]

		current = current + offset
	}
}
