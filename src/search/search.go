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
    Content string `json:"content"`
}

var pathCtt = []string{}

type searchCtt []searchTitle

func SearchRoutes(w http.ResponseWriter, r *http.Request) {
    key = r.FormValue("name")
    setype := r.FormValue("type")
    if key != "" {
        if len(pathCtt) == 0 {
            filepath.Walk(docPath, walkFn)
        }
        searchFn(w, r, setype)
    }
}

// 标题搜索
func searchFn(w http.ResponseWriter, r *http.Request, setype string) {
    filterRs := collectRs(setype)
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
    // TODO 返回对象
    io.WriteString(w, string(jsonRs))
}

func collectRs(setype string) []searchTitle{
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
            if setype == "" {
                replacedCtt := searchContentFn(content, title)
                st.Content = replacedCtt
            }

            stt = append(stt, st)
        }
    }
    return stt
}

// 内容搜索
func searchContentFn(content []byte, title string) string{
    keyRe := regexp.MustCompile(key)
    // TODO 
    replacedContent := keyRe.ReplaceAllString(string(content), "<span class='hljs-string'>$0</span>")
    return replacedContent
}