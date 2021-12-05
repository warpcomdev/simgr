package internal

import "log"

type Logger interface {
	Info(msg string, params ...interface{})
	Error(msg string, params ...interface{})
	Logger() *log.Logger
}

type logger struct {
	logger *log.Logger
	quiet  bool
}

func (l logger) Info(msg string, params ...interface{}) {
	if !l.quiet {
		l.logger.Printf(msg, params...)
	}
}

func (l logger) Error(msg string, params ...interface{}) {
	l.logger.Printf(msg, params...)
}

func NewLogger(quiet bool) logger {
	return logger{
		logger: log.Default(),
		quiet:  quiet,
	}
}

func (l logger) Logger() *log.Logger {
	return l.logger
}
