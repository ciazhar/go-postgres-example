package logger

import (
	"fmt"
	"go.uber.org/zap"
	"sync"
)

var once sync.Once

var Logger *zap.Logger

// Init returns an instance of logger
func Init() {

	once.Do(func() {
		config := Configuration{
			EnableConsole:     true,
			ConsoleLevel:      DebugD,
			ConsoleJSONFormat: true,
			EnableFile:        true,
			FileLevel:         InfoD,
			FileJSONFormat:    true,
			FileLocation:      "logs/app.log",
		}

		var err error
		Logger, err = newZapLogger(config)
		if err != nil {
			panic(err.Error())
		}
	})
}

func Info(format string, a ...interface{}) {
	Logger.Info(fmt.Sprintf(format, a))
}

func Error(format string) {
	Logger.Error(format)
}
