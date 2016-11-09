package walker

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type fileInfos struct {
	path  string
	title string
}

var fis fileInfos

var pathSlice = []string{}

type doctype struct {
	title string
	path  string
}

var (
	docPath    = "/Users/memee/Downloads/svn/ps-fe"
	ignoreDirs = []string{"img", ".git", ".svn", "courseware", "headline", "imgs", "js", "less", "assets"}
)

func ReadDirRs() {
	// walkErr := filepath.Walk(docPath, walkFunc)
	// if walkErr != nil {
	// 	log.Fatal(walkErr)
	// }

	files, err := ioutil.ReadDir(docPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
}

func walkFunc(walkPath string, info os.FileInfo, err error) error {

	if info.IsDir() {
		// pathArr := strings.Split(walkPath, "/")

		// for _, v := range ignoreDirs {
		// 	if v != pathArr[len(pathArr)-1] {
		// 		pathSlice := append(pathSlice, walkPath)
		// 		fmt.Println(pathSlice)
		// 	}
		// }
	} else {
		// 如果是文件夹
		/*walkErr := filepath.Walk(walkPath, walkFunc)
		if walkErr != nil {
			log.Fatal(walkErr)
		}*/
		var fi fileInfos
		fi.path = walkPath

		// for _, v := range ignoreDirs {
		// 	if v != pathArr[len(pathArr)-1] {
		// 		docinfo = append(docinfo, path.Dir(walkPath))
		// 	}
		// }
	}
	return nil
}
