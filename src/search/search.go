package search

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"util"
	"github.com/huichen/sego"
	"fmt"
)

var (
	mdReg    = ".+.md$"
	key      = ""
	width    = 333
	deep     = 3
	max      = 100
	imgRe    = regexp.MustCompile("<img.*?>")
	headRe   = regexp.MustCompile("<[h1|h2|h3|h4|h5|h6]>.*?</[h1|h2|h3|h4|h5|h6]>")
	strongRe = regexp.MustCompile("strong")
)

var keySl []string

type searchTitle struct {
	Path    string `json:"path"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type searchCtt []searchTitle

    
func SearchRoutes(w http.ResponseWriter, r *http.Request) {
	key = r.FormValue("name")
	if UsePinyin {
		text := []byte(key)
		segments := Segmenter.Segment(text)
		key = sego.SegmentsToString(segments, false)
		keySl = sego.SegmentsToSlice(segments, false)
		key = strings.Join(keySl, "|")
	}
	
	setype := r.FormValue("type")
	if key != "" {
		if len(PathCtt) == 0 {
			PreProcess()
		}
		searchFn(w, r, setype)
	}
}

// 标题搜索
func searchFn(w http.ResponseWriter, r *http.Request, setype string) {
	filterRs := collectRs(setype)
	returnJSON(filterRs, w, r)
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

func collectRs(setype string) []searchTitle {
	// 替换标题
	keyRe := regexp.MustCompile(key)
	var stt searchCtt
	var titleMatched searchCtt
	var rsCt searchCtt
	var ttTemp searchCtt
	for i, v := range PathCtt {
		if i < max {
			content := util.GetConent(v)
			ext := filepath.Ext(v)
			title := util.GetTitle(ext, content)

			var st = searchTitle{}
			st.Path = strings.Replace(v, DocPath, "", -1)
			ok, _ := regexp.MatchString(key, title)

			replacedTitle := keyRe.ReplaceAllString(title, "<span class='hljs-string'>$0</span>")
			st.Title = replacedTitle

			// 标题为空则跳过
			if st.Title == "" {
				continue
			}

			// 内容搜索
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
				// 标题搜索
				if ok {
					// 标题拼音搜索
					if UsePinyin {
						for _, v := range ScArr {
							//isGet := false
							title := v.title
							spell := v.spell
							pos := v.pos
							initials := v.initials
							ems := keySl[:]
							sIdx := strings.Index(spell, key)
							if sIdx > -1 {
								pIdx := util.IntIndexOf(pos, sIdx)
								if pIdx > -1 {
									wordCount := 0
									for i = pIdx; i < len(pos); i++ {
										if sIdx + len(key) <= pos[i] {
											break
										}
										wordCount ++
									}
									fmt.Println("~~~~~~~", title, len(title), pIdx, wordCount)
									sele := title[pIdx:wordCount]
									ems = append(ems, sele)
									//isGet = true
								}
							}
							
							// initials检索
							iIdex := strings.Index(initials, key)
							if iIdex > -1 {
								iele := title[iIdex:len(key)]
								ems = append(ems, iele)
								//isGet = true
							}

							// 去重
							ems = util.StringUniq(ems)
							emkeys := strings.Join(ems, " ")
							r := regexp.MustCompile("\\s+")
							s := r.ReplaceAllString(emkeys, "|")
							reg := regexp.MustCompile("^(\\|)*|(\\|)*$")
							rsString := reg.ReplaceAllString(s, "")

							tReg := regexp.MustCompile(rsString)
							matchTitle := tReg.MatchString(title)

							if /* isGet || */matchTitle {
								rpTitle := tReg.ReplaceAllString(title, "<span class='hljs-string'>$0</span>")
								st.Title = rpTitle
							}
						}
					}
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

// 标题搜索
func searchTitleFn() {
	// 载入词典
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
		if crtDp < deep {
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
			crtDp++
		}
	}
	return strings.Join(matchedContent, "...")
}
