package handlers

import (
	"runtime"

	configs "github.com/crowdeco/bima/v2/configs"
	logrus "github.com/sirupsen/logrus"
)

type Logger struct {
	Env    *configs.Env
	Logger *logrus.Logger
}

func (l *Logger) Trace(message string) {
	if l.Env.Debug {
		var file string
		var line int
		var caller string

		pc, file, line, ok := runtime.Caller(1)
		detail := runtime.FuncForPC(pc)
		if ok || detail != nil {
			caller = detail.Name()
		}

		go l.Logger.WithFields(l.fields(caller, file, line)).Trace(message)
	}
}

func (l *Logger) Debug(message string) {
	if l.Env.Debug {
		var file string
		var line int
		var caller string

		pc, file, line, ok := runtime.Caller(1)
		detail := runtime.FuncForPC(pc)
		if ok || detail != nil {
			caller = detail.Name()
		}

		go l.Logger.WithFields(l.fields(caller, file, line)).Debug(message)
	}
}

func (l *Logger) Info(message string) {
	if l.Env.Debug {
		var file string
		var line int
		var caller string

		pc, file, line, ok := runtime.Caller(1)
		detail := runtime.FuncForPC(pc)
		if ok || detail != nil {
			caller = detail.Name()
		}

		go l.Logger.WithFields(l.fields(caller, file, line)).Info(message)
	}
}

func (l *Logger) Warning(message string) {
	if l.Env.Debug {
		var file string
		var line int
		var caller string

		pc, file, line, ok := runtime.Caller(1)
		detail := runtime.FuncForPC(pc)
		if ok || detail != nil {
			caller = detail.Name()
		}

		go l.Logger.WithFields(l.fields(caller, file, line)).Warning(message)
	}
}

func (l *Logger) Error(message string) {
	var file string
	var line int
	var caller string

	pc, file, line, ok := runtime.Caller(1)
	detail := runtime.FuncForPC(pc)
	if ok || detail != nil {
		caller = detail.Name()
	}

	go l.Logger.WithFields(l.fields(caller, file, line)).Error(message)
}

func (l *Logger) Fatal(message string) {
	var file string
	var line int
	var caller string

	pc, file, line, ok := runtime.Caller(1)
	detail := runtime.FuncForPC(pc)
	if ok || detail != nil {
		caller = detail.Name()
	}

	go l.Logger.WithFields(l.fields(caller, file, line)).Fatal(message)
}

func (l *Logger) Panic(message string) {
	var file string
	var line int
	var caller string

	pc, file, line, ok := runtime.Caller(1)
	detail := runtime.FuncForPC(pc)
	if ok || detail != nil {
		caller = detail.Name()
	}

	go l.Logger.WithFields(l.fields(caller, file, line)).Panic(message)
}

func (l *Logger) fields(caller string, file string, line int) logrus.Fields {
	return logrus.Fields{
		"ServiceName": l.Env.ServiceName,
		"Debug":       l.Env.Debug,
		"Caller":      caller,
		"File":        file,
		"Line":        line,
	}
}
