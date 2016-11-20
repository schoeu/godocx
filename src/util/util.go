package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"os"
	//"strings"
	"github.com/russross/blackfriday"
	"html/template"
	"net/http"
)
var (
	docPath    = "/Users/memee/Downloads/svn/ps-fe-new"
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
	child *[]fileCache
}

var docTree = make([]fileCache,0)

type fileTrasName  map[string]string

var fileNameMap = fileTrasName{}

func ReadDirRs() {
	fileNameMap = fileTrasName{
		"aladdin":"阿拉丁",
		"www":"搜索结果页",
		"standards":"规范流程",
		"superframe":"superframe",
		"log":"日志",
		"transcode":"无线转码",
		"realtime":"时效性",
		"performance":"性能优化",
		"references":"资源引入",
		"data":"数据接口",
		"tools":"工具服务",
		"xueshu":"学术",
		"advertise":"广告",
		"rules":"规范",
		"santa":"圣玛利亚",
		"cardspeedup":"模板性能优化",
		"commonupdate":"通用升级",
		"develop":"圣玛利亚",
		"jscommmon":"js组件",
		"platform":"平台指南",
		"publish":"上线",
		"standard":"开发规范",
		"technicalarea":"技术专区",
		"tongji":"日志",
		"wise":"无线网页搜索",
		"devdocs":"开发指导",
		"static":"静态文件",
		"pc":"PC网页搜索",
		"midpage":"搜索中间页",
		"show":"展现日志",
		"click":"点击日志",
		"client":"客户端相关规范",
		"process":"使用和变更流程",
		"action":"异步日志",
		"framework":"架构",
		"aladdin-debug":"阿拉丁常见调试",
		"research":"技术调研",
		"courseware":"串讲文档",
		"grid":"栅格化",
		"general-dev":"通用开发相关",
		"frontend-dev":"前端开发相关",
		"environment-dev":"环境相关",
		"union":"联盟相关",
		"front-interface":"前后端接口",
		"refactor":"重构",
		"test":"测试相关",
		"async":"异步化",
		"spec":"规范",
		"component":"功能组件",
		"sample":"抽样相关",
		"schema":"Schemas标记",
		"aladdin-test":"阿拉丁（测试）",
		"tpldev":"开发平台",
		"wireless-dev":"无线开发",
		"new-reading":"新人必读",
		"pc-doc":"PC开发文档",
		"team":"团队介绍",
		"pc-other":"PC其他开发文档",
		"0-send-guide":"发送指南",
		"1-stat-guide":"统计指南",
		"todolist":"备忘列表",
		"pcspans":"PCspans",
		"pcuijs":"PC-js组件",
		"paduijs":"PAD组件",
		"pclog":"PC日志",    
	}

	makeDomTree(docPath, &docTree)
	// sep := string(filepath.Separator)
	// docPathReg := regexp.MustCompile(docPath)
	// filepath.Walk(docPath, func (fPath string, info os.FileInfo, err error) error {
	// 	relPath := docPathReg.ReplaceAllString(fPath, "")
	// 	justDir := filepath.Dir(relPath)
	// 	pathArr := strings.Split(justDir, sep)
	// 	fmt.Println(relPath,"~~~~",pathArr)
	// 	return nil
	// })

	fmt.Println(docTree, cap(docTree), len(docTree))
}

func makeDomTree (crtPath string, ctt *[]fileCache) {
	files, err := ioutil.ReadDir(crtPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fc := fileCache{}
		isDir := file.IsDir()
		docName := file.Name()
		fmt.Println("file",file)
		// 文件夹需要递归处理，文件则直接存容器
		if isDir {
			fileName := file.Name()
			//hitted := indexOf(ignoreDirs, fileName)
			//if hitted > -1 {
				subFileCache := make([]fileCache, 0)
				fc.title = fileNameMap[fileName]
				fc.path = fileName
				fc.ty = "dir"
				fc.child = &subFileCache
				// ctt := append(*ctt, fc)
				*ctt = append(*ctt, fc)
				
				makeDomTree(filepath.Join(crtPath, file.Name()), fc.child)
			//}
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
					//_ := append(*ctt, fc)
					*ctt = append(*ctt, fc)
				}
			}
		}
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
	fmt.Println("xxx")
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

// []string indexOf
func indexOf(s []string, oriVal string) int{
	for idx, val := range s {
		if val == oriVal {
			return idx
		}
	}
	return -1
}