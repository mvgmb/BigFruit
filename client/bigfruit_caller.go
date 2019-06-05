package client

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mvgmb/BigFruit/util"
)

const maxNoConcurrentRequestsPerServer = 3

func Call(reqCh chan proto.Message, resCh chan proto.Message, options []*util.Options, replicate bool) error {
	// Initialize requestors
	var requestors []*Requestor
	for i := range options {
		requestor, err := NewRequestor(options[i])
		if err != nil {
			continue
		}
		requestors = append(requestors, requestor)
	}
	if len(requestors) == 0 {
		return fmt.Errorf("Unable to connect to any given servers")
	}
	defer func() {
		for i := range requestors {
			requestors[i].Close()
		}
	}()

	// Initialize internal channels
	var internal []chan proto.Message
	if replicate {
		internal = make([]chan proto.Message, len(requestors))
	} else {
		internal = make([]chan proto.Message, maxNoConcurrentRequestsPerServer*len(requestors))
	}

	errors := make([]chan error, len(internal))
	for i := 0; i < len(internal); i++ {
		internal[i] = make(chan proto.Message)
		errors[i] = make(chan error)
	}

	closed := false
	permission := make(chan bool)

	go func() {
		curID := 0
		requestorsRobin := 0
		noRoutinesPerReq := 1
		if replicate {
			noRoutinesPerReq = len(requestors)
		}

		for {
			request, more := <-reqCh
			if more {
				for i := 0; i < noRoutinesPerReq; i++ {
					<-permission

					go func(id, requestorIndex int) {
						index := id % len(internal)
						res, err := requestors[requestorIndex].Invoke(request)
						errors[index] <- err
						internal[index] <- res
					}(curID, requestorsRobin)

					curID++
					requestorsRobin++
					if requestorsRobin >= len(requestors) {
						requestorsRobin = 0
					}
				}
			} else {
				closed = true
				break
			}
		}
	}()

	// Initize as many go routines as possible
	for i := 0; i < len(internal); i++ {
		permission <- true
	}

	curID := 0
	// The idea is to consume and then initialize a new go routine
	for {
		i := curID % len(internal)

		err := <-errors[i]
		if err != nil {
			return err
		}

		res := <-internal[i]
		resCh <- res

		if closed {
			close(resCh)
			break
		}
		permission <- true
		curID++
	}

	return nil
}
