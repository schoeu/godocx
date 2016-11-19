package main

import (
	"log"
	"net/http"
	"regexp"
	"util"
	"path/filepath"
	"fmt"
	"os"
)

var (
	index    = "/readme.md"
	docPath  = "/Users/memee/Downloads/svn/ps-fe"
	docxConf = "./docx-conf.json"
	theme    = "default"
	port     = "8910"
	mdReg = ".+.md$"
	staticRoot = "../../themes/" + theme
)

// 路由容器
type regRoute struct {
	pattern string
	handler func(w http.ResponseWriter, r *http.Request)
}

var routes = []regRoute{}

var static = staticServer(staticRoot + "/static")

func main() {
	initial()

	err := http.ListenAndServe(":"+port, nil)
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
	//util.ReadDirRs()
	// 获取当前文件路径
	// _, runp, _, _ := runtime.Caller(1)
	// dirname := filepath.Dir(runp)
	http.HandleFunc("/", allRoutes)
	// setRegRoute(".+.md$", mdHandler)

}

// markdown 文件处理
func mdHandler(mdRelPath string, w http.ResponseWriter) {
	mdPath := filepath.Join(docPath, mdRelPath)
	content := util.GetRsHTML(mdPath)
	locals := make(map[string]interface{})
	locals["mdData"] = string(content)
	util.RenderTpl(staticRoot + "/views/main.tmpl", locals, w)
}

// 路由分发
func allRoutes(w http.ResponseWriter, r *http.Request) {
	routePath := r.URL.Path
	fmt.Println(routePath)
	if routePath == "/" {
		mdHandler(index, w)
	} else if isMd, err :=regexp.MatchString(mdReg, routePath); err != nil && isMd {
		mdHandler(routePath, w)
	} else {
		static(w, r)
	}
}

// 检测文件是否存在
func isExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return true
	}
	return os.IsExist(err)
}

func staticServer(prefix string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		file := filepath.Join(prefix, p)
		http.ServeFile(w, r, file)
		return
	}
}