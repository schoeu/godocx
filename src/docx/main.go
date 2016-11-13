package main

import (
	"fmt"
	"log"
	"net/http"
	//"html/template"
	"util"
	"regexp"
	"io/ioutil"
	"path/filepath"
	//"runtime"
)

var (
	index    = "readme.md"
	docPath     = "/Users/memee/Downloads/svn/ps-fe"
	docxConf = "./docx-conf.json"
	theme = "default"
	port = "8910"
)

// 路由容器
type regRoute struct {
    pattern string
    handler func (w http.ResponseWriter, r *http.Request)
}

var routes = []regRoute{}
var static func(http.ResponseWriter, *http.Request)

func main() {
	initial()
	http.HandleFunc("/", allRoutes)
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// 设置路由公共方法
func setRegRoute(p string, h func(http.ResponseWriter, *http.Request)) {
    routes = append(routes, regRoute{p, h})
}

// 预处理
func initial() {
	util.ReadDirRs()

	// 获取当前文件路径
	// _, runp, _, _ := runtime.Caller(1)
	// dirname := filepath.Dir(runp)

    setRegRoute(".+.md$", func(w http.ResponseWriter, r *http.Request) {
		rsHTML := util.GetRsHTML(filepath.Join(docPath, r.URL.Path))
        fmt.Fprintf(w, string(rsHTML))
    })

	docStatic := staticFn(docPath)
	setRegRoute(".+.[png|jpg|gif|js|css]$", docStatic)

	// staticFilePath := "../../themes/" + theme
	// themeStatic := staticFn(filepath.Join(dirname, staticFilePath))
	// setRegRoute(".+.png$", docStatic)
}

// 路由分发
func allRoutes(w http.ResponseWriter, r *http.Request) {
	isHit := false
    // 添加路由
    for _, v := range routes {
        reg, err := regexp.Compile(v.pattern)
        if err != nil {
            continue
        }
        if reg.MatchString(r.URL.Path) {
			isHit = true
            v.handler(w, r)
        }
    }

	if !isHit {
		// 静态文件路径
		//static(w, r)
	}
}

// 静态文件路由
func staticFn(parentPath string) func(http.ResponseWriter, *http.Request){
	return func (w http.ResponseWriter, r *http.Request) {
		fileP := r.URL.Path
		crtPath := filepath.Join(parentPath, fileP)

		fmt.Println("staticFn-> %s", crtPath)

		fileContent, err := ioutil.ReadFile(crtPath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, string(fileContent))
	}
}