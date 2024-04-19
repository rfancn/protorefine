package reader

import (
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"os"
	"regexp"
	"strings"
	"unicode"
)

type ProtoReader interface {
	ExtractTypeDefs(protoDir string, pbTypeNames []string) ([]*PbTypeDef, error) // extract protobuf type definitions
}

type pbReader struct {
	pbPkgName           string
	type2msgDescriptor  map[string]*desc.MessageDescriptor
	type2enumDescriptor map[string]*desc.EnumDescriptor
}

func New(pbPkgName string) ProtoReader {
	return &pbReader{
		pbPkgName:           pbPkgName,
		type2msgDescriptor:  make(map[string]*desc.MessageDescriptor),
		type2enumDescriptor: make(map[string]*desc.EnumDescriptor),
	}
}

func (r pbReader) ExtractTypeDefs(protoDir string, pbTypeNames []string) ([]*PbTypeDef, error) {
	err := r.readProtoDescriptors(protoDir)
	if err != nil {
		return nil, err
	}

	if len(r.type2msgDescriptor) == 0 || len(r.type2enumDescriptor) == 0 {
		return nil, fmt.Errorf("proto descriptors not found")
	}

	// filter duplicate pbTypeName corresponding descriptors
	type2descriptors := make(map[string][]desc.Descriptor)
	for _, pbTypeName := range pbTypeNames {
		foundType, foundDescriptors := r.guess(pbTypeName)
		if len(foundDescriptors) == 0 {
			// protoc_gen may generate the protobuf type in source code with different naming convention
			// e,g: it capitalized the first character of the word
			// we should check this case, pb.TheWord ==> message theWord {}
			foundType, foundDescriptors = r.guess(convertFirstLetterToLower(pbTypeName))
			if len(foundDescriptors) == 0 {
				return nil, fmt.Errorf("pb definition %s not found", pbTypeName)
			}
		}

		if _, exists := type2descriptors[foundType]; !exists {
			type2descriptors[foundType] = foundDescriptors
		}
	}

	// filter duplicate descriptors
	descName2descriptor := make(map[string]desc.Descriptor)
	for _, descriptors := range type2descriptors {
		for _, d := range descriptors {
			if _, exists := descName2descriptor[d.GetName()]; !exists {
				descName2descriptor[d.GetName()] = d
			}
		}
	}

	results := make([]*PbTypeDef, 0)
	for _, descriptor := range descName2descriptor {
		t, err := newPbTypeDef(descriptor)
		if err != nil {
			return nil, err
		}

		results = append(results, t)
	}
	return results, nil
}

func (r pbReader) readProtoDescriptors(protoDir string) error {
	protoFiles, err := findProtoFiles(protoDir)
	if err != nil {
		return err
	}

	pbParser := &protoparse.Parser{
		ImportPaths: []string{protoDir},
	}

	results, err := pbParser.ParseFiles(protoFiles...)
	if err != nil {
		return err
	}

	for _, result := range results {
		for _, t := range result.GetMessageTypes() {
			r.type2msgDescriptor[t.GetName()] = t
		}

		for _, t := range result.GetEnumTypes() {
			r.type2enumDescriptor[t.GetName()] = t
		}
	}

	return nil
}

// guess split the name to words and recursive from end to begin to find if pb type can be found
func (r pbReader) guess(pbType string) (string, []desc.Descriptor) {
	words := splitWords(pbType)
	for i := len(words); i > 0; i-- {
		checkType := strings.Join(words[:i], "")
		if descriptor, exist := r.type2msgDescriptor[checkType]; exist {
			// recursively search descendant message types and enums
			checker := newRecursiveChecker()
			checker.traverse(r.pbPkgName, descriptor)
			return checkType, pie.Values(checker.founds)
		}

		if descriptor, exist := r.type2enumDescriptor[checkType]; exist {
			return checkType, []desc.Descriptor{
				descriptor,
			}
		}
	}

	return "", nil
}

func convertFirstLetterToLower(s string) string {
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

func findProtoFiles(baseDir string) ([]string, error) {
	files, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	founds := make([]string, 0)
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".proto") {
			founds = append(founds, f.Name())
		}
	}
	return founds, nil
}

func splitWords(input string) []string {
	var words []string
	if strings.Contains(input, "_") {
		// if name is snake case (例如：my_variable_name)
		words = strings.Split(input, "_")
	} else {
		reg := regexp.MustCompile(`([a-z])([A-Z])`)
		spaceSeparated := reg.ReplaceAllString(input, `${1} ${2}`)

		// 使用 strings.Fields 分隔成单词
		words = strings.Fields(spaceSeparated)
	}
	return words
}
