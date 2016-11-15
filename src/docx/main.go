package main

import (
	// "fmt"
	"log"
	"net/http"
	//"html/template"
	"util"
	"regexp"
	// "io/ioutil"
	"path/filepath"
	//"runtime"
	//"io"
	"os"
	"fmt"
)

var (
	index    = "readme.md"
	docPath     = "/Users/memee/Downloads/svn/ps-fe"
	docxConf = "./docx-conf.json"
	theme = "default"
	port = "8910"
)

const (
	listDir = 0x0001
)

// 路由容器
type regRoute struct {
    pattern string
    handler func (w http.ResponseWriter, r *http.Request)
}

var routes = []regRoute{}

var static = staticServer("../../themes/" + theme + "/static")

func main() {
	initial()
	http.HandleFunc("/", allRoutes)
	//staticServer("static", "../../themes/" + theme, 0)
	

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
	//util.ReadDirRs()

	// 获取当前文件路径
	// _, runp, _, _ := runtime.Caller(1)
	// dirname := filepath.Dir(runp)

    setRegRoute(".+.md$", func(w http.ResponseWriter, r *http.Request) {
		rsHTML := util.GetRsHTML(filepath.Join(docPath, r.URL.Path))
		locals := make(map[string]interface{})
		locals["mdData"] = string(rsHTML)
        //io.WriteString(w, string(rsHTML))
		util.RenderTpl("../../themes/" + theme + "/views/main.tpl", locals, w)
    })
	fmt.Println(routes)
}

// 路由分发
func allRoutes(w http.ResponseWriter, r *http.Request) {

    // 添加路由
    for _, v := range routes {
        reg, err := regexp.Compile(v.pattern)
        if err != nil {
            continue
        }
        if reg.MatchString(r.URL.Path) {
            v.handler(w, r)
        }
    }

	static(w,r)
}

// 检测文件是否存在
func isExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return true
	}
	return os.IsExist(err)
}

// 静态文件服务器
// func staticServer(prefix string, staticDir string, flags int){
// 	http.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
// 		file := staticDir + r.URL.Path[len(prefix)-1:]
// 		if (flags & listDir) == 0 {
// 			if exists := isExists(file); !exists { 
// 				http.NotFound(w, r)
// 				return
// 			} 
// 		}
// 		http.ServeFile(w, r, file)
// 	})
// }


func staticServer(prefix string) func(http.ResponseWriter,*http.Request){
	return func (w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		file := filepath.Join(prefix, p)
		fmt.Println("file", file)
		http.ServeFile(w, r, file)
		return
	}
}