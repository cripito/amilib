// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v4.24.2
// source: protos/event.proto

package event

import (
	node "./node"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID        string     `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	PayLoad   string     `protobuf:"bytes,2,opt,name=PayLoad,proto3" json:"PayLoad,omitempty"`
	Node      *node.Node `protobuf:"bytes,3,opt,name=Node,proto3" json:"Node,omitempty"`
	TimeStamp string     `protobuf:"bytes,4,opt,name=TimeStamp,proto3" json:"TimeStamp,omitempty"`
	Type      string     `protobuf:"bytes,5,opt,name=Type,proto3" json:"Type,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_protos_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_protos_event_proto_rawDescGZIP(), []int{0}
}

func (x *Event) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Event) GetPayLoad() string {
	if x != nil {
		return x.PayLoad
	}
	return ""
}

func (x *Event) GetNode() *node.Node {
	if x != nil {
		return x.Node
	}
	return nil
}

func (x *Event) GetTimeStamp() string {
	if x != nil {
		return x.TimeStamp
	}
	return ""
}

func (x *Event) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

var File_protos_event_proto protoreflect.FileDescriptor

var file_protos_event_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6e, 0x6f, 0x64,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7e, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44,
	0x12, 0x18, 0x0a, 0x07, 0x50, 0x61, 0x79, 0x4c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x50, 0x61, 0x79, 0x4c, 0x6f, 0x61, 0x64, 0x12, 0x19, 0x0a, 0x04, 0x4e, 0x6f,
	0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52,
	0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61,
	0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74,
	0x61, 0x6d, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_event_proto_rawDescOnce sync.Once
	file_protos_event_proto_rawDescData = file_protos_event_proto_rawDesc
)

func file_protos_event_proto_rawDescGZIP() []byte {
	file_protos_event_proto_rawDescOnce.Do(func() {
		file_protos_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_event_proto_rawDescData)
	})
	return file_protos_event_proto_rawDescData
}

var file_protos_event_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_protos_event_proto_goTypes = []any{
	(*Event)(nil),     // 0: Event
	(*node.Node)(nil), // 1: Node
}
var file_protos_event_proto_depIdxs = []int32{
	1, // 0: Event.Node:type_name -> Node
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_protos_event_proto_init() }
func file_protos_event_proto_init() {
	if File_protos_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_event_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Event); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protos_event_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protos_event_proto_goTypes,
		DependencyIndexes: file_protos_event_proto_depIdxs,
		MessageInfos:      file_protos_event_proto_msgTypes,
	}.Build()
	File_protos_event_proto = out.File
	file_protos_event_proto_rawDesc = nil
	file_protos_event_proto_goTypes = nil
	file_protos_event_proto_depIdxs = nil
}
