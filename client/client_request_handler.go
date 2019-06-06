package client

import (
	"fmt"
	"log"
	"math"
	"net"

	"github.com/mvgmb/BigFruit/util"
)

type clientRequestHandler struct {
	netConn net.Conn
}

func newClientRequestHandler(options *util.Options) (*clientRequestHandler, error) {
	netConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", options.Host, options.Port))
	if err != nil {
		return nil, err
	}
	e := &clientRequestHandler{
		netConn: netConn,
	}
	log.Println("connected to", fmt.Sprintf("%s:%d", options.Host, options.Port))
	return e, nil
}

func (e *clientRequestHandler) close() error {
	err := e.netConn.Close()
	if err != nil {
		return err
	}
	log.Println("disconnected from", e.netConn.RemoteAddr())
	return nil
}

func (e *clientRequestHandler) send(message *[]byte) error {
	_, err := e.netConn.Write(*message)
	return err
}

func (e *clientRequestHandler) receive() ([]byte, error) {
	buffer := make([]byte, math.MaxInt16)

	n, err := e.netConn.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}
