package main

import (
	"conf"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"util"
)

var (
	index        = "/readme.md"
	docPath      = "/Users/memee/Downloads/svn/ps-fe"
	theme        = "default"
	mdReg        = ".+.md$"
	staticPrefix = "static"
	staticRoot   = "../../themes/" + theme

	// 配置文件变量
	port        = conf.DocxConf.GetJson("port")
	supportInfo = conf.DocxConf.GetJson("supportInfo")
	title       = conf.DocxConf.GetJson("title")
	headText    = conf.DocxConf.GetJson("headText")
	links       = conf.DocxConf.GetJson("extUrls.links")
	label       = conf.DocxConf.GetJson("extUrls.label")
)

type PageData struct {
	MdData      template.HTML
	NavData     template.HTML
	BrandData   []string
	SupportInfo string
	Title       string
	HeadText    string
	Links       []interface{}
	Label       string
}

// 入口函数
func main() {
	// godocx 初始化
	initial()

	// 监听端口
	err := http.ListenAndServe(":"+port.(string), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

var navStr string

// 预处理
func initial() {

	// domtree 处理
	dirData := util.ReadDirRs()
	navStr = util.MakeNav(&dirData)

	http.HandleFunc("/", allRoutes)
}

// markdown 文件处理
func mdHandler(mdRelPath string, w http.ResponseWriter, r *http.Request) {
	mdPath := filepath.Join(docPath, mdRelPath)
	content := util.GetRsHTML(mdPath)

	brandArr := util.GetPjaxContent(mdRelPath)

	// pjax branch
	isPjax := r.Header.Get("x-pjax") == "true"
	// 如果是pajx请求则返回片段，其他返回整模板
	if isPjax {
		//fmt.Fprintf(w, string(content))
		brandPd := PageData{
			MdData:    template.HTML(content),
			BrandData: brandArr,
			HeadText:  headText.(string),
		}
		util.RenderTpl(staticRoot+"/views/pjax.tmpl", brandPd, w)
	} else {
		pd := PageData{
			MdData:      template.HTML(content),
			NavData:     template.HTML(navStr),
			SupportInfo: supportInfo.(string),
			Title:       title.(string),
			HeadText:    headText.(string),
			Links:       links.([]interface{}),
			Label:       label.(string),
			BrandData:   brandArr,
		}
		util.RenderTpl(staticRoot+"/views/main.tmpl", pd, w)

	}
}

// 路由分发
func allRoutes(w http.ResponseWriter, r *http.Request) {
	routePath := r.URL.Path
	isMd, _ := regexp.MatchString(mdReg, routePath)
	if routePath == "/" {
		mdHandler(index, w, r)
	} else if isMd {
		mdHandler(routePath, w, r)
	} else {
		staticServer(w, r)
	}
}

// 静态文件服务
func staticServer(w http.ResponseWriter, r *http.Request) {
	var staticRou string
	p := r.URL.Path
	pathSp := strings.Split(p, "/")
	if pathSp[1] == staticPrefix {
		staticRou = filepath.Join(staticRoot, p)
	} else {
		staticRou = filepath.Join(docPath, p)
	}

	http.ServeFile(w, r, staticRou)
}
