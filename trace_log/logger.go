package trace_log

import (
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger zerolog.Logger

func Init(fileName string, enableLog string) {
	if len(fileName) == 0 {
		fileName = "socket_chat.log"
	}
	logFile := &lumberjack.Logger{
		Filename:   fileName, // Log file name
		MaxSize:    10,       // Maximum size in megabytes before log rotation
		MaxBackups: 10,       // Maximum number of old log files to retain
		MaxAge:     28,       // Maximum number of days to retain log files
		Compress:   true,     // Whether to compress old log files
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = MarshalStack
	output := zerolog.ConsoleWriter{Out: logFile, TimeFormat: time.Kitchen}
	Logger = zerolog.New(output).With().Timestamp().Logger().Level(zerolog.ErrorLevel)
	if enableLog == "1" {
		Logger = Logger.Level(zerolog.DebugLevel)
	}
}
