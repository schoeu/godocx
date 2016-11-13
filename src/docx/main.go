package main

import (
	"fmt"
	"log"
	"net/http"
	//"html/template"
	"walker"
	"regexp"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

var (
	index    = "readme.md"
	docPath     = "/Users/memee/Downloads/svn/ps-fe"
	docxConf = "./docx-conf.json"
	theme = "default"
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
	err := http.ListenAndServe(":8911", nil)
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
	// walker.ReadDirRs()

    setRegRoute(".+.md$", func(w http.ResponseWriter, r *http.Request) {
		rsHTML := walker.GetRsHTML(r.URL.Path)
        fmt.Fprintf(w, string(rsHTML))
    })

	// 获取当前文件路径
	_, runp, _, _ := runtime.Caller(1)
	dirname := filepath.Dir(runp)

	// staticFilePath := "../../themes/" + theme
	//static = staticFn(filepath.Join(dirname, staticFilePath))
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