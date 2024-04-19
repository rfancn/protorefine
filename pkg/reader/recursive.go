package reader

import (
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/types/descriptorpb"
)

type RecursiveChecker struct {
	nestedNames map[string]struct{}
	founds      map[string]desc.Descriptor
}

func newRecursiveChecker() *RecursiveChecker {
	return &RecursiveChecker{
		nestedNames: make(map[string]struct{}),
		founds:      make(map[string]desc.Descriptor),
	}
}

func (rc *RecursiveChecker) traverse(pbPkgName string, currentMsgDescriptor *desc.MessageDescriptor) {
	// if curren msg descriptor exists in nested then ignore it
	// because protoc_gen will compile it as expected
	for _, t := range currentMsgDescriptor.GetNestedMessageTypes() {
		rc.nestedNames[t.GetFullyQualifiedName()] = struct{}{}
	}

	current := getToBeCheckedMessageType(pbPkgName, currentMsgDescriptor)
	if current == nil {
		return
	}

	if _, exists1 := rc.nestedNames[current.GetFullyQualifiedName()]; !exists1 {
		if _, exists2 := rc.founds[current.GetFullyQualifiedName()]; !exists2 {
			rc.founds[current.GetFullyQualifiedName()] = current
		}
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
