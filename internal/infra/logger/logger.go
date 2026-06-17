package logger

import (
	"os"
	"strings"

	"github.com/shuyou-ai/shuyou-go/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New(cfg config.LogConfig) (*zap.Logger, error) {
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(strings.ToLower(cfg.Level))); err != nil {
		level.SetLevel(zap.InfoLevel)
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	if strings.ToLower(cfg.Format) == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	var writeSyncer zapcore.WriteSyncer
	if strings.ToLower(cfg.Output) == "file" {
		writeSyncer = zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.File.Filename,
			MaxSize:    cfg.File.MaxSize,
			MaxBackups: cfg.File.MaxBackups,
			MaxAge:     cfg.File.MaxAge,
			Compress:   cfg.File.Compress,
		})
	} else {
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	core := zapcore.NewCore(encoder, writeSyncer, level)
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)), nil
}
