// Code generated by protoc-gen-go. DO NOT EDIT.
// source: broker.proto

package tars_broker

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// broker event
type Event struct {
	//the event name  e.g login
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// event headers
	Header map[string][]byte `protobuf:"bytes,4,rep,name=header,proto3" json:"header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// the event data  e.g proto.marshal data
	Data                 []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_f209535e190f2bed, []int{0}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Event) GetHeader() map[string][]byte {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Event) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*Event)(nil), "tars.broker.Event")
	proto.RegisterMapType((map[string][]byte)(nil), "tars.broker.Event.HeaderEntry")
}

func init() { proto.RegisterFile("broker.proto", fileDescriptor_f209535e190f2bed) }

var fileDescriptor_f209535e190f2bed = []byte{
	// 164 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0x2a, 0xca, 0xcf,
	0x4e, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2e, 0x49, 0x2c, 0x2a, 0xd6, 0x83,
	0x08, 0x29, 0x2d, 0x62, 0xe4, 0x62, 0x75, 0x2d, 0x4b, 0xcd, 0x2b, 0x11, 0x12, 0xe2, 0x62, 0xc9,
	0x4b, 0xcc, 0x4d, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb3, 0x85, 0xcc, 0xb8, 0xd8,
	0x32, 0x52, 0x13, 0x53, 0x52, 0x8b, 0x24, 0x58, 0x14, 0x98, 0x35, 0xb8, 0x8d, 0xe4, 0xf4, 0x90,
	0xf4, 0xea, 0x81, 0xf5, 0xe9, 0x79, 0x80, 0x15, 0xb8, 0xe6, 0x95, 0x14, 0x55, 0x06, 0x41, 0x55,
	0x83, 0xcc, 0x4a, 0x49, 0x2c, 0x49, 0x94, 0x60, 0x52, 0x60, 0xd4, 0xe0, 0x09, 0x02, 0xb3, 0xa5,
	0x2c, 0xb9, 0xb8, 0x91, 0x94, 0x0a, 0x09, 0x70, 0x31, 0x67, 0xa7, 0x56, 0x42, 0x6d, 0x03, 0x31,
	0x85, 0x44, 0xb8, 0x58, 0xcb, 0x12, 0x73, 0x4a, 0x53, 0xa1, 0xba, 0x20, 0x1c, 0x2b, 0x26, 0x0b,
	0xc6, 0x24, 0x36, 0xb0, 0xc3, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x3e, 0xd1, 0xff, 0xde,
	0xc8, 0x00, 0x00, 0x00,
}
