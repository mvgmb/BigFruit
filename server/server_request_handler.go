package server

import (
	"fmt"
	"log"
	"math"
	"net"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/mvgmb/BigFruit/util"
)

type ServerRequestHandler struct {
	Options  *util.Options
	listener *net.Listener
}

func NewServerRequestHandler(options *util.Options) (*ServerRequestHandler, error) {
	addr := fmt.Sprintf("%s:%d", options.Host, options.Port)

	if options.Port == 0 {
		addr = fmt.Sprintf("%s:", options.Host)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	resultAddr := strings.Split(listener.Addr().String(), ":")
	options.Host = resultAddr[0]

	num, err := strconv.ParseUint(resultAddr[1], 10, 16)
	if err != nil {
		return nil, err
	}

	options.Port = uint16(num)

	e := &ServerRequestHandler{
		Options:  options,
		listener: &listener,
	}

	return e, nil
}

func (e *ServerRequestHandler) Loop() {
	storageInvoker := NewStorageObjectInvoker()

	log.Printf("Listening at %s:%d", e.Options.Host, e.Options.Port)

	for {
		netConn, err := (*e.listener).Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		for {
			buffer := make([]byte, math.MaxInt16)

			n, err := netConn.Read(buffer)
			if err != nil {
				log.Println(err)
				break
			}

			var res proto.Message

			bytes := buffer[:n]

			req, objectName, requestName, err := util.UnwrapMessage(bytes)
			if err != nil {
				log.Println(err)
			}

			switch objectName {
			case "*storage_object":
				res, err = storageInvoker.Invoke(requestName, req)
				if err != nil {
					res = util.ErrBadRequest
				}
			default:
				res = util.ErrBadRequest
			}

			bytes, err = util.WrapMessage(res)
			if err != nil {
				// TODO handle error
				log.Println(err)
				return
			}

			_, err = netConn.Write(bytes)
			if err != nil {
				// TODO handle error
				log.Println(err)
				return
			}
		}
		netConn.Close()
	}
}
