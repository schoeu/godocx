package conf

import (
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
)

var (
	configPath = "./docx-conf.json"
)

type Config struct {
	path    string
	content string
}

var DocxConf = &Config{path: configPath}

// 获取配置文件
func (c *Config) getConf() {
	if c.content == "" {
		if leng := len(os.Args); leng > 1 && os.Args[1] != "" {
			c.path = os.Args[1]
		}
		config, err := ioutil.ReadFile(c.path)
		if err == nil {
			c.content = string(config)
		}
	}
}

// 获取参数配置
func (c *Config) GetJson(param string) interface{} {
	c.getConf()
	return gjson.Get(c.content, param).Value()
}
