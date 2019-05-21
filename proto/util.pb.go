// Code generated by protoc-gen-go.
// source: util.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	util.proto

It has these top-level messages:
	Status
	Message
*/
package proto

import proto1 "github.com/golang/protobuf/proto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal

type Status struct {
	Code    uint64 `protobuf:"varint,1,opt,name=code" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto1.CompactTextString(m) }
func (*Status) ProtoMessage()    {}

type Message struct {
	Status  *Status `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	Key     string  `protobuf:"bytes,2,opt,name=key" json:"key,omitempty"`
	RawData []byte  `protobuf:"bytes,3,opt,name=raw_data,proto3" json:"raw_data,omitempty"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto1.CompactTextString(m) }
func (*Message) ProtoMessage()    {}

func (m *Message) GetStatus() *Status {
	if m != nil {
		return m.Status
	}
	return nil
}

func init() {
}
