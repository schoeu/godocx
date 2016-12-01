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
    width = 333
    deep = 3
    max = 100
    imgRe = regexp.MustCompile("<img.*?>")
    headRe = regexp.MustCompile("<[h1|h2|h3|h4|h5|h6]>.*?</[h1|h2|h3|h4|h5|h6]>")
    strongRe = regexp.MustCompile("strong")
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
    fmt.Println("")
}

func collectRs(setype string) []searchTitle{
    // 替换标题
    keyRe := regexp.MustCompile(key)
    var stt searchCtt
    var titleMatched searchCtt
    var rsCt searchCtt
    var ttTemp searchCtt
    for i, v := range pathCtt {
        if i < max {
            content := util.GetConent(v)
            ext := filepath.Ext(v)
            title := util.GetTitle(ext, content)

            var st = searchTitle{}
            st.Path = strings.Replace(v, docPath, "", -1)
            ok, _ := regexp.MatchString(key, title)
            
            replacedTitle := keyRe.ReplaceAllString(title, "<span class='hljs-string'>$0</span>")
            st.Title = replacedTitle

            // 标题为空则跳过
            if st.Title == "" {
                continue
            }

            if setype == "" {
                replacedCtt := searchContentFn(string(content))
                if len(replacedCtt) > 0 {
                    st.Content = replacedCtt
                    // 标题匹配优先级高于内容匹配
                    if ok {
                        ttTemp = append(ttTemp, st)
                    } else {
                        stt = append(stt, st)
                    }
                }
            } else {
                if  ok {
                    titleMatched = append(titleMatched, st)
                }
            } 
        }
    }

    if setype == "title" {
        rsCt = titleMatched
    } else {
        rsCt = append(ttTemp, stt...)
    }

    return rsCt
}

// 内容搜索
func searchContentFn(content string) string {
    var matchedContent = []string{}
    keyRe := regexp.MustCompile(key)
    idxArr := keyRe.FindAllStringIndex(content, -1)
    crtDp := 0
    rlc := keyRe.ReplaceAllString(content, "<span class='hljs-string'>$0</span>")
    contentLength := len(rlc)
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
            cutPart := rlc[start:end]
            cutPart = imgRe.ReplaceAllString(cutPart, "")
            cutPart = headRe.ReplaceAllString(cutPart, "")
            cutPart = strongRe.ReplaceAllString(cutPart, "")
            matchedContent = append(matchedContent, cutPart)
            crtDp ++
        }
    }
    return strings.Join(matchedContent, "...")
}
