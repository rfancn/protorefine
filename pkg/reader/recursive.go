package reader

import (
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/types/descriptorpb"
)

type RecursiveChecker struct {
	nestedNames map[string]struct{}
	founds      []desc.Descriptor
}

func newRecursiveChecker() *RecursiveChecker {
	return &RecursiveChecker{
		nestedNames: make(map[string]struct{}),
		founds:      make([]desc.Descriptor, 0),
	}
}

func (rc *RecursiveChecker) traverse(pbPkgName string, currentMsgDescriptor *desc.MessageDescriptor) {
	// if curren msg descriptor exists in nested then ignore it
	// because protoc_gen will compile it as expected
	for _, t := range currentMsgDescriptor.GetNestedMessageTypes() {
		rc.nestedNames[t.GetFullyQualifiedName()] = struct{}{}
	}

	if _, exists := rc.nestedNames[currentMsgDescriptor.GetFullyQualifiedName()]; !exists {
		if currentMsgDescriptor.GetFile().GetPackage() == pbPkgName {
			rc.founds = append(rc.founds, currentMsgDescriptor)
		}
	}

	for _, field := range currentMsgDescriptor.GetFields() {
		switch field.GetType() {
		case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
			// avoid recursive definition
			// message AreaNode{
			//  ...
			//  AreaNode children = 1;
			//  ...
			// }
			parent := field.GetParent().GetFullyQualifiedName()
			owner := field.GetOwner().GetFullyQualifiedName()
			if parent != owner {
				rc.traverse(pbPkgName, field.GetMessageType())
			}

		case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
			if field.GetFile().GetPackage() == pbPkgName {
				rc.founds = append(rc.founds, field.GetEnumType())
			}
		}

	}
}