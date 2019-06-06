package client

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/mvgmb/BigFruit/app/proto/storage_object"
	"github.com/mvgmb/BigFruit/util"
)

type BigFruit struct {
	proxy interface{}
}

const maxNoConcurrentRequestsPerServer = 1

func NewBigFruit() *BigFruit {
	return &BigFruit{}
}

func (e *BigFruit) Call(objectName, methodName string, options []*util.Options, replicate bool, reqCh, resCh chan proto.Message) error {
	// Initialize proxy
	switch objectName {
	case "StorageObject":
		proxy, err := NewStorageObjectProxy(options)
		if err != nil {
			return err
		}
		defer proxy.Close()
		e.proxy = proxy
	default:
		return fmt.Errorf("Object requested not found")
	}

	// Initialize internal channels
	var internal []chan proto.Message
	if replicate {
		internal = make([]chan proto.Message, len(options))
	} else {
		internal = make([]chan proto.Message, maxNoConcurrentRequestsPerServer*len(options))
	}

	errors := make([]chan error, len(internal))
	for i := 0; i < len(internal); i++ {
		internal[i] = make(chan proto.Message)
		errors[i] = make(chan error)
	}

	lastRequestID := -1
	requestsCurID := 0
	permission := make(chan bool)
	waiting := make(chan bool)

	go func() {
		noRoutinesPerReq := 1
		if replicate {
			noRoutinesPerReq = len(options)
		}
		for {
			request, more := <-reqCh
			if more {
				for i := 0; i < noRoutinesPerReq; i++ {
					waiting <- true
					<-permission

					go func(id int, request proto.Message) {
						index := id % len(internal)
						protoMessage, err := callObject(objectName, methodName, e.proxy, request)
						if err != nil {
							log.Println(err)
						}
						errors[index] <- err
						internal[index] <- protoMessage
					}(requestsCurID, request)

					requestsCurID++
				}
			} else {
				close(waiting)
				lastRequestID = requestsCurID - 1
				break
			}
		}
	}()

	more := true
	// Initize as many go routines as possible
	for i := 0; i < len(internal) && more; i++ {
		_, more = <-waiting
		permission <- true
	}

	responsesCurID := 0

	// The idea is to consume and then initialize a new go routine
	for {
		i := responsesCurID % len(internal)

		err := <-errors[i]
		if err != nil {
			return err
		}
		resCh <- <-internal[i]

		if responsesCurID == lastRequestID {
			close(resCh)
			break
		}

		_, more = <-waiting
		if more {
			permission <- true
		}

		responsesCurID++
	}
	return nil
}

func callObject(objectName, methodName string, proxy interface{}, req proto.Message) (proto.Message, error) {
	switch objectName {
	case "StorageObject":
		return callStorageObjectMethod(methodName, proxy.(*StorageObjectProxy), req)
	default:
		return nil, fmt.Errorf("Object not found")
	}
}

func callStorageObjectMethod(methodName string, proxy *StorageObjectProxy, req proto.Message) (proto.Message, error) {
	switch methodName {
	case "Upload":
		return proxy.Upload(req.(*storage_object.UploadRequest))
	case "Download":
		return proxy.Download(req.(*storage_object.DownloadRequest))
	default:
		return nil, fmt.Errorf("Method not found")
	}
}
