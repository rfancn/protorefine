package render

import (
	"github.com/elliotchance/pie/v2"
	"github.com/rfancn/protorefine/config"
	"github.com/rfancn/protorefine/pkg/reader"
	"regexp"
)

func (p *ProtoRender) findImportAndDependents(pbTypeDefs []*reader.PbTypeDef) ([]string, []string, error) {
	var err error
	importRegexps := make([]*regexp.Regexp, len(config.Config.Import.Rules))
	for i, rule := range config.Config.Import.Rules {
		importRegexps[i], err = regexp.Compile(rule.Match)
		if err != nil {
			return nil, nil, err
		}
	}

	// check import files and dependents
	imports := make(map[string]struct{})
	dependents := make(map[string]struct{})
	for _, pbTypeDef := range pbTypeDefs {
		for i, importRule := range config.Config.Import.Rules {
			if importRegexps[i].MatchString(pbTypeDef.Definition) {
				// directly import proto files
				if _, exists := imports[importRule.ProtoFile]; !exists {
					imports[importRule.ProtoFile] = struct{}{}
				}

				// dependent proto files need copy
				if _, exists := dependents[importRule.ProtoFile]; !exists {
					dependents[importRule.ProtoFile] = struct{}{}
				}
				for _, d := range importRule.Dependents {
					if _, exists := dependents[d]; !exists {
						dependents[d] = struct{}{}
					}
				}
			}
		}
	}

	importFiles := pie.SortUsing(pie.Keys(imports), func(a, b string) bool {
		return a < b
	})

	dependentFiles := pie.Keys(dependents)

	return importFiles, dependentFiles, nil
}
