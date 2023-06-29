package log

import (
	"os"

	"github.com/Bloxico/exchange-gateway/sofija/projectpath"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	LogDirectory        string = "/logs"
	LogFileRelativePath string = "/logs/out.log"
)

type ZapLogger struct {
	zlog *zap.SugaredLogger
}

func NewLogger(pretty bool) (*ZapLogger, error) {
	zlog, err := newZapLogger(pretty)
	if err != nil {
		return nil, errors.Wrap(err, "new zap logger")
	}

	zap.ReplaceGlobals(zlog)
	return &ZapLogger{
		zlog: zlog.Sugar(),
	}, nil
}

func newZapLogger(pretty bool) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	if pretty {
		cfg = zap.NewDevelopmentConfig()
	}

	// Disable sampling of log messages until needed.
	cfg.Sampling = nil

	//Define log path
	if _, err := os.Stat(projectpath.Root + LogDirectory); os.IsNotExist(err) {
		os.Mkdir(projectpath.Root+LogDirectory, 0700)
	}

	logPath := projectpath.Root + LogFileRelativePath
	os.OpenFile(logPath, os.O_CREATE, 0666)

	cfg.OutputPaths = []string{
		"stderr",
		logPath,
	}

	return cfg.Build(zap.AddCallerSkip(1))
}

func (log *ZapLogger) Debug(msg string, fields ...interface{}) {
	log.zlog.Debugw(msg, fields...)
}

func (log *ZapLogger) Info(msg string, fields ...interface{}) {
	log.zlog.Infow(msg, fields...)
}

func (log *ZapLogger) Infof(msg string, fields ...interface{}) {
	log.zlog.Infof(msg, fields...)
}

func (log *ZapLogger) Warn(msg string, fields ...interface{}) {
	log.zlog.Warnw(msg, fields...)
}

func (log *ZapLogger) Error(msg string, fields ...interface{}) {
	log.zlog.Errorw(msg, fields...)
}

func (log *ZapLogger) Errorf(msg string, fields ...interface{}) {
	log.zlog.Errorf(msg, fields...)
}
