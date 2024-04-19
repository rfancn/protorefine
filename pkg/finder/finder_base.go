package finder

type baseFinder struct {
	sourceDir       string
	fileSuffix      string // sourceDir code file fileSuffix
	pbPkgImportPath string // protobuf package name
}

func newFinder(sourceDir, fileSuffix, pbPkgPath string) *baseFinder {
	return &baseFinder{
		sourceDir:       sourceDir,
		fileSuffix:      fileSuffix,
		pbPkgImportPath: pbPkgPath,
	}
}
