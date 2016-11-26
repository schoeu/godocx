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
    width = 360
    deep = 3
    max = 100
    imgRe = regexp.MustCompile(`/<img.*?>/`)
    headRe = regexp.MustCompile(`/<img.*?>/`)
    strongRe = regexp.MustCompile(`/strong/`)
)

type searchTitle struct {
    Path string `json:"path"`
    Title string `json:"title"`
    Content string `json:"content"`
}

var pathCtt = []string{}
var matchedContent = []string{}

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
    for i, v := range pathCtt {
        if i < max {
            content := util.GetConent(v)
            ext := filepath.Ext(v)
            title := util.GetTitle(ext, content)

            var st = searchTitle{}
            ok, _ := regexp.MatchString(key, title)
            if  ok {
                replacedTitle := keyRe.ReplaceAllString(title, "<span class='hljs-string'>$0</span>")
                if len(replacedTitle) > 0 && len(st.Path) > 0{
                    st.Path = strings.Replace(v, docPath, "", -1)
                    st.Title = replacedTitle
                    titleMatched = append(titleMatched, st)
                }
            }
            if setype == "" {
                replacedCtt := searchContentFn(string(content), title)
                if len(replacedCtt) > 0 {
                    st.Content = replacedCtt
                    stt = append(stt, st)
                }
            }
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
    idxArr := keyRe.FindAllStringIndex(content, -1)
    contentLength := len(content)
    crtDp := 0
    replacedContent := keyRe.ReplaceAllString(content, "<span class='hljs-string'>$0</span>")
    replacedContent = imgRe.ReplaceAllString(replacedContent, "")
    replacedContent = headRe.ReplaceAllString(replacedContent, "")
    replacedContent = strongRe.ReplaceAllString(replacedContent, "")
    for _, v := range idxArr {
        if crtDp < deep{
            start := v[0] - width
            if start < 0 {
                start = 0
            }
            end := v[1] + width
            if end > contentLength {
                end = contentLength
            }
            cutPart := replacedContent[start: end]
            matchedContent = append(matchedContent, cutPart)
            crtDp ++
        }
        
    }
    return strings.Join(matchedContent, "...")
}