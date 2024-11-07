package logger

import (
	"errors"

	"github.com/koma2211/you-meal/internal/config"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	InfoLog  zerolog.Logger
	ErrorLog zerolog.Logger
	WarnLog  zerolog.Logger
	DebugLog zerolog.Logger
}

func InitLogger(logs config.Logs) (*Logger, error) {
	if !(logs.Env == config.EnvLocal || logs.Env == config.EnvDev || logs.Env == config.EnvProd) {
		return nil, errors.New("env-config logger is not valid")
	}

	zeroInfo := zerolog.New(&lumberjack.Logger{
		Filename:   logs.InfoPath,   // File name
		MaxSize:    logs.MaxSize,    // Size in MB before file gets rotated
		MaxBackups: logs.MaxBackups, // Max number of files kept before being overwritten
		MaxAge:     logs.MaxAge,     // Max number of days to keep the files
		Compress:   true,            // Whether to compress log files using gzip
		LocalTime:  true,
	}).With().Str("env", logs.Env).Caller().Timestamp().Logger().Level(zerolog.InfoLevel)

	zeroError := zerolog.New(&lumberjack.Logger{
		Filename:   logs.ErrorPath,  // File name
		MaxSize:    logs.MaxSize,    // Size in MB before file gets rotated
		MaxBackups: logs.MaxBackups, // Max number of files kept before being overwritten
		MaxAge:     logs.MaxAge,     // Max number of days to keep the files
		Compress:   true,            // Whether to compress log files using gzip
		LocalTime:  true,
	}).With().Str("env", logs.Env).Caller().Timestamp().Logger().Level(zerolog.ErrorLevel)

	zeroWarn := zerolog.New(&lumberjack.Logger{
		Filename:   logs.WarnPath,   // File name
		MaxSize:    logs.MaxSize,    // Size in MB before file gets rotated
		MaxBackups: logs.MaxBackups, // Max number of files kept before being overwritten
		MaxAge:     logs.MaxAge,     // Max number of days to keep the files
		Compress:   true,            // Whether to compress log files using gzip
		LocalTime:  true,
	}).With().Str("env", logs.Env).Caller().Timestamp().Logger().Level(zerolog.WarnLevel)

	zeroDebug := zerolog.New(&lumberjack.Logger{
		Filename:   logs.DebugPath,  // File name
		MaxSize:    logs.MaxSize,    // Size in MB before file gets rotated
		MaxBackups: logs.MaxBackups, // Max number of files kept before being overwritten
		MaxAge:     logs.MaxAge,     // Max number of days to keep the files
		Compress:   true,            // Whether to compress log files using gzip
		LocalTime:  true,
	}).With().Str("env", logs.Env).Caller().Timestamp().Logger().Level(zerolog.DebugLevel)

	return &Logger{
		InfoLog:  zeroInfo,
		ErrorLog: zeroError,
		DebugLog: zeroDebug,
		WarnLog:  zeroWarn,
	}, nil
}

func InitAccessLog(logs config.Logs) (zerolog.Logger, error) {
	if !(logs.Env == config.EnvLocal || logs.Env == config.EnvDev || logs.Env == config.EnvProd) {
		return zerolog.Logger{}, errors.New("env-config logger is not valid")
	}

	accessLog := zerolog.New(&lumberjack.Logger{
		Filename:   logs.AccessPath, // File name
		MaxSize:    logs.MaxSize,    // Size in MB before file gets rotated
		MaxBackups: logs.MaxBackups, // Max number of files kept before being overwritten
		MaxAge:     logs.MaxAge,     // Max number of days to keep the files
		Compress:   true,            // Whether to compress log files using gzip
		LocalTime:  true,
	}).With().Str("env", logs.Env).Caller().Timestamp().Logger()

	return accessLog, nil
}
