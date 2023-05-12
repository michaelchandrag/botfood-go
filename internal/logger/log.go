package logger

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AppLogger struct {
	LocalAgent *zap.Logger
}

var Agent AppLogger

func InitLogger() {
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
	localAgent := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	Agent = AppLogger{
		LocalAgent: localAgent,
	}
}

func getCaller() (filename string, fn string) {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		panic("Could not get context info for logger!")
	}
	filenameSplit := strings.Split(file, "/")
	targetfilename := file
	if len(filenameSplit) >= 5 {
		targetfilename = ""
		for i := 5; i < len(filenameSplit); i++ {
			targetfilename = targetfilename + "/" + filenameSplit[i]
		}
	}

	filename = targetfilename + ":" + strconv.Itoa(line)
	funcname := runtime.FuncForPC(pc).Name()
	fn = funcname[strings.LastIndex(funcname, ".")+1:]
	return filename, fn
}

func (log *AppLogger) Info(message string) {
	filename, fn := getCaller()
	log.LocalAgent.Info(message,
		zap.String("filename", filename),
		zap.String("function", fn),
	)
}

func (log *AppLogger) Error(message string) {
	filename, fn := getCaller()
	log.LocalAgent.Error(message,
		zap.String("filename", filename),
		zap.String("function", fn),
	)
}
