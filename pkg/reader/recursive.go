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
		focusedType := getToBeCheckedMessageType(pbPkgName, t)
		if focusedType != nil {
			rc.nestedNames[focusedType.GetFullyQualifiedName()] = struct{}{}
		}
	}

	current := getToBeCheckedMessageType(pbPkgName, currentMsgDescriptor)
	if current == nil {
		return
	}

	if !rc.contains(current) {
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
			if !rc.contains(next) {
				rc.traverse(pbPkgName, next)
			}
		case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
			v := field.GetEnumType()
			if !rc.contains(v) {
				rc.founds[v.GetFullyQualifiedName()] = v
			}
		}

	}
}

func (rc *RecursiveChecker) contains(d desc.Descriptor) bool {
	descName := d.GetFullyQualifiedName()
	if _, exists1 := rc.nestedNames[descName]; !exists1 {
		if _, exists2 := rc.founds[descName]; !exists2 {
			return false
		}
	}
	return true
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
