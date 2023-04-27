package logger

import "go.uber.org/zap"

var Log *zap.Logger

func InitLogger() {
	if Log == nil {
		Log, _ = zap.NewProduction()
	}
}

func Sync() {
	Log.Sync()
}
