// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/centrifuge/centrifuge.proto

package com_github_romatroskin_viqchat_centrifuge_service_centrifuge

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

type ServiceMessageType int32

const (
	ServiceMessageType_Init   ServiceMessageType = 0
	ServiceMessageType_Add    ServiceMessageType = 1
	ServiceMessageType_Remove ServiceMessageType = 2
)

var ServiceMessageType_name = map[int32]string{
	0: "Init",
	1: "Add",
	2: "Remove",
}

var ServiceMessageType_value = map[string]int32{
	"Init":   0,
	"Add":    1,
	"Remove": 2,
}

func (x ServiceMessageType) String() string {
	return proto.EnumName(ServiceMessageType_name, int32(x))
}

func (ServiceMessageType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_afc5021d42f92d6e, []int{0}
}

type ServiceMessage struct {
	Type                 ServiceMessageType `protobuf:"varint,1,opt,name=type,proto3,enum=com.github.romatroskin.viqchat.centrifuge.service.centrifuge.ServiceMessageType" json:"type,omitempty"`
	Data                 string             `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ServiceMessage) Reset()         { *m = ServiceMessage{} }
func (m *ServiceMessage) String() string { return proto.CompactTextString(m) }
func (*ServiceMessage) ProtoMessage()    {}
func (*ServiceMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_afc5021d42f92d6e, []int{0}
}

func (m *ServiceMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServiceMessage.Unmarshal(m, b)
}
func (m *ServiceMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServiceMessage.Marshal(b, m, deterministic)
}
func (m *ServiceMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServiceMessage.Merge(m, src)
}
func (m *ServiceMessage) XXX_Size() int {
	return xxx_messageInfo_ServiceMessage.Size(m)
}
func (m *ServiceMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ServiceMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ServiceMessage proto.InternalMessageInfo

func (m *ServiceMessage) GetType() ServiceMessageType {
	if m != nil {
		return m.Type
	}
	return ServiceMessageType_Init
}

func (m *ServiceMessage) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func init() {
	proto.RegisterEnum("com.github.romatroskin.viqchat.centrifuge.service.centrifuge.ServiceMessageType", ServiceMessageType_name, ServiceMessageType_value)
	proto.RegisterType((*ServiceMessage)(nil), "com.github.romatroskin.viqchat.centrifuge.service.centrifuge.ServiceMessage")
}

func init() { proto.RegisterFile("proto/centrifuge/centrifuge.proto", fileDescriptor_afc5021d42f92d6e) }

var fileDescriptor_afc5021d42f92d6e = []byte{
	// 195 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2c, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0x4e, 0xcd, 0x2b, 0x29, 0xca, 0x4c, 0x2b, 0x4d, 0x4f, 0x45, 0x62, 0xea, 0x81,
	0xe5, 0x84, 0x6c, 0x92, 0xf3, 0x73, 0xf5, 0xd2, 0x33, 0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x8a, 0xf2,
	0x73, 0x13, 0x4b, 0x8a, 0xf2, 0x8b, 0xb3, 0x33, 0xf3, 0xf4, 0xca, 0x32, 0x0b, 0x93, 0x33, 0x12,
	0x4b, 0xf4, 0x90, 0x34, 0x14, 0xa7, 0x16, 0x95, 0x65, 0x26, 0xa7, 0x22, 0x09, 0x29, 0x75, 0x31,
	0x72, 0xf1, 0x05, 0x43, 0x84, 0x7d, 0x53, 0x8b, 0x8b, 0x13, 0xd3, 0x53, 0x85, 0x52, 0xb8, 0x58,
	0x4a, 0x2a, 0x0b, 0x52, 0x25, 0x18, 0x15, 0x18, 0x35, 0xf8, 0x8c, 0x02, 0xf4, 0x28, 0x31, 0x5f,
	0x0f, 0xd5, 0xec, 0x90, 0xca, 0x82, 0xd4, 0x20, 0xb0, 0xe9, 0x42, 0x42, 0x5c, 0x2c, 0x29, 0x89,
	0x25, 0x89, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x60, 0xb6, 0x96, 0x31, 0x97, 0x10, 0xa6,
	0x7a, 0x21, 0x0e, 0x2e, 0x16, 0xcf, 0xbc, 0xcc, 0x12, 0x01, 0x06, 0x21, 0x76, 0x2e, 0x66, 0xc7,
	0x94, 0x14, 0x01, 0x46, 0x21, 0x2e, 0x2e, 0xb6, 0xa0, 0xd4, 0xdc, 0xfc, 0xb2, 0x54, 0x01, 0xa6,
	0x24, 0x36, 0x70, 0x30, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x82, 0xa2, 0x7a, 0x8a, 0x2b,
	0x01, 0x00, 0x00,
}