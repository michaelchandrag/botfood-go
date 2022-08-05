package lib_logger

import (
	"os"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"github.com/google/uuid"
)

var Logger *zap.Logger

func init() {
	config := zap.NewProductionEncoderConfig()
    config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	filename := fmt.Sprintf("%s_%s.log", time.Now().Format("2006-01-02"), uuid.New())
	logDirectory := fmt.Sprintf("public/Logs/%s", filename)

	logFile, _ := os.OpenFile(logDirectory, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}