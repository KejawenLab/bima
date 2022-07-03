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
		verbose bool
		service string
		data    logrus.Fields
		Engine  *logrus.Logger
	}
)

func Default(service string) {
	Configure(true, service, LoggerExtension{})
}

func Configure(debug bool, service string, extensions LoggerExtension) {
	Engine := logrus.New()
	if debug {
		Engine.SetLevel(logrus.DebugLevel)
	}

	Engine.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	for _, e := range extensions.Extensions {
		Engine.AddHook(e)
	}

	Logger = &logger{
		verbose: debug,
		service: service,
		Engine:  Engine,
		data:    logrus.Fields{},
	}
}

func (l *LoggerExtension) Register(extensions []logrus.Hook) {
	l.Extensions = extensions
}

func (l *logger) Add(key string, value interface{}) {
	l.data[key] = value
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

	go l.Engine.WithFields(l.data).Trace(message)
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

	go l.Engine.WithFields(l.data).Debug(message)
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

	go l.Engine.WithFields(l.data).Info(message)
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

	go l.Engine.WithFields(l.data).Warning(message)
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

	go l.Engine.WithFields(l.data).Error(message)
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

	go l.Engine.WithFields(l.data).Fatal(message)
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

	go l.Engine.WithFields(l.data).Panic(message)
}

func (l *logger) fields(caller string, file string, line int) {
	workDir, _ := os.Getwd()
	l.data["debug"] = l.verbose
	l.data["service"] = l.service
	l.data["trace"] = map[string]interface{}{
		"caller": caller,
		"file":   strings.Replace(file, workDir, ".", 1),
		"line":   line,
	}
}
