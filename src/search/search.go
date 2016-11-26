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
    "fmt"

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
    var titleMatched searchCtt
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
            titleMatched = append(titleMatched, st)
        }
        if setype == "" {
            replacedCtt := searchContentFn(string(content), title)
            st.Content = replacedCtt
            stt = append(stt, st)
        }
        
    }

    //fmt.Println(strings.Index("chicken", "ken"))
    //fmt.Println(strings.Index("我是测试啊水电费", "测"))
    // 标题匹配优先内容匹配
    return append(titleMatched, stt...)
}

// 内容搜索
func searchContentFn(content, title string) string{
    keyRe := regexp.MustCompile(key)
    // idx := keyRe.FindAllStringIndex(content, -1)
    // TODO 
    idxArr := keyRe.FindAllStringIndex(content, -1)
    // for _, v := range idxArr {
    //     fmt.Println(v[0])
    // }
    fmt.Println(idxArr)

    replacedContent := keyRe.ReplaceAllString(content, "<span class='hljs-string'>$0</span>")
    
    return replacedContent
}