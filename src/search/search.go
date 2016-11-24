package search

import (
    "path/filepath"
    "net/http"
    "os"
    "regexp"
    "encoding/json"
    "log"
    "io"
    "strings"

    "conf"
    "util"
)

var (
    docPath = conf.DocxConf.GetJson("path").(string)
    mdReg   = ".+.md$"
    key = ""
)

type searchTitle struct {
    Path string `json:"path"`
    Title string `json:"title"`
}

// type searchTitle map[string]string

var pathCtt = []string{}

type searchCtt []searchTitle

func SearchRoutes(w http.ResponseWriter, r *http.Request) {
    key = r.FormValue("name")

    if key != "" {
        search(w, r)
    }
}

// 搜索
func search(w http.ResponseWriter, r *http.Request) {
    if len(pathCtt) == 0 {
        filepath.Walk(docPath, walkFn)
    }
    
    filterRs := collectRs()
    returnJSON(filterRs, w, r)
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

// 返回搜索结果
func returnJSON(js []searchTitle, w http.ResponseWriter, r *http.Request) {
    jsonRs, err := json.Marshal(js)
	if err != nil {
		log.Fatal(err)
	}
    io.WriteString(w, string(jsonRs))
}

func collectRs() []searchTitle{
    // 替换标题
    keyRe := regexp.MustCompile(key)
    var stt searchCtt
    for _, v := range pathCtt {
        content := util.GetConent(v)
        ext := filepath.Ext(v)
        title := util.GetTitle(ext, content)


        var st = searchTitle{}
        ok, _ := regexp.MatchString(key, title)
        if  ok {
            replacedTitle := keyRe.ReplaceAllString(title, "<span class='hljs-string'>$0</span>")
            st.Path = strings.Replace(v, docPath, "", -1)
            st.Title = replacedTitle
            stt = append(stt, st)
        }
    }
    return stt
}