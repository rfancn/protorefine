package reader

import (
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoprint"
)

type PbTypeKind int

const (
	PbDeclarationKindUnknown = iota
	PbTypeKindEnum
	PbTypeKindMessage
)

type PbTypeDef struct {
	Name       string
	Kind       PbTypeKind
	Definition string
}

func newPbTypeDef(descriptor desc.Descriptor) (*PbTypeDef, error) {
	pbPrinter := &protoprint.Printer{
		Compact:                              true,
		ShortOptionsExpansionThresholdLength: 200,
	}

	var kind PbTypeKind
	switch descriptor.(type) {
	case *desc.MessageDescriptor:
		kind = PbTypeKindMessage
	case *desc.EnumDescriptor:
		kind = PbTypeKindEnum
	default:
		return nil, fmt.Errorf("unsupported descriptor, %v", descriptor)
	}

	typeDef, err := pbPrinter.PrintProtoToString(descriptor)
	if err != nil {
		return nil, err
	}
	return &PbTypeDef{Name: descriptor.GetName(), Kind: kind, Definition: typeDef}, nil
}
