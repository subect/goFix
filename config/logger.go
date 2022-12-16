package config

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"runtime"
)

var sugarLogger *zap.SugaredLogger

func (config *Configuration) Logger() *zap.SugaredLogger {
	InitLogger()
	defer func(sugarLogger *zap.SugaredLogger) {
		err := sugarLogger.Sync()
		if err != nil {

		}
	}(sugarLogger)
	return Logger()
}

func Logger() *zap.SugaredLogger {
	if sugarLogger == nil {
		panic("init logger error")
	}
	return sugarLogger
}

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogPath() string {
	if runtime.GOOS == "windows" {
		//return os.Getenv("GOPATH") + "\\src\\goFix\\logger\\goFix.log"
		return "./logger/goFix.log"
	} else {
		return os.Getenv("GOPATH") + "/goFix/logger/goFix.log"
	}
	//return os.Getenv("GOPATH") + "/goFix/logger/goFix.log"
}

func getLogWriter() zapcore.WriteSyncer {

	lumberJackLogger := &lumberjack.Logger{
		Filename:   getLogPath(), //日志文件的位置
		MaxSize:    1,            //在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 5,            //保留旧文件的最大个数
		MaxAge:     30,           //保留旧文件的最大天数
		Compress:   false,        //是否压缩/归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger)
}
