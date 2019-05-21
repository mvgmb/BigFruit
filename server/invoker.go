package server

import (
	"log"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/mvgmb/BigFruit/bigfruit"
	pb "github.com/mvgmb/BigFruit/proto"
	"github.com/mvgmb/BigFruit/util"
)

type Invoker struct {
	serverRequestHandler *serverRequestHandler
	marshaller           *util.Marshaller
}

func NewInvoker(options *util.Options) (*Invoker, error) {
	rh, err := newServerRequestHandler(*options)
	if err != nil {
		return nil, err
	}

	marsh, err := util.NewMarshaller()
	if err != nil {
		return nil, err
	}

	e := &Invoker{
		serverRequestHandler: rh,
		marshaller:           marsh,
	}
	return e, nil
}

func (e *Invoker) Invoke() {
	log.Printf("Listening at %s:%d\n", e.serverRequestHandler.options.Host, e.serverRequestHandler.options.Port)

	bigFruit := bigfruit.NewBigFruitServer()
	bigFruit.RegisterFilePath("data.txt", "/home/mario/Documents/git/BigFruit/core_idea/server/data.txt")

	// TODO register service on "naming" service

	for {
		err := e.serverRequestHandler.accept()
		if err != nil {
			log.Println(err)
			continue
		}
		defer e.serverRequestHandler.close()

		for {
			bytes, err := e.serverRequestHandler.receive()
			if err != nil {
				// TODO handle error
				log.Println(err)
				e.serverRequestHandler.close()
				break
			}

			var res proto.Message

			req := pb.Message{}

			err = e.marshaller.Unmarshal(&bytes, &req)
			if err != nil {
				// TODO handle error
				log.Println(err)
			}

			if req.Status.Code == 200 {
				switch req.Key {
				case "OpenFile":
					err = bigFruit.OpenFile(string(req.RawData))
					if err != nil {
						res = util.ErrBadRequest
					} else {
						res = util.NewMessage([]byte("Ready to start transfer!"), string(req.RawData), "OK", 200)
					}
					break
				case "SendBytes":
					from, err := strconv.Atoi(string(req.RawData))
					if err != nil {
						res = util.ErrBadRequest
						break
					}
					bytes, err := bigFruit.SeekBytes(int64(from))
					if err != nil {
						res = util.NewMessage([]byte(err.Error()), "SeekBytesError", "Not Found", 404)
						break
					}
					res = util.NewMessage(bytes, string(req.RawData), "Continue", 100)
					break
				default:
					res = util.ErrBadRequest
				}
			} else {
				res = util.ErrBadRequest
			}

			bytes, err = e.marshaller.Marshal(&res)
			if err != nil {
				// TODO handle error
				log.Println(err)
			}

			err = e.serverRequestHandler.send(&bytes)
			if err != nil {
				// TODO handle error
				log.Println(err)
			}
		}
	}
}
