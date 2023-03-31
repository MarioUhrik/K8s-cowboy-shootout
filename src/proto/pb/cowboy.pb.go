// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cowboy.proto

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	cowboy.proto

It has these top-level messages:
	GetShotRequest
	GetShotResponse
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type GetShotRequest struct {
	ShooterName    string `protobuf:"bytes,1,opt,name=shooterName" json:"shooterName,omitempty"`
	IncomingDamage int32  `protobuf:"varint,2,opt,name=incomingDamage" json:"incomingDamage,omitempty"`
}

func (m *GetShotRequest) Reset()                    { *m = GetShotRequest{} }
func (m *GetShotRequest) String() string            { return proto.CompactTextString(m) }
func (*GetShotRequest) ProtoMessage()               {}
func (*GetShotRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *GetShotRequest) GetShooterName() string {
	if m != nil {
		return m.ShooterName
	}
	return ""
}

func (m *GetShotRequest) GetIncomingDamage() int32 {
	if m != nil {
		return m.IncomingDamage
	}
	return 0
}

type GetShotResponse struct {
	VictimName      string `protobuf:"bytes,1,opt,name=victimName" json:"victimName,omitempty"`
	RemainingHealth int32  `protobuf:"varint,2,opt,name=remainingHealth" json:"remainingHealth,omitempty"`
}

func (m *GetShotResponse) Reset()                    { *m = GetShotResponse{} }
func (m *GetShotResponse) String() string            { return proto.CompactTextString(m) }
func (*GetShotResponse) ProtoMessage()               {}
func (*GetShotResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *GetShotResponse) GetVictimName() string {
	if m != nil {
		return m.VictimName
	}
	return ""
}

func (m *GetShotResponse) GetRemainingHealth() int32 {
	if m != nil {
		return m.RemainingHealth
	}
	return 0
}

func init() {
	proto.RegisterType((*GetShotRequest)(nil), "GetShotRequest")
	proto.RegisterType((*GetShotResponse)(nil), "GetShotResponse")
}

func init() { proto.RegisterFile("cowboy.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 196 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x8f, 0xc1, 0xcb, 0x82, 0x40,
	0x10, 0x47, 0x3f, 0x3f, 0xc8, 0x6a, 0x0a, 0x8d, 0x3d, 0x49, 0x87, 0x10, 0x0f, 0xe1, 0x69, 0x83,
	0xba, 0x74, 0xae, 0xa0, 0x4e, 0x1d, 0xec, 0x66, 0xa7, 0x55, 0x06, 0x5d, 0x68, 0x77, 0xcc, 0xdd,
	0x8a, 0xfe, 0xfb, 0x40, 0x24, 0xcc, 0xdb, 0xf0, 0x0e, 0x6f, 0x7e, 0x0f, 0xa6, 0x39, 0xbd, 0x32,
	0x7a, 0xf3, 0xaa, 0x26, 0x4b, 0x51, 0x0a, 0xde, 0x11, 0xed, 0xa5, 0x24, 0x9b, 0xe0, 0xfd, 0x81,
	0xc6, 0xb2, 0x10, 0x26, 0xa6, 0x24, 0xb2, 0x58, 0x9f, 0x85, 0xc2, 0xc0, 0x09, 0x9d, 0x78, 0x9c,
	0x74, 0x11, 0x5b, 0x82, 0x27, 0x75, 0x4e, 0x4a, 0xea, 0xe2, 0x20, 0x94, 0x28, 0x30, 0xf8, 0x0f,
	0x9d, 0x78, 0x90, 0xf4, 0x68, 0x74, 0x05, 0xff, 0xeb, 0x36, 0x15, 0x69, 0x83, 0x6c, 0x01, 0xf0,
	0x94, 0xb9, 0x95, 0xaa, 0xe3, 0xee, 0x10, 0x16, 0x83, 0x5f, 0xa3, 0x12, 0x52, 0x4b, 0x5d, 0x9c,
	0x50, 0xdc, 0x6c, 0xd9, 0xba, 0xfb, 0x78, 0xbd, 0x05, 0x77, 0xdf, 0x84, 0x30, 0x0e, 0xc3, 0xf6,
	0x0d, 0xf3, 0xf9, 0x6f, 0xcc, 0x7c, 0xc6, 0x7b, 0x0b, 0xa2, 0xbf, 0x1d, 0xa4, 0xa3, 0xa6, 0x7d,
	0x55, 0x65, 0x99, 0xdb, 0x5c, 0x9b, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x18, 0x9d, 0x43, 0xd0,
	0x15, 0x01, 0x00, 0x00,
}
