package logger

import (
	"fmt"
	"os"
	"platform/internal/tracer"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Logger struct {
	logger *zerolog.Logger
}

func NewLogger() *Logger {
	l := zerolog.New(os.Stdout).With().Caller().Stack().Logger()

	// pretty output, harder parse. For development only
	l = l.Output(zerolog.ConsoleWriter{
		Out: os.Stdout,
		FormatTimestamp: func(i interface{}) string {
			return ""
		},
		NoColor: false,
	})

	return &Logger{logger: &l}
}

type LogExtra struct {
	key   string
	value any
}

func Extra(key string, value any) LogExtra {
	return LogExtra{key: key, value: value}
}

func (l *Logger) Info(c *fiber.Ctx, events ...any) {
	l.log(c, zerolog.InfoLevel, events...)
}

func (l *Logger) Error(c *fiber.Ctx, events ...any) {
	l.log(c, zerolog.ErrorLevel, events...)
}

func (l *Logger) log(c *fiber.Ctx, level zerolog.Level, events ...any) {
	var event *zerolog.Event
	switch level {
	case zerolog.ErrorLevel:
		event = l.logger.Error()
	case zerolog.WarnLevel:
		event = l.logger.Warn()
	case zerolog.InfoLevel:
		event = l.logger.Info()
	case zerolog.DebugLevel:
		event = l.logger.Debug()
	case zerolog.TraceLevel:
		event = l.logger.Trace()
	default:
		event = l.logger.Info()
	}

	b := strings.Builder{}
	for _, extra := range events {
		// check if LogExtra
		if extra, ok := extra.(LogExtra); ok {
			event = event.Str(extra.key, fmt.Sprintf("%+v", extra.value))
			continue
		}

		if b.Len() > 0 {
			b.WriteString("\t")
		}
		b.WriteString(fmt.Sprintf("%+v", extra))
	}

	reqID := tracer.GetRequestID(c)
	if reqID != "" {
		event = event.Str("requestID", reqID)
	}

	event.Msg(b.String())
}

func Middleware(l *Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := ctx.Next()

		// Log details after request processing
		l.Info(ctx,
			Extra("method", ctx.Method()),
			Extra("path", ctx.Path()),
			Extra("status", ctx.Response().StatusCode()),
			Extra("latency", time.Since(start).String()),
			Extra("ip", ctx.IP()),
			Extra("status", ctx.Response().StatusCode()),
		)

		return err
	}
}
