package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"os"
	"strings"
	"github.com/russross/blackfriday"
	"html/template"
	"net/http"
)
var (
	docPath    = "/Users/memee/Downloads/svn/ps-fe"
	ignoreDirs = []string{"img", ".git", ".svn", "courseware", "headline", "imgs", "js", "less", "assets"}
	mdReg = ".+.md$"
	// /^\s*#+\s?([^#\r\n]+)/
	titleReg = regexp.MustCompile("^\\s*#+\\s?([^#\\r\\n]+)")
	// /<title>(.+?)<\/title>/
	htmlTitleReg = regexp.MustCompile("<title>(.+?)<\\/title>")
)

type fileCache struct {
	title string
	path string
	ty string
	child interface{}
}

var docTree = make([]interface{},0)

func ReadDirRs() {
	//makeDomTree(docPath, &docTree)

	sep := string(filepath.Separator)

	docPathReg := regexp.MustCompile(docPath)
	filepath.Walk(docPath, func (fPath string, info os.FileInfo, err error) error {
		// isDir := info.IsDir()
		relPath := docPathReg.ReplaceAllString(fPath, "")
		justDir := filepath.Dir(relPath)
		pathArr := strings.Split(justDir, sep)
		fmt.Println(relPath,"~~~~",pathArr)
		// if isDir {
		// 	fmt.Println(relPath, "dir")
		// } else {
		// 	fmt.Println(relPath, "file")
		// }
		return nil
	})

	// fmt.Println(docTree, cap(docTree), len(docTree))
}

func makeDomTree (crtPath string, ctt *[]interface{}) {
	files, err := ioutil.ReadDir(crtPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fc := fileCache{}
		isDir := file.IsDir()
		docName := file.Name()

		// 文件夹需要递归处理，文件则直接存容器
		if isDir {
			fc.title = file.Name()
			fc.path = file.Name()
			fc.ty = "dir"
			fc.child = fc
			aCtt := append(*ctt, fc)
			fmt.Println(aCtt)
			makeDomTree(filepath.Join(crtPath, file.Name()), &aCtt)
			
		} else {
			// fmt.Println("file.Name()", file.Name())
			// title := re.FindString("tablett")
			fc.path = filepath.Join(docPath, docName)
			isMd, err := regexp.MatchString(mdReg, fc.path)
			extName := filepath.Ext(fc.path)
			if err == nil {
				if isMd || extName == ".html" || extName == ".htm"{
					content := GetConent(fc.path)
					fc.title = GetTitle(extName, content)
					fc.ty = "file"
					// markdown转换html
					ConvMd(content)
					//aCtt := append(*ctt, fc)
				}
			}
		}


		// for _, ignPath := range ignoreDirs {
		// 	fmt.Println(docName, ignPath)
		// 	if docName != ignPath {
		// 		// 文件夹需要递归处理，文件则直接存容器
		// 		if isDir {
		// 			fc.title = file.Name()
		// 			fc.path = file.Name()
		// 			fc.ty = "dir"
		// 			fc.child = fc
		// 			*ctt = append(*ctt, fc)
		// 			makeDomTree(filepath.Join(crtPath, file.Name()), ctt)
		// 		} else {
		// 			// fmt.Println("file.Name()", file.Name())
		// 			// title := re.FindString("tablett")
		// 			fc.path = filepath.Join(docPath, docName)
		// 			isMd, err := regexp.MatchString(mdReg, fc.path)
		// 			extName := filepath.Ext(fc.path)
		// 			if err == nil {
		// 				if isMd || extName == ".html" || extName == ".htm"{
		// 					content := GetConent(fc.path)
		// 					fc.title = GetTitle(extName, content)
		// 					fc.ty = "file"
		// 					// markdown转换html
		// 					ConvMd(content)
		// 					*ctt = append(*ctt, fc)
		// 				}
		// 			}
		// 		}
		// 	} else {
		// 		continue
		// 	}
		//}
	}
}

func GetRsHTML (path string) []byte{
	fmt.Println("path",path)
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
	nocontent := []byte("")
	stat, _ := PathExists(path)

	if stat {
		content, err := ioutil.ReadFile(path)

		if err != nil {
			log.Fatal(err)
		}
		return content
	}
	return nocontent
}

// 判断文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 渲染模板
func RenderTpl(path string, data interface{}, w http.ResponseWriter) {
	t, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}
