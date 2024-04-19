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
		TypeDefs:    make([]string, len(pbTypeDefs)),
	}

	// sort pb type defs
	pbTypeDefs = pie.SortUsing(pbTypeDefs, func(a, b *reader.PbTypeDef) bool {
		return a.Kind > b.Kind && a.Name < b.Name
	})

	for i, pbTypeDef := range pbTypeDefs {
		item.TypeDefs[i] = pbTypeDef.Definition
	}
	return item
}
