package utils

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//ResolvePath combine path with current folder
func ResolvePath(fileName string) string {
	if filepath.IsAbs(fileName) {
		return fileName
	}

	currentDir := GetWd()

	return filepath.Join(currentDir, fileName)
}

//
func IsTesting() bool {
	arg := os.Args[0]
	return strings.HasSuffix(arg, ".test.exe") || flag.Lookup("test.v") != nil
}

//
func GetWd() string {
	var currentDir string
	var err error

	if IsTesting() {
		currentDir, err = os.Getwd()
	} else {
		currentDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	}

	if err != nil {
		log.Fatal(err)
	}

	return currentDir
}
