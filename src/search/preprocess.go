package search

import (
	"conf"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
	"util"
	"fmt"

	"github.com/huichen/sego"
	"github.com/mozillazg/go-pinyin"
)

var (
	DocPath   = conf.DocxConf.GetJson("path").(string)
	UsePinyin = conf.DocxConf.GetJson("usePinyin").(bool)
	PathCtt   = []string{}
	DictPath  = "/Users/memee/Downloads/svn/git/go/godocx/dict/dictionary.txt"
)

type searchContent struct {
	title    string
	pos      []int
	path     string
	spell    string
	initials string
}

var ScArr = []searchContent{}
var pyArgs = pinyin.NewArgs()
var spellArgs = pinyin.NewArgs()

// 载入词典
var Segmenter sego.Segmenter

func PreProcess() {
	if UsePinyin {
		Segmenter.LoadDictionary(DictPath)
	}

	filepath.Walk(DocPath, walkFn)
	getTitleInfo()
	fmt.Println(ScArr)
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
		pyArgs.Fallback = func(r rune, a pinyin.Args) []string {
			return []string{string(r)}
		}
		sc.path = strings.Replace(v, DocPath, "", -1)
		sc.initials = strings.Join(pinyin.LazyPinyin(title, pyArgs), "")
		spellArgs.Fallback = func(r rune, a pinyin.Args) []string {
			return []string{string(r)}
		}
		spell := pinyin.LazyPinyin(title, spellArgs)
		for _, l := range spell {
			pos = append(pos, idx)
			idx += utf8.RuneCountInString(l)
		}
		sc.spell = strings.Join(spell, "")
		sc.pos = pos
		sc.title = util.GetTitle(ext, content)

		ScArr = append(ScArr, sc)
	}
}
