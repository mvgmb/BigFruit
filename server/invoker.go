package server

import (
	"log"
	"strconv"

	"github.com/golang/protobuf/proto"
	pb "github.com/mvgmb/BigFruit/proto"
	"github.com/mvgmb/BigFruit/util"
)

type Invoker struct {
	Options              *util.Options
	serverRequestHandler *serverRequestHandler
	marshaller           *util.Marshaller
}

func NewInvoker(options *util.Options) (*Invoker, error) {
	rh, err := newServerRequestHandler(*options)
	if err != nil {
		return nil, err
	}

	marsh := util.NewMarshaller()

	e := &Invoker{
		serverRequestHandler: rh,
		marshaller:           marsh,
		Options:              &rh.options,
	}
	return e, nil
}

func (e *Invoker) Invoke() {
	log.Printf("Listening at %s:%d\n", e.Options.Host, e.Options.Port)

	object := NewObject()

	for {
		err := e.serverRequestHandler.accept()
		if err != nil {
			log.Println(err)
			e.serverRequestHandler.close()
			continue
		}

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
				break
			}

			// Demultiplex
			if req.Status.Code == 200 {
				switch req.Key {
				case "Object.Upload":
					if len(req.RawData) != 3 {
						log.Println("Not enough arguments: needed 3, got:", len(req.RawData))
						res = util.ErrBadRequest
						break
					}
					start, err := strconv.Atoi(string(req.RawData[1]))
					if err != nil {
						log.Println(err)
						res = util.ErrBadRequest
						break
					}
					err = object.Upload(string(req.RawData[0]), int64(start), req.RawData[2])
					if err != nil {
						log.Println(err)
						res = util.ErrUnknown
						break
					}
					res = util.NewMessage(200, "OK", "Object.Upload")
					break
				case "Object.Download":
					if len(req.RawData) != 3 {
						log.Println("Not enough arguments: needed 3, got:", len(req.RawData))
						res = util.ErrBadRequest
						break
					}
					start, err := strconv.Atoi(string(req.RawData[1]))
					if err != nil {
						log.Println(err)
						res = util.ErrBadRequest
						break
					}
					offset, err := strconv.Atoi(string(req.RawData[2]))
					if err != nil {
						log.Println(err)
						res = util.ErrBadRequest
						break
					}
					bytes, err := object.Download(string(req.RawData[0]), int64(start), offset)
					if err != nil {
						log.Println(err)
						res = util.ErrBadRequest
						break
					}
					res = util.NewMessage(200, "OK", "Object.Download", bytes)
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
				break
			}

			err = e.serverRequestHandler.send(&bytes)
			if err != nil {
				// TODO handle error
				log.Println(err)
				break
			}
		}
		e.serverRequestHandler.close()
	}
}
