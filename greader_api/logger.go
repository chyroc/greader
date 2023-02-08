package greader_api

import (
	"context"
	"log"
	"os"
)

type ILogger interface {
	Info(ctx context.Context, format string, args ...interface{})
	Error(ctx context.Context, format string, args ...interface{})
}

type logger struct {
	log *log.Logger
}

func NewDefaultLogger() ILogger {
	return &logger{log.New(os.Stderr, "[greader] ", log.LstdFlags)}
}

func (r *logger) Info(ctx context.Context, format string, args ...interface{}) {
	r.log.Printf(format, args...)
}

func (r *logger) Error(ctx context.Context, format string, args ...interface{}) {
	r.log.Printf(format, args...)
}
