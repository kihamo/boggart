package di

import (
	"sync"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

const (
	LoggerDefaultBufferedRecordsLimit uint64 = 100
	LoggerDefaultBufferedRecordsLevel        = LogLevel(zap.DebugLevel)
)

type LogLevel zapcore.Level

func (l LogLevel) MarshalYAML() (interface{}, error) {
	return int64(l), nil
}

type LoggerContainerSupport interface {
	SetLogger(*LoggerContainer)
	Logger() *LoggerContainer
	LastRecords() []observer.LoggedEntry
	Clean()
}

func LoggerContainerBind(bind boggart.Bind) (*LoggerContainer, bool) {
	if support, ok := bind.(LoggerContainerSupport); ok {
		container := support.Logger()
		return container, container != nil
	}

	return nil, false
}

type LoggerBind struct {
	mutex     sync.RWMutex
	container *LoggerContainer
}

func (b *LoggerBind) SetLogger(container *LoggerContainer) {
	b.mutex.Lock()
	b.container = container
	b.mutex.Unlock()
}

func (b *LoggerBind) Logger() *LoggerContainer {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.container
}

func (b *LoggerBind) LastRecords() []observer.LoggedEntry {
	c := b.Logger()
	if c != nil {
		return c.getObserver().Records()
	}

	return nil
}

func (b *LoggerBind) Clean() {
	if c := b.Logger(); c != nil {
		c.getObserver().Clean()
	}
}

type loggerObserver struct {
	mutex   sync.RWMutex
	level   zapcore.Level
	limit   int
	logs    []observer.LoggedEntry
	context []zapcore.Field
}

func (o *loggerObserver) Enabled(level zapcore.Level) bool {
	if o.limit == 0 {
		return false
	}

	return o.level.Enabled(level)
}

func (o *loggerObserver) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if o.Enabled(ent.Level) {
		return ce.AddCore(ent, o)
	}

	return ce
}

func (o *loggerObserver) With(fields []zapcore.Field) zapcore.Core {
	if o.limit == 0 {
		return o
	}

	return &loggerObserver{
		level:   o.level,
		limit:   o.limit,
		logs:    o.logs,
		context: append(o.context[:len(o.context):len(o.context)], fields...),
	}
}

func (o *loggerObserver) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	if o.limit == 0 {
		return nil
	}

	all := make([]zapcore.Field, 0, len(fields)+len(o.context))
	all = append(all, o.context...)
	all = append(all, fields...)

	record := observer.LoggedEntry{
		Entry:   entry,
		Context: all,
	}

	o.mutex.Lock()

	l := len(o.logs)
	limit := o.limit - 1

	if l >= limit {
		o.logs = append(o.logs[l-limit:], record)
	} else {
		o.logs = append(o.logs, record)
	}

	o.mutex.Unlock()

	return nil
}

func (o *loggerObserver) Sync() error {
	return nil
}

func (o *loggerObserver) Records() []observer.LoggedEntry {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	tmp := make([]observer.LoggedEntry, len(o.logs))
	copy(tmp, o.logs)

	return tmp
}

func (o *loggerObserver) Clean() {
	o.mutex.Lock()
	o.logs = o.logs[:0]
	o.mutex.Unlock()
}

type LoggerContainer struct {
	logging.Logger
	bindItem boggart.BindItem

	sugar    *zap.SugaredLogger
	observer *loggerObserver
	once     sync.Once
}

func NewLoggerContainer(bindItem boggart.BindItem, logger logging.Logger) *LoggerContainer {
	return &LoggerContainer{
		Logger:   logging.NewLazyLogger(logger, logger.Name()+"."+bindItem.ID()),
		bindItem: bindItem,
	}
}

func (c *LoggerContainer) Debug(message string, args ...interface{}) {
	if c == nil {
		return
	}

	c.Logger.Debug(message, args...)
	c.getSugar().Debugw(message, args...)
}

func (c *LoggerContainer) Debugf(template string, args ...interface{}) {
	if c == nil {
		return
	}

	c.Logger.Debugf(template, args...)
	c.getSugar().Debugf(template, args...)
}

func (c *LoggerContainer) Info(message string, args ...interface{}) {
	if c == nil {
		return
	}

	c.Logger.Info(message, args...)
	c.getSugar().Infow(message, args...)
}

func (c *LoggerContainer) Infof(template string, args ...interface{}) {
	if c == nil {
		return
	}

	c.Logger.Infof(template, args...)
	c.getSugar().Infof(template, args...)
}

func (c *LoggerContainer) Warn(message string, args ...interface{}) {
	if c == nil {
		return
	}

	c.Logger.Warn(message, args...)
	c.getSugar().Warnw(message, args...)
}

func (c *LoggerContainer) Warnf(template string, args ...interface{}) {
	if c == nil {
		return
	}

	c.Logger.Warnf(template, args...)
	c.getSugar().Warnf(template, args...)
}

func (c *LoggerContainer) Error(message string, args ...interface{}) {
	if c == nil {
		return
	}

	c.Logger.Error(message, args...)
	c.getSugar().Errorw(message, args...)
}

func (c *LoggerContainer) Errorf(template string, args ...interface{}) {
	if c == nil {
		return
	}

	c.Logger.Errorf(template, args...)
	c.getSugar().Errorf(template, args...)
}

func (c *LoggerContainer) Panic(message string, args ...interface{}) {
	if c == nil {
		return
	}

	c.Logger.Panic(message, args...)
	c.getSugar().Panicw(message, args...)
}

func (c *LoggerContainer) Panicf(template string, args ...interface{}) {
	if c == nil {
		return
	}

	c.Logger.Panicf(template, args...)
	c.getSugar().Panicf(template, args...)
}

func (c *LoggerContainer) Fatal(message string, args ...interface{}) {
	if c == nil {
		return
	}

	c.Logger.Fatal(message, args...)
	c.getSugar().Fatalw(message, args...)
}

func (c *LoggerContainer) Fatalf(template string, args ...interface{}) {
	if c == nil {
		return
	}

	c.Logger.Fatalf(template, args...)
	c.getSugar().Fatalf(template, args...)
}

func (c *LoggerContainer) getSugar() *zap.SugaredLogger {
	return c.init().sugar
}

func (c *LoggerContainer) getObserver() *loggerObserver {
	return c.init().observer
}

func (c *LoggerContainer) init() *LoggerContainer {
	if c == nil {
		return nil
	}

	c.once.Do(func() {
		limit := LoggerDefaultBufferedRecordsLimit
		level := LoggerDefaultBufferedRecordsLevel

		if config, ok := ConfigForBind(c.bindItem.Bind()); ok {
			if loggerConfig, ok := config.(LoggerBufferedConfig); ok {
				limit = loggerConfig.LoggerBufferedRecordsLimit()
				level = loggerConfig.LoggerBufferedRecordsLevel()
			}
		}

		c.observer = &loggerObserver{
			level: zapcore.Level(level),
			limit: int(limit),
		}
		c.sugar = zap.New(c.observer).Sugar()
	})

	return c
}

type LoggerBufferedConfig interface {
	LoggerBufferedRecordsLimit() uint64
	LoggerBufferedRecordsLevel() LogLevel
}

type LoggerConfig struct {
	BufferedRecordsLimit uint64   `mapstructure:"logger_buffered_records_limit" yaml:"logger_buffered_records_limit"`
	BufferedRecordsLevel LogLevel `mapstructure:"logger_buffered_records_level" yaml:"logger_buffered_records_level"`
}

func LoggerConfigDefaults() (c LoggerConfig) {
	c.BufferedRecordsLimit = LoggerDefaultBufferedRecordsLimit
	c.BufferedRecordsLevel = LoggerDefaultBufferedRecordsLevel

	return c
}

func (c LoggerConfig) LoggerBufferedRecordsLimit() uint64 {
	return c.BufferedRecordsLimit
}

func (c LoggerConfig) LoggerBufferedRecordsLevel() LogLevel {
	if min := LogLevel(zap.DebugLevel); c.BufferedRecordsLevel < min {
		return min
	}

	if max := LogLevel(zap.FatalLevel); c.BufferedRecordsLevel > max {
		return max
	}

	return c.BufferedRecordsLevel
}
