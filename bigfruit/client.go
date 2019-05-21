package bigfruit

import (
	"container/list"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/mvgmb/BigFruit/client"
	"github.com/mvgmb/BigFruit/util"

	pb "github.com/mvgmb/BigFruit/proto"
)

const (
	noThreads = 1
	noBytes   = 3
)

type Client struct {
	requestor *client.Requestor
	options   [noThreads]util.Options

	done   [noThreads]chan []byte
	errors [noThreads]chan error
	offset int

	fileName, outputDirectory string
}

func NewBigFruitClient() (*Client, error) {
	e := &Client{
		offset: noThreads * noBytes,
	}
	return e, nil
}

func (e *Client) DownloadBigFile(_fileName, _outputDirectory string) error {
	e.fileName = _fileName
	e.outputDirectory = _outputDirectory

	var err error
	e.requestor, err = client.NewRequestor()
	if err != nil {
		return err
	}

	// test
	e.options[0] = util.Options{
		Host:     "localhost",
		Port:     8080,
		Protocol: "tcp",
	}

	// TODO get server ports from "naming" service

	list := list.New()

	for i := 0; i < noThreads; i++ {
		list.PushBack(i)
		e.done[i] = make(chan []byte)
		e.errors[i] = make(chan error)

		go e.work(i)
	}

	file, err := os.Create(e.fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	current := 0
	it := list.Front()

	for it != nil {
		curWorker := it
		it = it.Next()

		err = <-e.errors[curWorker.Value.(int)]

		if err != nil {
			log.Println(err)
			list.Remove(curWorker)
			continue
		}

		bytes := <-e.done[curWorker.Value.(int)]

		log.Println(string(bytes))

		_, err := file.WriteAt(bytes, int64(current))
		if err != nil {
			log.Println(err)
		}

		current += noBytes

		if it == nil {
			it = list.Front()
		}
	}

	return nil
}

func (e *Client) work(index int) {
	err := e.requestor.Open(&e.options[index])
	if err != nil {
		e.errors[index] <- err
		return
	}
	defer e.requestor.Close(&e.options[index])

	req := util.NewMessage([]byte(e.fileName), "OpenFile", "OK", 200)

	res, err := e.requestor.Invoke(&req, &e.options[index])
	if err != nil {
		log.Println("essage")
		e.errors[index] <- err
		return
	}
	message, ok := res.(*pb.Message)
	if !ok {
		e.errors[index] <- fmt.Errorf("Not a Message")
		return
	}

	if message.Status.Code != 200 {
		e.errors[index] <- fmt.Errorf(message.Status.Message)
		return
	}

	current := index * noBytes

	for {
		req := util.NewMessage([]byte(strconv.Itoa(current)), "SendBytes", "OK", 200)

		res, err := e.requestor.Invoke(&req, &e.options[index])
		if err != nil {
			e.errors[index] <- err
			break
		}

		message, ok := res.(*pb.Message)
		if !ok {
			e.errors[index] <- fmt.Errorf("Not a Message")
			break
		}

		if message.Status.Code == 100 {
			e.errors[index] <- nil
			e.done[index] <- message.RawData
		} else {
			e.errors[index] <- fmt.Errorf(message.Status.Message)
			break
		}

		current = current + e.offset
	}
}
