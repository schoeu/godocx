package main

import (
	"html/template"
	"net/http"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"conf"
	"search"
	"util"
	"zap"
)

var (
	index        = "/readme.md"
	theme        = "default"
	mdReg        = ".+.md$"
	staticPrefix = "static"
	staticRoot   = "../themes/" + theme

	// 配置文件变量
	docPath     = conf.DocxConf.GetJson("path").(string)
	port        = conf.DocxConf.GetJson("port").(string)
	supportInfo = conf.DocxConf.GetJson("supportInfo").(string)
	title       = conf.DocxConf.GetJson("title").(string)
	headText    = conf.DocxConf.GetJson("headText").(string)
	links       = conf.DocxConf.GetJson("extUrls.links").([]interface{})
	label       = conf.DocxConf.GetJson("extUrls.label").(string)
	zlog        = zap.GetLogger()
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
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		// log.Fatal("ListenAndServe: ", err)
	}
}

var navStr string

// 预处理
func initial() {

	// domtree 处理
	dirData := util.ReadDirRs()
	navStr = util.MakeNav(&dirData)

	search.PreProcess()

	http.HandleFunc("/api/update", updateRoutes)

	http.HandleFunc("/api/search", search.SearchRoutes)

	http.HandleFunc("/", allRoutes)
}

// markdown 文件处理
func mdHandler(mdRelPath string, w http.ResponseWriter, r *http.Request) {
	mdPath := filepath.Join(docPath, mdRelPath)
	content := util.GetRsHTML(mdPath)

	brandArr := util.GetPjaxContent(mdRelPath)

	// pjax branch
	isPjax := r.Header.Get("x-pjax") == "true"
	zlog.Info(mdPath)
	// 如果是pajx请求则返回片段，其他返回整模板
	if isPjax {
		brandPd := PageData{
			MdData:    template.HTML(content),
			BrandData: brandArr,
			HeadText:  headText,
		}
		util.RenderTpl(staticRoot+"/views/pjax.tmpl", brandPd, w)
	} else {
		pd := PageData{
			MdData:      template.HTML(content),
			NavData:     template.HTML(navStr),
			SupportInfo: supportInfo,
			Title:       title,
			HeadText:    headText,
			Links:       links,
			Label:       label,
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

func updateRoutes(w http.ResponseWriter, r *http.Request) {
	docPath = "./"
	cmd := exec.Command(docPath, "git status")
	cmd.Output()
	cmd.Start()
}
