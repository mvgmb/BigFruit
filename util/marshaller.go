package util

import (
	"github.com/golang/protobuf/proto"
)

// Marshaller struct
type Marshaller struct {
}

// NewMarshaller constructs a new Mashaller
func NewMarshaller() (*Marshaller, error) {
	return &Marshaller{}, nil
}

// Marshal serializes the message into bytes
func (e *Marshaller) Marshal(message *proto.Message) ([]byte, error) {
	return proto.Marshal(*message)
}

// Unmarshal retrieves the serialized message
func (e *Marshaller) Unmarshal(bytes *[]byte, pb proto.Message) error {
	return proto.Unmarshal(*bytes, pb)
}
