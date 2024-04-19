package finder

//
//type regexFinder struct {
//	*baseFinder
//	exp *regexp.Regexp
//}
//
//var (
//	regexPbVariable = `(^pb\.\w+$)|(^&pb\.\w+$)`
//)
//
//func newRegexFinder(sourceDir, fileSuffix, pbPkgImportPath string) *regexFinder {
//	exp, _ := regexp.Compile(regexPbVariable)
//	return &regexFinder{baseFinder: newFinder(sourceDir, fileSuffix, pbPkgImportPath), exp: exp}
//}
//
//// Match 尝试找到
//func (r regexFinder) Find(skipDirs ...string) ([]string, error) {
//	st, err := os.Stat(r.sourceDir)
//	if err != nil {
//		return nil, err
//	}
//
//	if !st.IsDir() {
//		return nil, fmt.Errorf("invalid source code dir, dir: %s", r.sourceDir)
//	}
//
//	results := make(map[string]struct{})
//	_ = filepath.Walk(r.sourceDir, func(path string, info os.FileInfo, err error) error {
//		if info.IsDir() && pie.Contains(skipDirs, info.Name()) {
//			return filepath.SkipDir
//		}
//
//		if !info.IsDir() && strings.HasSuffix(path, r.fileSuffix) {
//			f, err := os.Open(path)
//			if err != nil {
//				return err
//			}
//			defer func() {
//				_ = f.Close()
//			}()
//
//			scanner := bufio.NewScanner(f)
//			for scanner.Scan() {
//				s := scanner.Text()
//				if r.exp.MatchString(s) {
//					matches := r.exp.FindAllString(s, -1)
//					for _, matched := range matches {
//						if matched != "" {
//							results[matched] = struct{}{}
//						}
//					}
//				}
//			}
//		}
//		return nil
//	})
//
//	return pie.Keys(results), nil
//}
