package main

import (
	"fmt"
	"os"
)

type Logger interface {
	Infof(format string, a ...interface{})
	Debugf(format string, a ...interface{})
}

func NewStdOutLogger(debug bool) Logger {
	return &StdOutLogger{
		debug: debug,
	}
}

type StdOutLogger struct {
	debug bool
}

func (l *StdOutLogger) Infof(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}

func (l *StdOutLogger) Debugf(format string, a ...interface{}) {
	if l.debug {
		fmt.Fprintf(os.Stdout, format, a...)
	}
}
