package naming

import (
	"fmt"
	"math"
	"net"

	"github.com/mvgmb/BigFruit/util"
)

type serverRequestHandler struct {
	options  util.Options
	netConn  net.Conn
	listener net.Listener
}

func newServerRequestHandler(options util.Options) (*serverRequestHandler, error) {
	addr := fmt.Sprintf("%s:%d", options.Host, options.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	e := &serverRequestHandler{
		options:  options,
		listener: listener,
	}

	return e, nil
}

func (e *serverRequestHandler) accept() error {
	if e.netConn != nil {
		return fmt.Errorf("Connection Already Accepted")
	}

	newConn, err := e.listener.Accept()
	if err != nil {
		return err
	}

	e.netConn = newConn
	return nil
}

func (e *serverRequestHandler) open(options *util.Options) error {
	if e.netConn != nil {
		return fmt.Errorf("Connection Already Connected")
	}

	netConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", options.Host, options.Port))
	if err != nil {
		return err
	}

	e.netConn = netConn
	return nil
}

func (e *serverRequestHandler) close() error {
	if e.netConn == nil {
		return fmt.Errorf("Connection Already Closed")
	}
	netConn := e.netConn
	e.netConn = nil

	err := netConn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (e *serverRequestHandler) send(message *[]byte) error {
	if e.netConn == nil {
		return fmt.Errorf("Connection Not Accepted")
	}

	_, err := e.netConn.Write(*message)
	return err
}

func (e *serverRequestHandler) receive() ([]byte, error) {
	if e.netConn == nil {
		return nil, fmt.Errorf("Connection Not Accepted")
	}

	buffer := make([]byte, math.MaxInt16)

	n, err := e.netConn.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}
