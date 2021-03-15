package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	rootLogrusConfig *logrusConfig
)

func (f Format) toLogrusFormatter() logrus.Formatter {
	switch f {
	case JSON:
		return new(logrus.JSONFormatter)
	default:
		return new(logrus.TextFormatter)
	}
}

func (o Out) toLogrusOut() (io.Writer, error) {
	return getOutput(o)
}

func (n Level) toLogrusLevel() logrus.Level {
	return logrus.Level(n)
}

type logrusLogger struct {
	baseLogger
	logger *logrus.Logger
}

func (l logrusLogger) toLogrusFields(fields ...Field) logrus.Fields {
	logrusFields := make(map[string]interface{})
	for _, v := range fields {
		logrusFields[v.key] = v.val
	}
	return logrusFields
}

func (l logrusLogger) toInterfaceSlice(fields ...Field) []interface{} {
	logrusFields := make([]interface{}, len(fields))
	for i, v := range fields {
		logrusFields[i] = v.val
	}
	return logrusFields
}

func (l logrusLogger) Debug(message string, fields ...Field) {
	if len(fields) <= 0 {
		l.logger.Debug(message)
		return
	}
	if DefaultConfig.Format == LOGRUSFmtfText {
		l.logger.Debugf(message, l.toInterfaceSlice(fields...)...)
		return
	}
	l.logger.WithFields(l.toLogrusFields(fields...)).Debug(message)
}

func (l logrusLogger) Info(message string, fields ...Field) {
	if len(fields) <= 0 {
		l.logger.Info("[Logrus]" + message)
	}
	if DefaultConfig.Format == LOGRUSFmtfText {
		l.logger.Infof("[Logrus]"+message, l.toInterfaceSlice(fields...)...)
		return
	}
	l.logger.WithFields(l.toLogrusFields(fields...)).Info("[Logrus]" + message)
}

func (l logrusLogger) Warn(message string, fields ...Field) {
	if len(fields) <= 0 {
		l.logger.Warn(message)
	}
	if DefaultConfig.Format == LOGRUSFmtfText {
		l.logger.Warnf(message, l.toInterfaceSlice(fields...)...)
		return
	}
	l.logger.WithFields(l.toLogrusFields(fields...)).Warn(message)
}

func (l logrusLogger) Error(message string, fields ...Field) {
	if len(fields) <= 0 {
		l.logger.Error(message)
	}
	if DefaultConfig.Format == LOGRUSFmtfText {
		l.logger.Errorf(message, l.toInterfaceSlice(fields...)...)
		return
	}
	l.logger.WithFields(l.toLogrusFields(fields...)).Error(message)
}

func (l logrusLogger) Panic(message string, fields ...Field) {
	if len(fields) <= 0 {
		l.logger.Panic(message)
	}
	if DefaultConfig.Format == LOGRUSFmtfText {
		l.logger.Panicf(message, l.toInterfaceSlice(fields...)...)
		return
	}
	l.logger.WithFields(l.toLogrusFields(fields...)).Panic(message)
}

func (l logrusLogger) Fatal(message string, fields ...Field) {
	if len(fields) <= 0 {
		l.logger.Fatal(message)
	}
	if DefaultConfig.Format == LOGRUSFmtfText {
		l.logger.Fatalf(message, l.toInterfaceSlice(fields...)...)
		return
	}
	l.logger.WithFields(l.toLogrusFields(fields...)).Fatal(message)
}

func setupLogrus(loggerConfig Configuration) error {
	logrusConf, errs := createLogrusConfig(loggerConfig)
	if errs != nil && len(errs) > 0 {
		return fmt.Errorf("SetupLogrusErr[Errs=%v]", errs)
	}
	rootLogrusConfig = &logrusConf
	logrus.SetLevel(rootLogrusConfig.level)
	logrus.SetFormatter(rootLogrusConfig.formatter)
	logrus.SetOutput(rootLogrusConfig.output)
	loggerFactory = newLogrus
	return nil
}

func newLogrus(config Configuration) Logger {
	logrusConf, _ := createLogrusConfig(config)
	logger := new(logrusLogger)
	logger.logger = &logrus.Logger{
		Level:     logrusConf.level,
		Formatter: logrusConf.formatter,
		//Hooks: make(LevelHooks),
		Out: logrusConf.output,
	}
	return logger
}

type logrusConfig struct {
	output    io.Writer
	formatter logrus.Formatter
	level     logrus.Level
}

func createLogrusConfig(cfg Configuration) (logrusConfig, []error) {
	var errs []error
	var output io.Writer
	if cfg.Out == Out("") {
		output = os.Stdout
	} else if tmpWriter, tmpErr := cfg.Out.toLogrusOut(); tmpErr != nil {
		errs = append(errs, tmpErr)
		output = os.Stdout
	} else {
		output = tmpWriter
	}
	var level logrus.Level
	if cfg.Level == Level(0) {
		level = logrus.DebugLevel
	} else {
		level = logrus.Level(cfg.Level)
	}
	return logrusConfig{
		level:     level,
		formatter: cfg.Format.toLogrusFormatter(),
		output:    output,
	}, errs
}
