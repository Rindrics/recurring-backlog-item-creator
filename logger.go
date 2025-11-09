package main

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	debug *log.Logger
	info  *log.Logger
}

var logger *Logger

func init() {
	logger = &Logger{
		debug: log.New(io.Discard, "", log.LstdFlags),
		info:  log.New(os.Stderr, "", log.LstdFlags),
	}
}

func SetDebugMode(enabled bool) {
	if enabled {
		logger.debug = log.New(os.Stderr, "[DEBUG] ", log.LstdFlags)
	} else {
		logger.debug = log.New(io.Discard, "", log.LstdFlags)
	}
}

func Debug(v ...interface{}) {
	logger.debug.Println(v...)
}

func Debugf(format string, v ...interface{}) {
	logger.debug.Printf(format, v...)
}
