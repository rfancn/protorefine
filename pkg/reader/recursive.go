package reader

import (
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/types/descriptorpb"
)

type RecursiveChecker struct {
	founds map[string]desc.Descriptor
}

func newRecursiveChecker() *RecursiveChecker {
	return &RecursiveChecker{
		founds: make(map[string]desc.Descriptor),
	}
}

func (rc *RecursiveChecker) traverse(pbPkgName string, currentMsgDescriptor *desc.MessageDescriptor) {
	current := getToBeCheckedMessageType(pbPkgName, currentMsgDescriptor)
	if current == nil {
		return
	}

	if _, exists := rc.founds[current.GetFullyQualifiedName()]; !exists {
		rc.founds[current.GetFullyQualifiedName()] = current
	}

	for _, field := range current.GetFields() {
		switch field.GetType() {
		case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
			next := getToBeCheckedMessageType(pbPkgName, field.GetMessageType())
			if next == nil {
				continue
			}

			// avoid recursive definition
			// message AreaNode{
			//  ...
			//  AreaNode children = 1;
			//  ...
			// }
			if _, exists := rc.founds[next.GetFullyQualifiedName()]; !exists {
				rc.traverse(pbPkgName, next)
			}
		case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
			v := field.GetEnumType()
			if _, exists := rc.founds[v.GetFullyQualifiedName()]; !exists {
				rc.founds[v.GetFullyQualifiedName()] = v
			}
		}

	}
}

func getToBeCheckedMessageType(pbPkgName string, d *desc.MessageDescriptor) *desc.MessageDescriptor {
	// if it is not the same package, ignore it
	if d.GetFile().GetPackage() != pbPkgName {
		return nil
	}

	// if it is not map, return it
	if !d.IsMapEntry() {
		return d
	}

	// if it is map, return the value message type
	mapValue := d.GetFields()[1]
	if d.GetFields()[1].GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
		return mapValue.GetMessageType()
	}
	return nil
}
