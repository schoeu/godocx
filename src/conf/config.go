package conf

import(
	"io/ioutil"
	"github.com/tidwall/gjson"
)

var (
	configPath = "../../docx-conf.json"
)

type Config struct {
	path string
	content string
}

var DocxConf = &Config{path: configPath}

// 获取配置文件
func (c *Config)getConf(){
	if c.content == "" {
		config, err := ioutil.ReadFile(c.path)
		if err == nil {
			c.content = string(config)
		}
	}
}

// 获取参数配置
func (c *Config)GetJson(param string) interface{}{
	c.getConf()
	return gjson.Get(c.content, param)
}