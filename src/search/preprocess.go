package search

import (
    "path/filepath"
    "os"
    "regexp"
)

var (
    PathCtt = []string{}
)

func PreProcess() {
    filepath.Walk(docPath, walkFn)
}

func walkFn(walkPath string, info os.FileInfo, err error) error {
	isMd, err := regexp.MatchString(mdReg, walkPath)
	if !info.IsDir() {
		if err == nil && isMd {
			PathCtt = append(PathCtt, walkPath)
		}
	}
	return nil
}