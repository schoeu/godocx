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
	// /^\s*#+\s?([^#\r\n]+)/
	titleReg = regexp.MustCompile("^[\t\n\f\r ]*#+[\t\n\f\r ]?([^#\r\n]+)")
	// /<title>(.+?)<\/title>/
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
					content := GetConent(fc.path)
					fc.title = GetTitle(extName, content)
					// markdown转换html
					ConvMd(content)
				}
			}
			
		}
	}
}

func GetRsHTML (path string) []byte{
	content := GetConent(path)
	rsHtml := ConvMd(content)
	return rsHtml
}

// md转html
func ConvMd(content []byte) []byte{
	output := blackfriday.MarkdownBasic(content)
	//ioutil.WriteFile(".+.md$", output, 0777)
	return output
}

// 获取标题
func GetTitle(extName string, content []byte) string {
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
func GetConent(path string) []byte {
	absPath := filepath.Join(docPath,path)
	content, err := ioutil.ReadFile(absPath)

	if err != nil {
		log.Fatal(err)
	}
	return content
}