package utility

import (
	"fmt"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

//Logger New
var Logger = NewLogger()

//NewLogger func
func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetReportCaller(true) // show file name and line number
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}

	return logger
}
