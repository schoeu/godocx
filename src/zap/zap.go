package zap

import (
	"github.com/uber-go/zap"
	"io/ioutil"

	"conf"
)

var Logger zap.Logger

func GetLogger() zap.Logger {
	if Logger == nil {
		logPath := conf.DocxConf.GetJson("logPath").(string)

		f, err := ioutil.TempFile(logPath, "log")
		if err != nil {
			panic("failed to create temporary file")
		}

		Logger = zap.New(
			zap.NewJSONEncoder(),
			zap.Output(f),
		)
	}

	return Logger
}
