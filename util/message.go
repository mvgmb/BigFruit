package util

import (
	"reflect"

	"github.com/golang/protobuf/proto"
	pb "github.com/mvgmb/BigFruit/proto"
)

var (
	// ErrUnknown            = NewMessage(000, "Unknown", "")
	ErrBadRequest = &pb.Error{Code: 400, Message: "Bad Request"}
	// ErrUnauthorized       = NewMessage(401, "Unauthorized", "")
	// ErrForbidden          = NewMessage(403, "Forbidden", "")
	// ErrNotFound           = NewMessage(404, "Not found", "")
	// ErrMethodNotAllowed   = NewMessage(405, "Method not allowed", "")
	// ErrPayloadTooLarge    = NewMessage(413, "Payload too large!", "")
	// ErrExpectationFailed  = NewMessage(417, "Expectation fail", "")
	// ErrServiceUnavailable = NewMessage(503, "Service unavailable", "")
)

type Options struct {
	Host     string
	Port     uint16
	Protocol string
}

var protoType = map[string]proto.Message{
	"*proto.Error": &pb.Error{},

	"*proto.StorageObjectUploadRequest":    &pb.StorageObjectUploadRequest{},
	"*proto.StorageObjectUploadResponse":   &pb.StorageObjectUploadResponse{},
	"*proto.StorageObjectDownloadRequest":  &pb.StorageObjectDownloadRequest{},
	"*proto.StorageObjectDownloadResponse": &pb.StorageObjectDownloadResponse{},

	"*proto.NamingServiceBindRequest":        &pb.NamingServiceBindRequest{},
	"*proto.NamingServiceBindResponse":       &pb.NamingServiceBindResponse{},
	"*proto.NamingServiceLookupRequest":      &pb.NamingServiceLookupRequest{},
	"*proto.NamingServiceLookupResponse":     &pb.NamingServiceLookupResponse{},
	"*proto.NamingServiceLookupManyRequest":  &pb.NamingServiceLookupManyRequest{},
	"*proto.NamingServiceLookupManyResponse": &pb.NamingServiceLookupManyResponse{},
	"*proto.NamingServiceLookupAllRequest":   &pb.NamingServiceLookupAllRequest{},
	"*proto.NamingServiceLookupAllResponse":  &pb.NamingServiceLookupAllResponse{},
}

func WrapMessage(message proto.Message) ([]byte, error) {
	bytes, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}

	wrapper := &pb.MessageWrapper{
		Type:    reflect.TypeOf(message).String(),
		Message: bytes,
	}

	bytes, err = proto.Marshal(wrapper)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func UnwrapMessage(wrapper *pb.MessageWrapper) (proto.Message, error) {
	message := protoType[wrapper.Type]
	err := proto.Unmarshal(wrapper.Message, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}
