package conf

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Site     siteConfig   `yaml:"site"`
	Server   serverConfig `yaml:"server"`
	ConfPath string
}

type serverConfig struct {
	port int    `yaml:port`
	path string `yaml:path`
}

type siteConfig struct {
	headText    string       `yaml:"headText"`
	title       string       `yaml:"title"`
	logPath     int          `yaml:"logPath"`
	//ignoreDir   []string     `yaml:"ignoreDir"`
	supportInfo string       `yaml:"supportInfo"`
	//extUrls     extUrlStruct `yaml:"extUrls"`
}

type extUrlStruct struct {
	links linksStruct `yaml:"links"`
	label string      `yaml:"label"`
}

type linksStruct struct {
	name string `yaml:"name"`
	url  string `yaml:"url"`
}

func (c *Config) parseConfig() error {
	data, err := ioutil.ReadFile(filepath.Join(c.ConfPath, "conf.yml"))
	if err != nil {
		return fmt.Errorf("config err: %v", err)
	}

	fmt.Println(yaml.Unmarshal(data, &c))

	if err = yaml.Unmarshal(data, &c); err != nil {
		return fmt.Errorf("Unmarshal config: %v", err)
	}

	return nil
}

func GetConf() (*Config, error) {
	var config = &Config{
		ConfPath: ".",
	}

	err := config.parseConfig()
	if err != nil {
		return nil, fmt.Errorf("parse config: %v", err)
	}
	return config, nil
}
