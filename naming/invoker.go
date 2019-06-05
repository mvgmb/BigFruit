package naming

import (
	"log"

	"github.com/golang/protobuf/proto"
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

	marsh := util.NewMarshaller()

	e := &Invoker{
		serverRequestHandler: rh,
		marshaller:           marsh,
	}
	return e, nil
}

func (e *Invoker) Invoke() {
	log.Printf("Listening at %s:%d\n", e.serverRequestHandler.options.Host, e.serverRequestHandler.options.Port)

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
				break
			}

			var res proto.Message

			req := pb.MessageWrapper{}

			err = e.marshaller.Unmarshal(&bytes, &req)
			if err != nil {
				// TODO handle error
				log.Println(err)
				break
			}

			switch req.Type {
			// case "Lookup":
			// 	result, err := lookup(string(req.RawData[0]))
			// 	if err != nil {
			// 		res = util.ErrNotFound
			// 		break
			// 	}
			// 	res = util.NewMessage(200, "OK", req.Key, []byte(result.String()))
			// case "LookupMany":
			// 	result, err := lookupMany(string(req.RawData[0]), int(binary.BigEndian.Uint64(req.RawData[1])))
			// 	if err != nil {
			// 		res = util.ErrNotFound
			// 		break
			// 	}
			// 	resultBytes, err := json.Marshal(result)
			// 	if err != nil {
			// 		res = util.ErrUnknown
			// 		break
			// 	}
			// 	res = util.NewMessage(200, "OK", req.Key, resultBytes)
			// case "LookupAll":
			// 	result, err := lookupAll(string(req.RawData[0]))
			// 	if err != nil {
			// 		res = util.ErrNotFound
			// 		break
			// 	}
			// 	resultBytes, err := json.Marshal(result)
			// 	if err != nil {
			// 		res = util.ErrUnknown
			// 		break
			// 	}
			// 	res = util.NewMessage(200, "OK", req.Key, resultBytes)
			// case "Bind":
			// 	aor, err := util.StringToAOR(string(req.RawData[0]))
			// 	if err != nil {
			// 		res = util.ErrBadRequest
			// 		break
			// 	}
			// 	bind(aor)
			// 	res = util.NewMessage(200, "OK", req.Key, []byte(""))
			default:
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

		err = e.serverRequestHandler.close()
		if err != nil {
			log.Println(err)
		}
	}
}
