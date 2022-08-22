package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Log     *zap.Logger
	Console *zap.Logger
)

func StartLogger() {
	fileLogger := &lumberjack.Logger{
		Filename: "./logs/logging.log",
		MaxSize:  50,
		MaxAge:   7,
	}

	fileCfg := zap.NewProductionEncoderConfig()
	fileCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(fileCfg)

	consoleCfg := zap.NewProductionEncoderConfig()
	consoleCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	consoleEncoder := zapcore.NewJSONEncoder(consoleCfg)

	Log = zap.New(zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(fileLogger), zap.LevelEnablerFunc(func(_ zapcore.Level) bool {
			return true
		})),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.ErrorLevel
		})),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stderr), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})),
	))

	Console = zap.New(zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.ErrorLevel
		})),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stderr), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})),
	))
}
