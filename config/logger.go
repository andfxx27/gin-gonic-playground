package config

import (
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

func NewMultiLevelWriter() io.Writer {
	var consoleLogWriter io.Writer = zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	fileLogger := &lumberjack.Logger{
		Filename: "logs/app.log",
		MaxSize:  50,
		Compress: true,
	}

	var aggregateWriter io.Writer = zerolog.MultiLevelWriter(consoleLogWriter, fileLogger)
	return aggregateWriter
}
