package client

// import (
// 	"container/list"
// 	"fmt"
// 	"log"
// 	"os"
// 	"sync"
// 	"time"

// 	"github.com/mvgmb/BigFruit/client"
// 	"github.com/mvgmb/BigFruit/util"

// 	pb "github.com/mvgmb/BigFruit/proto"
// )

// const (
// 	noThreads       = 1
// 	noBytes   int64 = 6
// )

// type Client struct {
// 	requestor *client.Requestor

// 	done      chan []byte
// 	errors    chan error
// 	listElem  chan *list.Element
// 	nextElem  chan *list.Element
// 	listMutex *sync.Mutex

// 	requestsList              *list.List
// 	fileName, outputDirectory string

// 	fileSize int64
// 	options  [noThreads]util.Options
// }

// func NewBigFruitClient() *Client {
// 	return &Client{}
// }

// func (e *Client) DownloadBigFile(fileName, outputDirectory string) error {
// 	e.fileName = fileName
// 	e.outputDirectory = outputDirectory

// 	e.done = make(chan []byte)
// 	e.errors = make(chan error)
// 	e.nextElem = make(chan *list.Element)
// 	e.listElem = make(chan *list.Element)
// 	e.listMutex = &sync.Mutex{}

// 	var err error
// 	e.requestor, err = client.NewRequestor()
// 	if err != nil {
// 		return err
// 	}

// 	// test
// 	e.options[0] = util.Options{
// 		Host:     "localhost",
// 		Port:     8080,
// 		Protocol: "tcp",
// 	}

// 	e.fileSize = int64(60)

// 	// TODO get server ports and file size from "naming" service

// 	go e.requestsListManager()

// 	for i := 0; i < noThreads; i++ {
// 		go e.smallFruitWorker(i)
// 	}

// 	file, err := os.Create(e.fileName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	for e.requestsList.Front() != nil { // it != nil {
// 		err = <-e.errors

// 		if err != nil {
// 			log.Println("d", err)
// 			continue
// 		}

// 		bytes := <-e.done
// 		it := <-e.listElem

// 		_, err := file.WriteAt(bytes, it.Value.(int64))
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		e.listMutex.Lock()
// 		e.requestsList.Remove(it)
// 		e.listMutex.Unlock()
// 	}

// 	fmt.Println("Done!")
// 	return nil
// }

// func (e *Client) requestsListManager() {
// 	e.requestsList = list.New()
// 	for i := int64(0); i < e.fileSize; i += noBytes {
// 		e.requestsList.PushBack(i)
// 	}

// 	it := e.requestsList.Front()

// 	for it != nil {
// 		cur := it

// 		e.listMutex.Lock()
// 		it = it.Next()
// 		if it == nil {
// 			it = e.requestsList.Front()
// 		}
// 		e.listMutex.Unlock()

// 		e.nextElem <- cur
// 		<-e.nextElem
// 	}
// }

// func (e *Client) smallFruitWorker(index int) {
// 	err := e.requestor.Open(&e.options[index])
// 	if err != nil {
// 		e.errors <- err
// 		return
// 	}
// 	defer e.requestor.Close(&e.options[index])

// 	err = e.openFileRequest(index)
// 	if err != nil {
// 		e.errors <- err
// 	}

// 	curElem := <-e.nextElem
// 	e.nextElem <- curElem

// 	for {
// 		time.Sleep(time.Second)

// 		err = e.sendBytesRequest(index, curElem)
// 		if err != nil {
// 			e.errors <- err
// 		}

// 		curElem = <-e.nextElem
// 		e.nextElem <- curElem
// 	}
// }

// func (e *Client) openFileRequest(index int) error {
// 	req := util.NewMessage([]byte(e.fileName), "OpenFile", "OK", 200)

// 	res, err := e.requestor.Invoke(&req, &e.options[index])
// 	if err != nil {
// 		return err
// 	}
// 	message, ok := res.(*pb.Message)
// 	if !ok {
// 		return fmt.Errorf("Not a Message")
// 	}

// 	if message.Status.Code != 200 {
// 		return fmt.Errorf(message.Status.Message)
// 	}
// 	return nil
// }

// func (e *Client) sendBytesRequest(index int, curElem *list.Element) error {
// 	args := fmt.Sprintf("%d,%d", curElem.Value.(int64), noBytes)

// 	req := util.NewMessage([]byte(args), "SendBytes", "OK", 200)

// 	res, err := e.requestor.Invoke(&req, &e.options[index])
// 	if err != nil {
// 		return err
// 	}
// 	message, ok := res.(*pb.Message)
// 	if !ok {
// 		return fmt.Errorf("Not a Message")
// 	}

// 	if message.Status.Code == 100 {
// 		e.errors <- nil
// 		e.done <- message.RawData
// 		e.listElem <- curElem
// 	} else {
// 		return fmt.Errorf(message.Status.Message)
// 	}

// 	return nil
// }
