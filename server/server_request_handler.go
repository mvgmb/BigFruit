package server

import (
	"fmt"
	"log"
	"math"
	"net"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	pb "github.com/mvgmb/BigFruit/proto"
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
	marshaller := util.NewMarshaller()
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

			go func() {
				var res proto.Message

				bytes := buffer[:n]

				wrapper := &pb.MessageWrapper{}

				err = marshaller.Unmarshal(&bytes, wrapper)
				if err != nil {
					// TODO handle error
					log.Println(err)
					return
				}

				req, err := util.UnwrapMessage(wrapper)
				if err != nil {
					log.Println(err)
				}

				if strings.HasPrefix(wrapper.Type, "*proto.StorageObject") {
					res, err = storageInvoker.Invoke(wrapper.Type, req)
					if err != nil {
						res = util.ErrBadRequest
					}
				} else {
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
			}()
		}
		netConn.Close()
	}
}
