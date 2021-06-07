package utility

import (
	"os"

	"github.com/withmandala/go-log"
)

//Logger New
var Logger = NewLogger()

//NewLogger func
func NewLogger() *log.Logger {
	logger := log.New(os.Stderr).WithColor()
	return logger
}
