package render

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/rfancn/protorefine/pkg/reader"
	"github.com/rfancn/protorefine/utils"
	"os"
	"path/filepath"
	"text/template"
)

const (
	protoTemplate = `
syntax = "proto3";

package pb;

{{range .ImportFiles}}
import "{{.}}";
{{end}}

{{range .TypeDefs}}
{{.}}
{{end}}
`
)

type ProtoRender struct {
	outputFilename string
	protoDir       string
	outputDir      string
}

func New(protoDir, outputDir, outputFileName string) *ProtoRender {
	return &ProtoRender{
		outputFilename: outputFileName,
		protoDir:       protoDir,
		outputDir:      outputDir,
	}
}

func (p *ProtoRender) Render(pbTypeDefs []*reader.PbTypeDef) error {
	importFiles, dependentFiles, err := p.findImportAndDependents(pbTypeDefs)
	if err != nil {
		return err
	}

	err = os.MkdirAll(p.outputDir, 0644)
	if err != nil {
		return errors.Wrapf(err, "make output dir, dir: %s", p.outputDir)
	}

	err = p.generateProtoFile(newProtoRenderDescriptor(importFiles, pbTypeDefs))
	if err != nil {
		return err
	}

	for _, f := range dependentFiles {
		err = utils.CopyFile(filepath.Join(p.protoDir, f), filepath.Join(p.outputDir, f))
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *ProtoRender) generateProtoFile(data *ProtoRenderDescriptor) error {
	t := template.Must(template.New("proto").Parse(protoTemplate))
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return err
	}

	outputPath := filepath.Join(p.outputDir, p.outputFilename+".proto")
	err := os.WriteFile(outputPath, buf.Bytes(), 0644)
	if err != nil {
		return errors.Wrapf(err, "write proto file, path: %s", outputPath)
	}
	return nil
}
