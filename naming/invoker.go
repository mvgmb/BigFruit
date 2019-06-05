package naming

import (
	"log"

	"github.com/golang/protobuf/proto"
	pb "github.com/mvgmb/BigFruit/proto"
	naming "github.com/mvgmb/BigFruit/proto/naming"
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
	defer log.Println("die")
	for {
		err := e.serverRequestHandler.accept()
		if err != nil {
			log.Println(err)
			e.serverRequestHandler.close()
			continue
		}

		bytes, err := e.serverRequestHandler.receive()
		if err != nil {
			// TODO handle error
			log.Println(err)
			break
		}

		var res proto.Message

		wrapper := pb.MessageWrapper{}

		err = e.marshaller.Unmarshal(&bytes, &wrapper)
		if err != nil {
			// TODO handle error
			log.Println(err)
			break
		}

		req, err := util.UnwrapMessage(&wrapper)
		if err != nil {
			log.Println(err)
		}

		switch wrapper.Type {
		case "*naming.BindRequest":
			res = bind(req.(*naming.BindRequest))
		case "*naming.LookupRequest":
			res = lookup(req.(*naming.LookupRequest))
		case "*naming.LookupManyRequest":
			res = lookupMany(req.(*naming.LookupManyRequest))
		case "*naming.LookupAllRequest":
			res = lookupAll(req.(*naming.LookupAllRequest))
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
