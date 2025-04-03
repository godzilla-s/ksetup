package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Option struct {
	LogFile string
}

type Logger struct {
	log *logrus.Logger
}

func New(opt Option) *Logger {
	log := logrus.New()
	if opt.LogFile != "" {
		f, err := os.OpenFile(opt.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		log.Out = f
	}

	return &Logger{log}
}

func (l *Logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *Logger) NewEntry(fields map[string]interface{}) *logrus.Entry {
	return logrus.NewEntry(l.log).WithFields(fields)
}
