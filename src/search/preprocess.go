package search

import (
    "path/filepath"
    "os"
    "regexp"
    "conf"
    "util"
    "strings"

    "github.com/mozillazg/go-pinyin"
)

var (
    DocPath  = conf.DocxConf.GetJson("path").(string)
    PathCtt = []string{}
)

type searchContent struct {
    title string
    pos []int
    path string
    spell string
    initials string
}

var ScArr = []searchContent{}
var pyArgs = pinyin.NewArgs()

func PreProcess() {
    filepath.Walk(DocPath, walkFn)
    getTitleInfo()
}

func walkFn(walkPath string, info os.FileInfo, err error) error {
	isMd, err := regexp.MatchString(mdReg, walkPath)
	if !info.IsDir() {
		if err == nil && isMd {
			PathCtt = append(PathCtt, walkPath)
		}
	}
	return nil
}

// 收集标题搜索数据
func getTitleInfo() {
    for _, v := range PathCtt {
        sc := searchContent{}
        pos := []int{}
        idx := 0
        content := util.GetConent(v)
        ext := filepath.Ext(v)
        title := util.GetTitle(ext, content)

        pyArgs.Style = pinyin.FIRST_LETTER
        sc.path = strings.Replace(v, DocPath, "", -1)
        sc.initials = strings.Join(pinyin.LazyPinyin(title, pyArgs), "")
        spell := pinyin.LazyConvert(title, nil)
        for _, l := range spell {
            pos = append(pos, idx)
            idx += len(l)
        }
        sc.spell = strings.Join(spell, "")
        sc.pos = pos
        sc.title = util.GetTitle(ext, content)

        ScArr = append(ScArr, sc)
    }
}