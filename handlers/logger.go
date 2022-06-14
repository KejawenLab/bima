package handlers

import (
	"context"
	"os"
	"runtime"
	"strings"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	Verbose bool
	Service configs.Service
	Logger  *logrus.Logger
	Data    logrus.Fields
}

func (l *Logger) Add(key string, value interface{}) {
	l.Data[key] = value
}

func (l *Logger) Trace(ctx context.Context, message string) {
	if l.Verbose {
		var file, caller string
		var line int

		pc, file, line, ok := runtime.Caller(1)
		detail := runtime.FuncForPC(pc)
		if ok && detail != nil {
			caller = detail.Name()
		}

		l.Add("scope", ctx.Value("scope"))
		l.fields(caller, file, line)

		go l.Logger.WithFields(l.Data).Trace(message)
	}
}

func (l *Logger) Debug(ctx context.Context, message string) {
	if l.Verbose {
		var file, caller string
		var line int

		pc, file, line, ok := runtime.Caller(1)
		detail := runtime.FuncForPC(pc)
		if ok && detail != nil {
			caller = detail.Name()
		}

		l.Add("scope", ctx.Value("scope"))
		l.fields(caller, file, line)

		go l.Logger.WithFields(l.Data).Info(message)
	}
}

func (l *Logger) Info(ctx context.Context, message string) {
	if l.Verbose {
		var file, caller string
		var line int

		pc, file, line, ok := runtime.Caller(1)
		detail := runtime.FuncForPC(pc)
		if ok && detail != nil {
			caller = detail.Name()
		}

		l.Add("scope", ctx.Value("scope"))
		l.fields(caller, file, line)

		go l.Logger.WithFields(l.Data).Info(message)
	}
}

func (l *Logger) Warning(ctx context.Context, message string) {
	if l.Verbose {
		var file, caller string
		var line int

		pc, file, line, ok := runtime.Caller(1)
		detail := runtime.FuncForPC(pc)
		if ok && detail != nil {
			caller = detail.Name()
		}

		l.Add("scope", ctx.Value("scope"))
		l.fields(caller, file, line)

		go l.Logger.WithFields(l.Data).Warning(message)
	}
}

func (l *Logger) Error(ctx context.Context, message string) {
	var file, caller string
	var line int

	pc, file, line, ok := runtime.Caller(1)
	detail := runtime.FuncForPC(pc)
	if ok && detail != nil {
		caller = detail.Name()
	}

	l.Add("scope", ctx.Value("scope"))
	l.fields(caller, file, line)

	go l.Logger.WithFields(l.Data).Error(message)
}

func (l *Logger) Fatal(ctx context.Context, message string) {
	var file, caller string
	var line int

	pc, file, line, ok := runtime.Caller(1)
	detail := runtime.FuncForPC(pc)
	if ok && detail != nil {
		caller = detail.Name()
	}

	l.Add("scope", ctx.Value("scope"))
	l.fields(caller, file, line)

	go l.Logger.WithFields(l.Data).Fatal(message)
}

func (l *Logger) Panic(ctx context.Context, message string) {
	var file, caller string
	var line int

	pc, file, line, ok := runtime.Caller(1)
	detail := runtime.FuncForPC(pc)
	if ok && detail != nil {
		caller = detail.Name()
	}

	l.Add("scope", ctx.Value("scope"))
	l.fields(caller, file, line)

	go l.Logger.WithFields(l.Data).Panic(message)
}

func (l *Logger) fields(caller string, file string, line int) {
	workDir, _ := os.Getwd()
	l.Data["debug"] = l.Verbose
	l.Data["service"] = map[string]string{
		"name": l.Service.ConnonicalName,
		"host": l.Service.Host,
	}
	l.Data["trace"] = map[string]interface{}{
		"caller": caller,
		"file":   strings.Replace(file, workDir, ".", 1),
		"line":   line,
	}
}
