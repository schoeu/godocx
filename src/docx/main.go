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

type PageData struct {
    mdData  string
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
func mdHandler(mdRelPath string, w http.ResponseWriter, r *http.Request) {
	mdPath := filepath.Join(docPath, mdRelPath)
	content := util.GetRsHTML(mdPath)
	
	//TODO pjax branch
	p := PageData{}
	p.mdData = string(content)
	fmt.Println(p)
	util.RenderTpl(staticRoot + "/views/main.tmpl", p, w)
}

// 路由分发
func allRoutes(w http.ResponseWriter, r *http.Request) {
	routePath := r.URL.Path
	isMd, _ :=regexp.MatchString(mdReg, routePath)
	fmt.Println(routePath)
	if routePath == "/" {
		mdHandler(index, w, r)
	} else if  isMd {
		fmt.Println("ismd", routePath)
		mdHandler(routePath, w, r)
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

