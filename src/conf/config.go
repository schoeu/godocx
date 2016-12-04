package conf

import (
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
)

var (
	configPath = ""
	defaultP   = "../docx-conf.json"
)

type Config struct {
	path    string
	content string
}

var DocxConf = &Config{path: configPath}

// 解析配置文件路径
func cmdParse() string {
	arg := ""
	args := os.Args
	if len(args) > 1 {
		arg = os.Args[1]
	}
	if arg == "" {
		arg = defaultP
	}
	return arg
}

// 获取配置文件
func (c *Config) getConf() {
	if c.content == "" {
		// 获取配置文件路径
		if configPath = cmdParse(); configPath != "" {
			c.path = configPath
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
