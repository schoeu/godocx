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
	docPath      = "/Users/memee/Downloads/svn/ps-fe"
	mdReg        = ".+.md$"
	titleReg     = regexp.MustCompile("^\\s*#+\\s?([^#\\r\\n]+)")
	htmlTitleReg = regexp.MustCompile("<title>(.+?)<\\/title>")

	// 配置文件变量
	ignoreDir = conf.DocxConf.GetJson("ignoreDir").([]interface{})
	docNames  = conf.DocxConf.GetJson("docName").([]interface{})
)

type fileCache struct {
	title string
	path  string
	ty    string
	child *[]fileCache
}

var docTree = make([]fileCache, 0)

// 文件夹排序用
var dirOrder = []string{}

var fileNameMap = map[string]string{}

func ReadDirRs() []fileCache {
	fileNameMap = getDocNames(docNames)

	// 生成树结构
	makeDomTree(docPath, &docTree)

	// 依配置文件循序排序
	var sortedDocTree = dirSort()

	return sortedDocTree
}

// 排序
func dirSort() []fileCache {
	var tempDirSlice = []fileCache{}
	var tempFileSlice = []fileCache{}
	var rsDirSlice = []fileCache{}
	for _, v := range docTree {
		ty := v.ty
		if ty == "file" {
			tempFileSlice = append(tempFileSlice, v)
		} else if ty == "dir" {
			tempDirSlice = append(tempDirSlice, v)
		}
	}

	for _, it := range dirOrder {
		for _, v := range docTree {
			if v.path == it {
				rsDirSlice = append(rsDirSlice, v)
			}
		}
	}

	return append(tempFileSlice, rsDirSlice...)
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
		fileName := file.Name()
		// 文件夹需要递归处理，文件则直接存容器
		if isDir {
			hitted := IndexOf(ignoreDir, fileName) > -1
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
			relFile := filepath.Join(crtPath, fileName)
			relPath := strings.Replace(crtPath, docPath, "", -1)
			fc.path = filepath.Join("/", relPath, fileName)
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

// []interface{} IndexOf
func IndexOf(s []interface{}, oriVal string) int {
	for i, val := range s {
		if val == oriVal {
			return i
		}
	}
	return -1
}

// []int IntIndexOf
func IntIndexOf(s []int, oriVal int) int {
	for i, val := range s {
		if val == oriVal {
			return i
		}
	}
	return -1
}

// []string uniq
func StringUniq(s []string) []string {
	rs := []string{}
	m := map[string]int{}
	for _, val := range s {
		has := m[val]
		if has == 0 {
			rs = append(rs, val)
			m[val] = 1
		}
	}
	return rs
}

// 转换文件名配置为map
func getDocNames(docs []interface{}) map[string]string {
	tempTrasName := map[string]string{}
	for _, v := range docs {
		t := v.(map[string]interface{})
		k := t["name"].(string)
		val := t["trans"].(string)
		tempTrasName[k] = val
		dirOrder = append(dirOrder, k)
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
