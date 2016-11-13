package main

import (
	"fmt"
	"log"
	"net/http"
	//"html/template"
	"walker"
	"regexp"
)

var (
	index    = "readme.md"
	path     = "/Users/memee/Downloads/svn/ps-fe"
	docxConf = "./docx-conf.json"
)

// 路由容器
type regRoute struct {
    pattern string
    handler func (w http.ResponseWriter, r *http.Request)
}

var routes = []regRoute{}

func main() {

	initial()

	http.HandleFunc("/", allRoutes)

	// 路由匹配&文档树
	//http.HandleFunc("/a", allHandle)

	// walker.ReadDirRs()

	err := http.ListenAndServe(":8911", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}



// 路由分发
func allRoutes(w http.ResponseWriter, r *http.Request) {
    // 添加路由
	fmt.Println("routes",routes)
    for _, v := range routes {
        reg, err := regexp.Compile(v.pattern)
		fmt.Println(`reg`,reg)
        if err != nil {
            continue
        }
        if reg.MatchString(r.URL.Path) {
			fmt.Println(v.handler)
            v.handler(w, r)
        }
    }
}

// 设置路由公共方法
func setRegRoute(p string, h func(http.ResponseWriter, *http.Request)) {
    routes = append(routes, regRoute{p, h})
}

// 添加路由
func initial() {

	// routes = append(routes, regRoute{"*.md", func(w http.ResponseWriter, r *http.Request) {
    //     fmt.Fprintf(w, "*.md")
    // }})

    setRegRoute(".+.md$", func(w http.ResponseWriter, r *http.Request) {
		rsHTML := walker.GetRsHTML(r.URL.Path)
        fmt.Fprintf(w, string(rsHTML))
    })
}
