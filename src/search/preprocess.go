package search

import (
    "path/filepath"
    "os"
    "regexp"
    "conf"
)

var (
    DocPath  = conf.DocxConf.GetJson("path").(string)
    PathCtt = []string{}
)

type searchContent struct {
    title string
    pos []int
    path string
    spell string
    initials string
}

var sc = []searchContent{}

func PreProcess() {
    filepath.Walk(DocPath, walkFn)
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