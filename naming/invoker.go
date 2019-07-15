package naming

import (
	"log"

	"github.com/golang/protobuf/proto"
	naming "github.com/mvgmb/BigFruit/proto/naming"
	"github.com/mvgmb/BigFruit/util"
)

type Invoker struct {
	serverRequestHandler *serverRequestHandler
}

func NewInvoker(options *util.Options) (*Invoker, error) {
	rh, err := newServerRequestHandler(*options)
	if err != nil {
		return nil, err
	}
	e := &Invoker{serverRequestHandler: rh}

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

		bytes, err := e.serverRequestHandler.receive()
		if err != nil {
			// TODO handle errors
			log.Println(err)
			break
		}

		var res proto.Message

		req, objectType, requestType, err := util.UnwrapMessage(bytes)
		if err != nil {
			log.Println(err)
		}

		switch objectType {
		case "*naming":
			switch requestType {
			case "BindRequest":
				res = bind(req.(*naming.BindRequest))
			case "LookupRequest":
				res = lookup(req.(*naming.LookupRequest))
			case "LookupManyRequest":
				res = lookupMany(req.(*naming.LookupManyRequest))
			case "LookupAllRequest":
				res = lookupAll(req.(*naming.LookupAllRequest))
			default:
				res = util.ErrBadRequest
			}
		default:
			res = util.ErrBadRequest
		}

		bytes, err = util.WrapMessage(res)
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

		err = e.serverRequestHandler.close()
		if err != nil {
			log.Println(err)
		}
	}
}
