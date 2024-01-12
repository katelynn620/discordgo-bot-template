package util

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogConfig struct {
	Level      string `json:"level"`       // Level is the minimum log level, DEBUG<INFO<WARN<ERROR<FATAL. For example: info --> collect logs of info level and above.
	FileName   string `json:"file_name"`   // FileName is the location of the log file.
	MaxSize    int    `json:"max_size"`    // MaxSize is the maximum size of the log file in megabytes before it is rotated. Default is 100MB.
	MaxAge     int    `json:"max_age"`     // MaxAge is the maximum number of days to retain old log files based on the timestamp encoded in the file name.
	MaxBackups int    `json:"max_backups"` // MaxBackups is the maximum number of old log files to retain. By default, all old log files are retained (although MaxAge may still cause them to be deleted).
}

var (
	logger *zap.Logger
)

// Responsible for setting the encoding format of the log
func getFileEncoder() zapcore.Encoder {
	// Get a specified EncoderConfig for customization
	encodeConfig := zap.NewProductionEncoderConfig()

	// Set the keys used for each log entry. If any key is empty, that part of the entry is omitted.

	// Serialize time. eg: 2022-09-01T19:11:35.921+0800
	encodeConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	// "time":"2022-09-01T19:11:35.921+0800"
	encodeConfig.TimeKey = "time"
	// Serialize the Level as an uppercase string. For example, serialize info level as INFO.
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encodeConfig)
}

func getConsoleEncoder() zapcore.Encoder {
	// Get a specified EncoderConfig for customization
	encodeConfig := zap.NewDevelopmentEncoderConfig()
	encodeConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encodeConfig)
}

func getFileWriter(filename string, maxsize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,  // path of filename
		MaxSize:    maxsize,   // max size of log file(MB)
		MaxAge:     maxAge,    // max age of old files
		MaxBackups: maxBackup, // max size of old files
		Compress:   false,     // compress old files or not
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getConsoleWriter() zapcore.WriteSyncer {
	return zapcore.AddSync(os.Stdout)
}

func InitLogger(cfg LogConfig) (err error) {
	var zapcores []zapcore.Core

	level := new(zapcore.Level)
	err = level.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		panic(fmt.Sprintf("unmarshal log level error: %v", err))
	}

	fileWriteSyncer := getFileWriter(cfg.FileName, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	fileEncoder := getFileEncoder()
	fileCore := zapcore.NewCore(fileEncoder, fileWriteSyncer, level)

	consoleWriteSyncer := getConsoleWriter()
	consoleEncoder := getConsoleEncoder()
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriteSyncer, level)

	zapcores = append(zapcores, fileCore, consoleCore)
	logger = zap.New(zapcore.NewTee(zapcores...), zap.AddCaller())

	// replace zap's global logger instance, and then you can call zap.L() in other packages
	zap.ReplaceGlobals(logger)
	return
}
