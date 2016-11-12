package walker

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"github.com/russross/blackfriday"
)
var (
	docPath    = "/Users/memee/Downloads/svn/ps-fe"
	ignoreDirs = []string{"img", ".git", ".svn", "courseware", "headline", "imgs", "js", "less", "assets"}
	mdReg = ".+.md$"
	titleReg = regexp.MustCompile("^[\t\n\f\r ]*#+[\t\n\f\r ]?([^#\r\n]+)")
	htmlTitleReg = regexp.MustCompile("<title>(.+?)</title>")
)

type fileCache struct {
	title string
	path string
}

func ReadDirRs() {
	fc := fileCache{}
	files, err := ioutil.ReadDir(docPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		isDir := file.IsDir()
		docName := file.Name()
		// 文件夹需要递归处理，文件则直接存容器
		if isDir {
			//fmt.Println(docName)
		} else {
			fmt.Println("file.Name()", file.Name())
			// title := re.FindString("tablett")
			fc.path = filepath.Join(docPath, docName)
			isMd, err := regexp.MatchString(mdReg, fc.path)
			extName := filepath.Ext(fc.path)
			if err == nil {
				if isMd || extName == ".html" || extName == ".htm"{
					content := getConent(fc.path)
					fc.title = getTitle(extName, content)
					fmt.Println("title->%s", fc.title)
					// markdown转换html
					convMd(content)
				}
			}
			
		}
	}
}
/*
		if (ext === '.md') {
            // /^\s*\#+\s?(.+)/
            // /^\s*#+\s?([^#\s]+)/
            // /^\s*\#+\s?([^\#]+)\s*\#?/
            titleArr =  /^\s*#+\s?([^#\r\n]+)/.exec(content) || [];
            return titleArr[1] || '';
        }
        else if (ext === '.html' || ext === '.htm'){
            titleArr = /<title>(.+?)<\/title>/.exec(content) || [];
            return titleArr[1] || '';
        }

*/


// md转html
func convMd(content []byte) {
	output := blackfriday.MarkdownBasic(content)
	ioutil.WriteFile("./md", output, 0777)
}

// 获取标题
func getTitle(extName string, content []byte) string {
	contentStr := string(content)
	var title string
	if extName == ".md" {
		title = titleReg.FindString(contentStr)
	} else if extName == ".html" || extName == ".htm" {
		title = htmlTitleReg.FindString(contentStr)
	}
	return title
}

// 获取文档内容
func getConent(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("error path->%s", path)
		log.Fatal(err)
	}
	return content
}