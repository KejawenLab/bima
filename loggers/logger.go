package loggers

import (
	"context"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var Logger *logger

type (
	LoggerExtension struct {
		Extensions []logrus.Hook
	}

	logger struct {
		Verbose bool
		Service string
		Engine  *logrus.Logger
		Data    logrus.Fields
	}
)

func Default(service string) {
	Configure(true, service, LoggerExtension{})
}

func Configure(debug bool, service string, extensions LoggerExtension) {
	engine := logrus.New()
	if debug {
		engine.SetLevel(logrus.DebugLevel)
	}

	engine.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	for _, e := range extensions.Extensions {
		engine.AddHook(e)
	}

	Logger = &logger{
		Verbose: debug,
		Service: service,
		Engine:  engine,
		Data:    logrus.Fields{},
	}
}

func (l *LoggerExtension) Register(extensions []logrus.Hook) {
	l.Extensions = extensions
}

func (l *logger) Add(key string, value interface{}) {
	l.Data[key] = value
}

func (l *logger) Trace(ctx context.Context, message string) {
	var file, caller string
	var line int

	pc, file, line, ok := runtime.Caller(1)
	detail := runtime.FuncForPC(pc)
	if ok && detail != nil {
		caller = detail.Name()
	}

	l.Add("scope", ctx.Value("scope"))
	l.fields(caller, file, line)

	go l.Engine.WithFields(l.Data).Trace(message)
}

func (l *logger) Debug(ctx context.Context, message string) {
	var file, caller string
	var line int

	pc, file, line, ok := runtime.Caller(1)
	detail := runtime.FuncForPC(pc)
	if ok && detail != nil {
		caller = detail.Name()
	}

	l.Add("scope", ctx.Value("scope"))
	l.fields(caller, file, line)

	go l.Engine.WithFields(l.Data).Debug(message)
}

func (l *logger) Info(ctx context.Context, message string) {
	var file, caller string
	var line int

	pc, file, line, ok := runtime.Caller(1)
	detail := runtime.FuncForPC(pc)
	if ok && detail != nil {
		caller = detail.Name()
	}

	l.Add("scope", ctx.Value("scope"))
	l.fields(caller, file, line)

	go l.Engine.WithFields(l.Data).Info(message)
}

func (l *logger) Warning(ctx context.Context, message string) {
	var file, caller string
	var line int

	pc, file, line, ok := runtime.Caller(1)
	detail := runtime.FuncForPC(pc)
	if ok && detail != nil {
		caller = detail.Name()
	}

	l.Add("scope", ctx.Value("scope"))
	l.fields(caller, file, line)

	go l.Engine.WithFields(l.Data).Warning(message)
}

func (l *logger) Error(ctx context.Context, message string) {
	var file, caller string
	var line int

	pc, file, line, ok := runtime.Caller(1)
	detail := runtime.FuncForPC(pc)
	if ok && detail != nil {
		caller = detail.Name()
	}

	l.Add("scope", ctx.Value("scope"))
	l.fields(caller, file, line)

	go l.Engine.WithFields(l.Data).Error(message)
}

func (l *logger) Fatal(ctx context.Context, message string) {
	var file, caller string
	var line int

	pc, file, line, ok := runtime.Caller(1)
	detail := runtime.FuncForPC(pc)
	if ok && detail != nil {
		caller = detail.Name()
	}

	l.Add("scope", ctx.Value("scope"))
	l.fields(caller, file, line)

	go l.Engine.WithFields(l.Data).Fatal(message)
}

func (l *logger) Panic(ctx context.Context, message string) {
	var file, caller string
	var line int

	pc, file, line, ok := runtime.Caller(1)
	detail := runtime.FuncForPC(pc)
	if ok && detail != nil {
		caller = detail.Name()
	}

	l.Add("scope", ctx.Value("scope"))
	l.fields(caller, file, line)

	go l.Engine.WithFields(l.Data).Panic(message)
}

func (l *logger) fields(caller string, file string, line int) {
	workDir, _ := os.Getwd()
	l.Data["debug"] = l.Verbose
	l.Data["service"] = l.Service
	l.Data["trace"] = map[string]interface{}{
		"caller": caller,
		"file":   strings.Replace(file, workDir, ".", 1),
		"line":   line,
	}
}
