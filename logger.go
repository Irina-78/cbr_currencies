package main

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const logFileName = "cbr_currencies.log"

func newProductionLogger() (*zap.Logger, error) {
	flags := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	file, err := os.OpenFile(logFileName, flags, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to create or open log file: %v", err)
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(newProductionConfig().EncoderConfig),
		zapcore.AddSync(file),
		zapcore.Level(zap.InfoLevel),
	)
	return zap.New(core), nil
}

func newDevelopmentLogger() (*zap.Logger, error) {
	return newDevelopmentConfig().Build()
}

func newProductionConfig() zap.Config {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05.000Z0700"))
	}
	return cfg
}

func newDevelopmentConfig() zap.Config {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	cfg.EncoderConfig.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}
	return cfg
}
