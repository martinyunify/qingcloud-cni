// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: contract.proto

/*
	Package messages is a generated protocol buffer package.

	It is generated from these files:
		contract.proto

	It has these top-level messages:
		AllocateNicMessage
		AllocateNicReplyMessage
		GetGatewayMessage
		GetGatewayReplyMessage
		DeleteNicMessage
*/
package messages

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/AsynkronIT/protoactor-go/actor"

import strings "strings"
import reflect "reflect"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type AllocateNicMessage struct {
	Name string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
}

func (m *AllocateNicMessage) Reset()                    { *m = AllocateNicMessage{} }
func (*AllocateNicMessage) ProtoMessage()               {}
func (*AllocateNicMessage) Descriptor() ([]byte, []int) { return fileDescriptorContract, []int{0} }

func (m *AllocateNicMessage) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type AllocateNicReplyMessage struct {
	Name       string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	NetworkID  string `protobuf:"bytes,2,opt,name=NetworkID,proto3" json:"NetworkID,omitempty"`
	EndpointID string `protobuf:"bytes,3,opt,name=EndpointID,proto3" json:"EndpointID,omitempty"`
	Address    string `protobuf:"bytes,4,opt,name=Address,proto3" json:"Address,omitempty"`
}

func (m *AllocateNicReplyMessage) Reset()                    { *m = AllocateNicReplyMessage{} }
func (*AllocateNicReplyMessage) ProtoMessage()               {}
func (*AllocateNicReplyMessage) Descriptor() ([]byte, []int) { return fileDescriptorContract, []int{1} }

func (m *AllocateNicReplyMessage) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *AllocateNicReplyMessage) GetNetworkID() string {
	if m != nil {
		return m.NetworkID
	}
	return ""
}

func (m *AllocateNicReplyMessage) GetEndpointID() string {
	if m != nil {
		return m.EndpointID
	}
	return ""
}

func (m *AllocateNicReplyMessage) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type GetGatewayMessage struct {
	Vxnet string `protobuf:"bytes,1,opt,name=vxnet,proto3" json:"vxnet,omitempty"`
}

func (m *GetGatewayMessage) Reset()                    { *m = GetGatewayMessage{} }
func (*GetGatewayMessage) ProtoMessage()               {}
func (*GetGatewayMessage) Descriptor() ([]byte, []int) { return fileDescriptorContract, []int{2} }

func (m *GetGatewayMessage) GetVxnet() string {
	if m != nil {
		return m.Vxnet
	}
	return ""
}

type GetGatewayReplyMessage struct {
	NetworkID  string `protobuf:"bytes,1,opt,name=NetworkID,proto3" json:"NetworkID,omitempty"`
	EndpointID string `protobuf:"bytes,2,opt,name=EndpointID,proto3" json:"EndpointID,omitempty"`
	Address    string `protobuf:"bytes,3,opt,name=Address,proto3" json:"Address,omitempty"`
}

func (m *GetGatewayReplyMessage) Reset()                    { *m = GetGatewayReplyMessage{} }
func (*GetGatewayReplyMessage) ProtoMessage()               {}
func (*GetGatewayReplyMessage) Descriptor() ([]byte, []int) { return fileDescriptorContract, []int{3} }

func (m *GetGatewayReplyMessage) GetNetworkID() string {
	if m != nil {
		return m.NetworkID
	}
	return ""
}

func (m *GetGatewayReplyMessage) GetEndpointID() string {
	if m != nil {
		return m.EndpointID
	}
	return ""
}

func (m *GetGatewayReplyMessage) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type DeleteNicMessage struct {
	Nicid   string `protobuf:"bytes,1,opt,name=nicid,proto3" json:"nicid,omitempty"`
	Nicname string `protobuf:"bytes,2,opt,name=nicname,proto3" json:"nicname,omitempty"`
}

func (m *DeleteNicMessage) Reset()                    { *m = DeleteNicMessage{} }
func (*DeleteNicMessage) ProtoMessage()               {}
func (*DeleteNicMessage) Descriptor() ([]byte, []int) { return fileDescriptorContract, []int{4} }

func (m *DeleteNicMessage) GetNicid() string {
	if m != nil {
		return m.Nicid
	}
	return ""
}

func (m *DeleteNicMessage) GetNicname() string {
	if m != nil {
		return m.Nicname
	}
	return ""
}

func init() {
	proto.RegisterType((*AllocateNicMessage)(nil), "messages.AllocateNicMessage")
	proto.RegisterType((*AllocateNicReplyMessage)(nil), "messages.AllocateNicReplyMessage")
	proto.RegisterType((*GetGatewayMessage)(nil), "messages.GetGatewayMessage")
	proto.RegisterType((*GetGatewayReplyMessage)(nil), "messages.GetGatewayReplyMessage")
	proto.RegisterType((*DeleteNicMessage)(nil), "messages.DeleteNicMessage")
}
func (this *AllocateNicMessage) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*AllocateNicMessage)
	if !ok {
		that2, ok := that.(AllocateNicMessage)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.Name != that1.Name {
		return false
	}
	return true
}
func (this *AllocateNicReplyMessage) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*AllocateNicReplyMessage)
	if !ok {
		that2, ok := that.(AllocateNicReplyMessage)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.Name != that1.Name {
		return false
	}
	if this.NetworkID != that1.NetworkID {
		return false
	}
	if this.EndpointID != that1.EndpointID {
		return false
	}
	if this.Address != that1.Address {
		return false
	}
	return true
}
func (this *GetGatewayMessage) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*GetGatewayMessage)
	if !ok {
		that2, ok := that.(GetGatewayMessage)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.Vxnet != that1.Vxnet {
		return false
	}
	return true
}
func (this *GetGatewayReplyMessage) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*GetGatewayReplyMessage)
	if !ok {
		that2, ok := that.(GetGatewayReplyMessage)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.NetworkID != that1.NetworkID {
		return false
	}
	if this.EndpointID != that1.EndpointID {
		return false
	}
	if this.Address != that1.Address {
		return false
	}
	return true
}
func (this *DeleteNicMessage) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*DeleteNicMessage)
	if !ok {
		that2, ok := that.(DeleteNicMessage)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.Nicid != that1.Nicid {
		return false
	}
	if this.Nicname != that1.Nicname {
		return false
	}
	return true
}
func (this *AllocateNicMessage) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&messages.AllocateNicMessage{")
	s = append(s, "Name: "+fmt.Sprintf("%#v", this.Name)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *AllocateNicReplyMessage) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.AllocateNicReplyMessage{")
	s = append(s, "Name: "+fmt.Sprintf("%#v", this.Name)+",\n")
	s = append(s, "NetworkID: "+fmt.Sprintf("%#v", this.NetworkID)+",\n")
	s = append(s, "EndpointID: "+fmt.Sprintf("%#v", this.EndpointID)+",\n")
	s = append(s, "Address: "+fmt.Sprintf("%#v", this.Address)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *GetGatewayMessage) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&messages.GetGatewayMessage{")
	s = append(s, "Vxnet: "+fmt.Sprintf("%#v", this.Vxnet)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *GetGatewayReplyMessage) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&messages.GetGatewayReplyMessage{")
	s = append(s, "NetworkID: "+fmt.Sprintf("%#v", this.NetworkID)+",\n")
	s = append(s, "EndpointID: "+fmt.Sprintf("%#v", this.EndpointID)+",\n")
	s = append(s, "Address: "+fmt.Sprintf("%#v", this.Address)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *DeleteNicMessage) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&messages.DeleteNicMessage{")
	s = append(s, "Nicid: "+fmt.Sprintf("%#v", this.Nicid)+",\n")
	s = append(s, "Nicname: "+fmt.Sprintf("%#v", this.Nicname)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringContract(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *AllocateNicMessage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AllocateNicMessage) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintContract(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	return i, nil
}

func (m *AllocateNicReplyMessage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AllocateNicReplyMessage) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintContract(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if len(m.NetworkID) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintContract(dAtA, i, uint64(len(m.NetworkID)))
		i += copy(dAtA[i:], m.NetworkID)
	}
	if len(m.EndpointID) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintContract(dAtA, i, uint64(len(m.EndpointID)))
		i += copy(dAtA[i:], m.EndpointID)
	}
	if len(m.Address) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintContract(dAtA, i, uint64(len(m.Address)))
		i += copy(dAtA[i:], m.Address)
	}
	return i, nil
}

func (m *GetGatewayMessage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GetGatewayMessage) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Vxnet) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintContract(dAtA, i, uint64(len(m.Vxnet)))
		i += copy(dAtA[i:], m.Vxnet)
	}
	return i, nil
}

func (m *GetGatewayReplyMessage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GetGatewayReplyMessage) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.NetworkID) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintContract(dAtA, i, uint64(len(m.NetworkID)))
		i += copy(dAtA[i:], m.NetworkID)
	}
	if len(m.EndpointID) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintContract(dAtA, i, uint64(len(m.EndpointID)))
		i += copy(dAtA[i:], m.EndpointID)
	}
	if len(m.Address) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintContract(dAtA, i, uint64(len(m.Address)))
		i += copy(dAtA[i:], m.Address)
	}
	return i, nil
}

func (m *DeleteNicMessage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DeleteNicMessage) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Nicid) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintContract(dAtA, i, uint64(len(m.Nicid)))
		i += copy(dAtA[i:], m.Nicid)
	}
	if len(m.Nicname) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintContract(dAtA, i, uint64(len(m.Nicname)))
		i += copy(dAtA[i:], m.Nicname)
	}
	return i, nil
}

func encodeFixed64Contract(dAtA []byte, offset int, v uint64) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	dAtA[offset+4] = uint8(v >> 32)
	dAtA[offset+5] = uint8(v >> 40)
	dAtA[offset+6] = uint8(v >> 48)
	dAtA[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Contract(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintContract(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *AllocateNicMessage) Size() (n int) {
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovContract(uint64(l))
	}
	return n
}

func (m *AllocateNicReplyMessage) Size() (n int) {
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovContract(uint64(l))
	}
	l = len(m.NetworkID)
	if l > 0 {
		n += 1 + l + sovContract(uint64(l))
	}
	l = len(m.EndpointID)
	if l > 0 {
		n += 1 + l + sovContract(uint64(l))
	}
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovContract(uint64(l))
	}
	return n
}

func (m *GetGatewayMessage) Size() (n int) {
	var l int
	_ = l
	l = len(m.Vxnet)
	if l > 0 {
		n += 1 + l + sovContract(uint64(l))
	}
	return n
}

func (m *GetGatewayReplyMessage) Size() (n int) {
	var l int
	_ = l
	l = len(m.NetworkID)
	if l > 0 {
		n += 1 + l + sovContract(uint64(l))
	}
	l = len(m.EndpointID)
	if l > 0 {
		n += 1 + l + sovContract(uint64(l))
	}
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovContract(uint64(l))
	}
	return n
}

func (m *DeleteNicMessage) Size() (n int) {
	var l int
	_ = l
	l = len(m.Nicid)
	if l > 0 {
		n += 1 + l + sovContract(uint64(l))
	}
	l = len(m.Nicname)
	if l > 0 {
		n += 1 + l + sovContract(uint64(l))
	}
	return n
}

func sovContract(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozContract(x uint64) (n int) {
	return sovContract(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *AllocateNicMessage) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AllocateNicMessage{`,
		`Name:` + fmt.Sprintf("%v", this.Name) + `,`,
		`}`,
	}, "")
	return s
}
func (this *AllocateNicReplyMessage) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AllocateNicReplyMessage{`,
		`Name:` + fmt.Sprintf("%v", this.Name) + `,`,
		`NetworkID:` + fmt.Sprintf("%v", this.NetworkID) + `,`,
		`EndpointID:` + fmt.Sprintf("%v", this.EndpointID) + `,`,
		`Address:` + fmt.Sprintf("%v", this.Address) + `,`,
		`}`,
	}, "")
	return s
}
func (this *GetGatewayMessage) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&GetGatewayMessage{`,
		`Vxnet:` + fmt.Sprintf("%v", this.Vxnet) + `,`,
		`}`,
	}, "")
	return s
}
func (this *GetGatewayReplyMessage) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&GetGatewayReplyMessage{`,
		`NetworkID:` + fmt.Sprintf("%v", this.NetworkID) + `,`,
		`EndpointID:` + fmt.Sprintf("%v", this.EndpointID) + `,`,
		`Address:` + fmt.Sprintf("%v", this.Address) + `,`,
		`}`,
	}, "")
	return s
}
func (this *DeleteNicMessage) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&DeleteNicMessage{`,
		`Nicid:` + fmt.Sprintf("%v", this.Nicid) + `,`,
		`Nicname:` + fmt.Sprintf("%v", this.Nicname) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringContract(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *AllocateNicMessage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowContract
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AllocateNicMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AllocateNicMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthContract
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipContract(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthContract
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *AllocateNicReplyMessage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowContract
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AllocateNicReplyMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AllocateNicReplyMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthContract
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NetworkID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthContract
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NetworkID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndpointID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthContract
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EndpointID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthContract
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipContract(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthContract
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *GetGatewayMessage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowContract
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GetGatewayMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GetGatewayMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Vxnet", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthContract
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Vxnet = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipContract(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthContract
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *GetGatewayReplyMessage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowContract
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GetGatewayReplyMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GetGatewayReplyMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NetworkID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthContract
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NetworkID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndpointID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthContract
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EndpointID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthContract
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipContract(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthContract
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *DeleteNicMessage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowContract
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DeleteNicMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DeleteNicMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nicid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthContract
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Nicid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nicname", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthContract
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Nicname = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipContract(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthContract
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipContract(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowContract
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowContract
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowContract
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthContract
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowContract
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipContract(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthContract = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowContract   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("contract.proto", fileDescriptorContract) }

var fileDescriptorContract = []byte{
	// 327 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0xbf, 0x4e, 0x32, 0x41,
	0x14, 0xc5, 0x77, 0xf8, 0xf3, 0x7d, 0x72, 0x0b, 0xa3, 0x13, 0xa2, 0x1b, 0x63, 0x26, 0x66, 0x2b,
	0x4c, 0x14, 0x0a, 0x13, 0x7b, 0x08, 0x86, 0x50, 0x48, 0x41, 0x7c, 0x81, 0x61, 0xf6, 0x06, 0x37,
	0x2c, 0x33, 0x9b, 0x9d, 0x51, 0xa4, 0xb3, 0xb1, 0xf7, 0x31, 0x7c, 0x14, 0x4b, 0x4a, 0x4b, 0x19,
	0x1b, 0x4b, 0x1e, 0xc1, 0x30, 0x2b, 0xc2, 0x36, 0xdb, 0xdd, 0x73, 0xcf, 0xc9, 0x99, 0x5f, 0xee,
	0xc0, 0xbe, 0x50, 0xd2, 0xa4, 0x5c, 0x98, 0x66, 0x92, 0x2a, 0xa3, 0xe8, 0xde, 0x14, 0xb5, 0xe6,
	0x63, 0xd4, 0x27, 0xd7, 0xe3, 0xc8, 0xdc, 0x3f, 0x8c, 0x9a, 0x42, 0x4d, 0x5b, 0x6d, 0x3d, 0x97,
	0x93, 0x54, 0xc9, 0xfe, 0x5d, 0xcb, 0xc5, 0xb8, 0x30, 0x2a, 0xbd, 0x1c, 0xab, 0x96, 0x1b, 0xb2,
	0x9d, 0xce, 0x1a, 0x82, 0x06, 0xd0, 0x76, 0x1c, 0x2b, 0xc1, 0x0d, 0x0e, 0x22, 0x71, 0x9b, 0xd5,
	0x51, 0x0a, 0x95, 0x01, 0x9f, 0xa2, 0x4f, 0xce, 0x48, 0xa3, 0x36, 0x74, 0x73, 0xf0, 0x42, 0xe0,
	0x78, 0x27, 0x3a, 0xc4, 0x24, 0x9e, 0x17, 0xe4, 0xe9, 0x29, 0xd4, 0x06, 0x68, 0x66, 0x2a, 0x9d,
	0xf4, 0xbb, 0x7e, 0xc9, 0x19, 0xdb, 0x05, 0x65, 0x00, 0x37, 0x32, 0x4c, 0x54, 0x24, 0x4d, 0xbf,
	0xeb, 0x97, 0x9d, 0xbd, 0xb3, 0xa1, 0x3e, 0xfc, 0x6f, 0x87, 0x61, 0x8a, 0x5a, 0xfb, 0x15, 0x67,
	0x6e, 0x64, 0x70, 0x0e, 0x87, 0x3d, 0x34, 0x3d, 0x6e, 0x70, 0xc6, 0xff, 0x00, 0xea, 0x50, 0x7d,
	0x7c, 0x92, 0x68, 0x7e, 0x09, 0x32, 0x11, 0x24, 0x70, 0xb4, 0x8d, 0xe6, 0x80, 0x73, 0x70, 0xa4,
	0x18, 0xae, 0x54, 0x04, 0x57, 0xce, 0xc3, 0x75, 0xe0, 0xa0, 0x8b, 0x31, 0xe6, 0x8e, 0x59, 0x87,
	0xaa, 0x8c, 0x44, 0x14, 0x6e, 0xd8, 0x9c, 0x58, 0x77, 0xc8, 0x48, 0xc8, 0xf5, 0xd5, 0xb2, 0x07,
	0x36, 0xb2, 0x73, 0xb1, 0x58, 0x32, 0xef, 0x63, 0xc9, 0xbc, 0xd5, 0x92, 0x91, 0x67, 0xcb, 0xc8,
	0x9b, 0x65, 0xe4, 0xdd, 0x32, 0xb2, 0xb0, 0x8c, 0x7c, 0x5a, 0x46, 0xbe, 0x2d, 0xf3, 0x56, 0x96,
	0x91, 0xd7, 0x2f, 0xe6, 0x8d, 0xfe, 0xb9, 0x7f, 0xbc, 0xfa, 0x09, 0x00, 0x00, 0xff, 0xff, 0x0c,
	0x55, 0xcb, 0x00, 0x1b, 0x02, 0x00, 0x00,
}
