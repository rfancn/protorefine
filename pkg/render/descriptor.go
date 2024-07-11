package render

import (
	"github.com/elliotchance/pie/v2"
	"github.com/rfancn/protorefine/pkg/reader"
)

type ProtoRenderDescriptor struct {
	ImportFiles []string
	TypeDefs    []string
}

func newProtoRenderDescriptor(importFiles []string, pbTypeDefs []*reader.PbTypeDef) *ProtoRenderDescriptor {
	item := &ProtoRenderDescriptor{
		ImportFiles: importFiles,
		TypeDefs:    make([]string, 0),
	}

	// sort pb type defs
	kind2typeDefs := pie.GroupBy(pbTypeDefs, func(item *reader.PbTypeDef) int {
		return int(item.Kind)
	})

	sorted := make(map[int][]*reader.PbTypeDef)
	for kind, typeDefs := range kind2typeDefs {
		typeDefs = pie.SortUsing(typeDefs, func(a, b *reader.PbTypeDef) bool {
			return a.Name < b.Name
		})
		sorted[kind] = typeDefs
	}

	sortedKinds := pie.SortUsing(pie.Keys(kind2typeDefs), func(a, b int) bool {
		return a < b
	})

	for _, kind := range sortedKinds {
		for _, def := range sorted[kind] {
			item.TypeDefs = append(item.TypeDefs, def.Definition)
		}
	}
	return item
}
