package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Fatalf(err error, format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	fmt.Printf("%s, error: %+v", s, err)
	os.Exit(1)
}

func Fatal(err error) {
	fmt.Printf("error: %+v", err)
	os.Exit(1)
}

func CopyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	err = os.MkdirAll(filepath.Dir(dest), 0644)
	if err != nil {
		return err
	}

	destFile, err := os.Create(dest) // creates if file doesn't exist
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}
	return nil
}

//
//func CreateTempDir(dirPrefix string) (string, error) {
//	os.MkdirTemp()
//	func generateTempDir(dirPrefix string) (string, error) {
//		dirname := make([]byte, 4)
//		_, err := rand.Read(dirname)
//		if err != nil {
//			return "", err
//		}
//		randomDirName := hex.EncodeToString(dirname)
//		randomDirPath := filepath.Join(os.TempDir(), fmt.Sprintf("%s_%s", dirPrefix, randomDirName))
//
//		//在默认临时文件路径下, 创建一个以dirPrefix为前缀的新的临时目录.
//		err = os.Mkdir(randomDirPath, 0644)
//		if err != nil {
//			return "", err
//		}
//		return randomDirPath, nil
//	}
//
//})
