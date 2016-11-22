package util

import (
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"conf"
)

var (
	docPath = "/Users/memee/Downloads/svn/ps-fe"
	mdReg   = ".+.md$"
	// /^\s*#+\s?([^#\r\n]+)/
	titleReg = regexp.MustCompile("^\\s*#+\\s?([^#\\r\\n]+)")
	// /<title>(.+?)<\/title>/
	htmlTitleReg = regexp.MustCompile("<title>(.+?)<\\/title>")

	// 配置文件变量
	ignoreDir = conf.DocxConf.GetJson("ignoreDir").([]interface{})
	// docNames = conf.DocxConf.GetJson("docName").([]map[string]string)
	docNames = conf.DocxConf.GetJson("docName").([]interface{})
)

type fileCache struct {
	title string
	path  string
	ty    string
	child *[]fileCache
}

var docTree = make([]fileCache, 0)

type fileTrasName map[string]string

var fileNameMap = fileTrasName{}

func ReadDirRs() []fileCache {
	fileNameMap = getDocNames(docNames)
	makeDomTree(docPath, &docTree)
	return docTree
}

// 遍历文件生成文档层级树
func makeDomTree(crtPath string, ctt *[]fileCache) {
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
			fileName := file.Name()
			hitted := indexOf(ignoreDir, fileName)
			if !hitted {
				subFileCache := make([]fileCache, 0)

				// 如果在目录map内有值，则目录名为设定值，没有则默认为文件夹名
				fc.title = fileName
				if fileNameMap[fileName] != "" {
					fc.title = fileNameMap[fileName]
				}

				fc.path = fileName
				fc.ty = "dir"
				fc.child = &subFileCache
				*ctt = append(*ctt, fc)

				makeDomTree(filepath.Join(crtPath, file.Name()), fc.child)
			}
		} else {
			relFile := filepath.Join(crtPath, docName)
			relPath := strings.Replace(crtPath, docPath, "", -1)
			fc.path = filepath.Join(relPath, docName)
			isMd, err := regexp.MatchString(mdReg, fc.path)
			extName := filepath.Ext(fc.path)
			if err == nil {
				if isMd || extName == ".html" || extName == ".htm" {
					content := GetConent(relFile)
					fc.title = GetTitle(extName, content)
					fc.ty = "file"
					// markdown转换html
					ConvMd(content)
					*ctt = append(*ctt, fc)
				}
			}
		}
	}
}

func MakeNav(treeData *[]fileCache) string {
	htmlStr := ""
	makeNavHtml(&htmlStr, treeData)
	return htmlStr
}

// 生成目录树，TODO: 使用template
func makeNavHtml(str *string, data *[]fileCache) {
	for _, v := range *data {
		fileType := v.ty
		if fileType == "file" {
			*str += "<li class='nav nav-title docx-files' data-path='" + v.path + " data-title='" + v.title + "><a href='" + v.path + "'>" + v.title + "</a></li>"
		} else if fileType == "dir" {
			*str += "<li data-dir='" + v.path + "' data-title='" + v.title + "' class='docx-dir'><a href='#' class='docx-dirsa'>" + v.title + "<span class='fa arrow'></span></a><ul class='docx-submenu'>"
			makeNavHtml(str, v.child)
			*str += "</ul></li>"
		}
	}
}

// 获取最终html字符串
func GetRsHTML(path string) []byte {
	content := GetConent(path)
	rsHtml := ConvMd(content)
	return rsHtml
}

// md转html
func ConvMd(content []byte) []byte {
	output := blackfriday.MarkdownBasic(content)
	//ioutil.WriteFile(".+.md$", output, 0777)
	return output
}

// 获取标题
func GetTitle(extName string, content []byte) string {
	var titleCt [][]string
	var title string
	contentStr := string(content)
	if extName == ".md" {
		titleCt = titleReg.FindAllStringSubmatch(contentStr, -1)
	} else if extName == ".html" || extName == ".htm" {
		titleCt = titleReg.FindAllStringSubmatch(contentStr, -1)
	}
	for _, v := range titleCt {
		title = v[1]
		break
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

// 检测文件是否存在
func isExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return true
	}
	return os.IsExist(err)
}

// []interface{} indexOf
func indexOf(s []interface{}, oriVal string) bool {
	for _, val := range s {
		if val == oriVal {
			return true
		}
	}
	return false
}

// 转换文件名配置为map
func getDocNames(docs []interface{}) map[string]string {
	tempTrasName := map[string]string{}
	for _, v := range docs {
		t := v.(map[string]interface{})
		k := t["name"].(string)
		val := t["trans"].(string)
		tempTrasName[k] = val
	}
	return tempTrasName
}

// 获取面包屑数据
func GetPjaxContent(path string) []string {
	pathCtt := strings.Split(path, "/")
	paths := []string{}

	for _, v := range pathCtt {
		paths = append(paths, fileNameMap[v])
	}
	return paths[1:]
}
