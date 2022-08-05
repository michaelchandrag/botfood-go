package lib_logger

import (
	"fmt"
	"strconv"

	utils "github.com/michaelchandrag/botfood-go/utils"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	sugarLogger *zap.SugaredLogger
	logger      *zap.Logger
	filename    string
	maxsize     int
	maxBackups  int
	maxAge      int
	compress    bool
	proddev     string
)

func init() {
	filename = utils.GetEnv("LOGGER_FILENAME", "foo.log")                  //true
	maxsize, _ = strconv.Atoi(utils.GetEnv("LOGGER_MAXSIZE", "100"))       //true
	maxBackups, _ = strconv.Atoi(utils.GetEnv("LOGGER_MAXBACKUPS", "100")) //true
	maxAge, _ = strconv.Atoi(utils.GetEnv("LOGGER_MAXAGE", "100"))         //true
	proddev = utils.GetEnv("LOGGER_DEBUG", "0")                            //true
	compress = false
	if utils.GetEnv("LOGGER_COMPRESS", "true") == "true" {
		compress = true
	}
	initLogger()
}
func initLogger() {
	fmt.Println("Init Logger")
	writerSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel)
	logger = zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	c := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "file",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewJSONEncoder(c)
}

// rotated. It defaults to 100 megabytes.
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxsize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func GetLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func GetLoggerEncodeFile() *zap.SugaredLogger {
	return sugarLogger
}