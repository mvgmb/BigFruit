package client

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/mvgmb/BigFruit/app/proto/storage_object"
	"github.com/mvgmb/BigFruit/util"
)

type BigFruit struct {
	proxies interface{}
}

func NewBigFruit() *BigFruit {
	return &BigFruit{}
}

func (e *BigFruit) Call(objectName, methodName string, options []*util.Options, replicate bool, reqCh, resCh chan proto.Message) error {
	// Initialize proxies
	switch objectName {
	case "StorageObject":
		var proxies []*StorageObjectProxy
		for i := range options {
			proxy, err := NewStorageObjectProxy(options[i])
			if err != nil {
				return err
			}
			defer proxy.Close()
			proxies = append(proxies, proxy)
		}
		e.proxies = proxies
	default:
		return fmt.Errorf("Object requested not found")
	}

	// Initialize internal channels
	internal := make([]chan proto.Message, len(options))

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
		noRequestsPerRoutine := 1
		if replicate {
			noRequestsPerRoutine = len(options)
		}
		for {
			request, more := <-reqCh
			if more {
				for i := 0; i < noRequestsPerRoutine; i++ {
					index := requestsCurID % len(options)

					waiting <- true
					<-permission

					go func() {
						protoMessage, err := e.callObject(index, objectName, methodName, request)
						if err != nil {
							log.Println(err)
						}
						errors[index] <- err
						internal[index] <- protoMessage
					}()

					requestsCurID++
				}
			} else {
				close(waiting)
				lastRequestID = requestsCurID
				break
			}
		}
	}()

	more := true
	// Initize as many go routines as possible
	for i := 0; i < len(options) && more; i++ {
		_, more = <-waiting
		permission <- true
	}

	// The idea is to consume and then initialize a new go routine
	responsesCurID := 0
	for {
		index := responsesCurID % len(internal)

		err := <-errors[index]
		if err != nil {
			return err
		}
		resCh <- <-internal[index]

		responsesCurID++
		if responsesCurID == lastRequestID {
			close(resCh)
			break
		}

		_, more = <-waiting
		if more {
			permission <- true
		}
	}
	return nil
}

func (e *BigFruit) callObject(requestorIndex int, objectName, methodName string, req proto.Message) (proto.Message, error) {
	switch objectName {
	case "StorageObject":
		return e.callStorageObjectMethod(requestorIndex, methodName, req)
	default:
		return nil, fmt.Errorf("Object not found")
	}
}

func (e *BigFruit) callStorageObjectMethod(requestorIndex int, methodName string, req proto.Message) (proto.Message, error) {
	switch methodName {
	case "Upload":
		return e.proxies.([]*StorageObjectProxy)[requestorIndex].Upload(req.(*storage_object.UploadRequest))
	case "Download":
		return e.proxies.([]*StorageObjectProxy)[requestorIndex].Download(req.(*storage_object.DownloadRequest))
	default:
		return nil, fmt.Errorf("Method not found")
	}
}
