package finder

import (
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/rfancn/protorefine/utils"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type golangAstFinder struct {
	*baseFinder
}

func newGolangAstFinder(sourceDir, pbPkgPath string) *golangAstFinder {
	return &golangAstFinder{baseFinder: newFinder(sourceDir, ".go", pbPkgPath)}
}

// Find all <pkg>.<pbtype> in golang source codes
func (a golangAstFinder) Find(skipDirs ...string) ([]string, error) {
	st, err := os.Stat(a.sourceDir)
	if err != nil {
		return nil, err
	}

	if !st.IsDir() {
		return nil, fmt.Errorf("invalid source code dir, dir: %s", a.sourceDir)
	}

	results := make(map[string]struct{})
	name2pkgImportPath := make(map[string]string)
	_ = filepath.Walk(a.sourceDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && pie.Contains(skipDirs, info.Name()) {
			return filepath.SkipDir
		}

		if !info.IsDir() && strings.HasSuffix(path, a.fileSuffix) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				utils.Fatalf(err, "golang ast parse file, path: %s", path)
			}

			ast.Inspect(f, func(node ast.Node) bool {
				switch n := node.(type) {
				case *ast.SelectorExpr:
					if ident, ok := n.X.(*ast.Ident); ok {
						// if type's package name equal to pbPkgImportPath
						if name2pkgImportPath[ident.Name] == a.pbPkgImportPath {
							if _, exists := results[n.Sel.Name]; !exists {
								results[n.Sel.Name] = struct{}{}
							}
						}
					}
				case *ast.ImportSpec: // record all name2pkgImportPath
					var alias string
					if n.Name != nil {
						alias = n.Name.Name
					}
					fullPath := n.Path.Value[1 : len(n.Path.Value)-1]
					pkgName := filepath.Base(fullPath)
					if alias == "" {
						name2pkgImportPath[pkgName] = fullPath
					} else {
						name2pkgImportPath[alias] = fullPath
					}
				}
				return true
			})
		}
		return nil
	})

	return pie.Keys(results), nil
}
