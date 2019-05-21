package util

import (
	"github.com/golang/protobuf/proto"
	pb "github.com/mvgmb/BigFruit/proto"
)

var (
	ErrUnknown            = NewMessage([]byte(""), "", "Unknown", 000)
	ErrBadRequest         = NewMessage([]byte(""), "", "Bad Request", 400)
	ErrUnauthorized       = NewMessage([]byte(""), "", "Unauthorized", 401)
	ErrForbidden          = NewMessage([]byte(""), "", "Forbidden", 403)
	ErrNotFound           = NewMessage([]byte(""), "", "Not found", 404)
	ErrMethodNotAllowed   = NewMessage([]byte(""), "", "Method not allowed", 405)
	ErrPayloadTooLarge    = NewMessage([]byte(""), "", "Payload too large!", 413)
	ErrExpectationFailed  = NewMessage([]byte(""), "", "Expectation fail", 417)
	ErrServiceUnavailable = NewMessage([]byte(""), "", "Service unavailable", 503)
)

// Options defines the options values
type Options struct {
	Host     string
	Port     uint16
	Protocol string
}

// NewMessage creates a message
func NewMessage(bytes []byte, key, statusMessage string, statusCode uint64) proto.Message {
	status := &pb.Status{
		Code:    statusCode,
		Message: statusMessage,
	}
	message := &pb.Message{
		Status:  status,
		Key:     key,
		RawData: bytes,
	}
	return message
}
