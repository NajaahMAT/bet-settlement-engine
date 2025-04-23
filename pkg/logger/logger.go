package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger() {
	var err error
	Log, err = zap.NewProduction() // or zap.NewDevelopment()
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}
}
