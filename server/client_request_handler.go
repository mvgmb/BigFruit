package server

import (
	"fmt"
	"math"
	"net"

	"github.com/mvgmb/BigFruit/util"
)

type clientRequestHandler struct {
	netConn net.Conn
}

func newClientRequestHandler() *clientRequestHandler {
	e := clientRequestHandler{
		netConn: nil,
	}
	return &e
}

func (e *clientRequestHandler) open(options util.Options) error {
	if e.netConn != nil {
		return fmt.Errorf("Connection already established, please close to open a new one")
	}

	netConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", options.Host, options.Port))
	if err != nil {
		return err
	}
	e.netConn = netConn
	return nil
}

func (e *clientRequestHandler) close() error {
	if e.netConn == nil {
		return fmt.Errorf("No connection established")
	}

	err := e.netConn.Close()
	if err != nil {
		return err
	}
	e.netConn = nil
	return nil
}

func (e *clientRequestHandler) send(message *[]byte) error {
	if e.netConn == nil {
		return fmt.Errorf("No connection established")
	}
	_, err := e.netConn.Write(*message)
	return err
}

func (e *clientRequestHandler) receive() ([]byte, error) {
	if e.netConn == nil {
		return nil, fmt.Errorf("No connection established")
	}

	buffer := make([]byte, math.MaxInt16)

	n, err := e.netConn.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}
