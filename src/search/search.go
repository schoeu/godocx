package search

import (
    "path/filepath"
    "net/http"
    "os"
    "fmt"
    "regexp"

    "conf"
    "util"
)

var (
    docPath = conf.DocxConf.GetJson("path").(string)
    mdReg   = ".+.md$"
    key = ""
)

type searchTitle struct {
    path string
    title string
}

var pathCtt = []string{}

type searchCtt []searchTitle

func SearchRoutes(w http.ResponseWriter, r *http.Request) {
    key = r.FormValue("name")

    if key != "" {
        search()
    }
}

// 搜索
func search() {
    filepath.Walk(docPath, walkFn)
    filterRs := collectRs()
    fmt.Println(filterRs)
}

func walkFn (walkPath string, info os.FileInfo, err error) error{
    isMd, err := regexp.MatchString(mdReg, walkPath)
    if !info.IsDir() {
        if err == nil && isMd {
            pathCtt = append(pathCtt, walkPath)
        }
    }
    return nil
}


func returnTpl() {
    //util.RenderTpl(staticRoot+"/views/pjax.tmpl", brandPd, w)
}

func collectRs() []searchTitle{
    // 替换标题
    keyRe := regexp.MustCompile(key)
    var stt searchCtt
    for _, v := range pathCtt {
        content := util.GetConent(v)
        ext := filepath.Ext(v)
        title := util.GetTitle(ext, content)

        ok, _ := regexp.MatchString(key, title)
        if  ok {
            replacedTitle := keyRe.ReplaceAllString(title, "<span class='hljs-string'>$0</span>")
            stt = append(stt, searchTitle{
                path: v,
                title: replacedTitle,
            })
        }
    }
    return stt
}