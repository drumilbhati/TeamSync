package logs

import (
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func InitLogger() {
	config := zap.Config{
		Level:         zap.NewAtomicLevelAt(zap.DebugLevel),
		Encoding:      "json",
		OutputPaths:   []string{"stdout", "logs/app.log"},
		EncoderConfig: zap.NewProductionEncoderConfig(),
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	Log = logger.Sugar()
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
