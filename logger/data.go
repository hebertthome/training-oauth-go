package logger

import (
	"errors"
	"fmt"
)

const (
	//STDOUT any message to stdout
	STDOUT Out = "stdout"
	//STDERR redirects any message to stderr
	STDERR Out = "stderr"
	//DISCARD set logger to ignore all message
	DISCARD Out = "discard"

	//LOGRUS is the github.com/Sirupsen/logrus id
	LOGRUS Provider = "logrus"
	//ZAP is the github.com/uber-go/zap id
	ZAP Provider = "zap"

	//TEXT is the text log format
	TEXT Format = "text"
	//JSON is the json log format
	JSON Format = "json"
	//LOGRUSFmtfText is the text with the logrus formatf approach
	LOGRUSFmtfText Format = "logrusFrmtfText"

	//PANIC is the panic level logger
	PANIC Level = iota
	//FATAL is the fatal level logger
	FATAL
	//ERROR is the error level logger
	ERROR
	//WARN is the warn level logger
	WARN
	//INFO is the info level logger
	INFO
	//DEBUG is the debug level logger
	DEBUG

	//StringField is a constant for string logger fields
	StringField FieldType = iota
	//IntField is a constant for string logger fields
	IntField
	//Int64Field is a constant for string logger fields
	Int64Field
	//FloatField is a constant for string logger fields
	FloatField
	//Float64Field is a constant for string logger fields
	Float64Field
	//DurationField is a constant for duration logger fields
	DurationField
	//TimeField is a constant for time logger fields
	TimeField
	//BoolField is a constant for string logger fields
	BoolField
	//StructField is a constant for string logger fields
	StructField
)

var (
	//ErrInvalidProvider is the err raised when an invalid provider was select
	ErrInvalidProvider = errors.New("Logger.InvalidProvider[Message='Avaible providers are: LOGRUS and ZAP']")
	//DefaultConfig holds the instance of the behavior parameters
	DefaultConfig *Configuration
)

//Provider is the back end implementor id of the logging feature
type Provider string

func (p Provider) String() string {
	return string(p)
}

//Out is the type for logger writer config
type Out string

func (o Out) String() string {
	return string(o)
}

// Option is used to set options for the logger.
type Option interface {
	apply(Logger)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(Logger) error

func (f optionFunc) apply(l Logger) {
	f(l)
}

//Level is the threshold of the logger
type Level int

// String returns a lower-case ASCII representation of the log level.
func (n Level) String() string {
	switch n {
	case DEBUG:
		return "debug"
	case INFO:
		return "info"
	case WARN:
		return "warn"
	case ERROR:
		return "error"
	case PANIC:
		return "panic"
	case FATAL:
		return "fatal"
	default:
		return fmt.Sprintf("Level(%d)", n)
	}
}

//Format is a parameter to controle the logger style
type Format string

//Configuration holds the log beahvior parameters
type Configuration struct {
	Provider Provider `json:"provider" mapstructure:"provider"`
	Level    Level    `json:"level" mapstructure:"level"`
	Format   Format   `json:"format" mapstructure:"format"`
	Out      Out      `json:"out" mapstructure:"out"`
}

func (l *Configuration) String() string {
	return fmt.Sprintf("Configuration[Provider=%v Level=%v Format=%v Out=%v]", l.Provider, l.Level, l.Format, l.Out)
}

//FieldType is a type identifier for logger fields
type FieldType int8

//Field is a struct to send paramaters to log messages
type Field struct {
	key     string
	val     interface{}
	valType FieldType
}

//Logger is an interface to write log messages
type Logger interface {
	Level() Level
	IsEnabled(Level) bool
	Debug(string, ...Field)
	Info(string, ...Field)
	Warn(string, ...Field)
	Error(string, ...Field)
	Panic(string, ...Field)
	Fatal(string, ...Field)
}

type baseLogger struct {
	level Level
}

func (b baseLogger) Level() Level {
	return b.level
}

func (b baseLogger) IsEnabled(level Level) bool {
	return b.level >= level
}
