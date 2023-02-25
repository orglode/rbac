package api

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func InitLogger() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	file1, _ := os.Create("./log/info_" + time.Now().Format("20060102") + ".log")
	file2, _ := os.Create("./log/error_" + time.Now().Format("20060102") + ".log")

	core1 := zapcore.NewCore(encoder, zapcore.AddSync(file1), zapcore.DebugLevel)
	core2 := zapcore.NewCore(encoder, zapcore.AddSync(file2), zapcore.ErrorLevel)
	// 将两个core合并成一个新的core
	newcore := zapcore.NewTee(core1, core2)
	// 创建logger
	logger := zap.New(newcore)
	return logger
}
