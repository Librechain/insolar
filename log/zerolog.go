//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package log

import (
	"fmt"
	"github.com/insolar/insolar/critlog"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"go.opencensus.io/stats"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
)

var insolarPrefix = "github.com/insolar/insolar/"

func trimInsolarPrefix(file string, line int) string {
	var skip = 0
	if idx := strings.Index(file, insolarPrefix); idx != -1 {
		skip = idx + len(insolarPrefix)
	}
	return file[skip:] + ":" + strconv.Itoa(line)
}

func init() {
	zerolog.TimeFieldFormat = timestampFormat
	zerolog.CallerMarshalFunc = trimInsolarPrefix
}

type callerHookConfig struct {
	enabled        bool
	skipFrameCount int
	funcname       bool
}

var _ insolar.Logger = &zerologAdapter{}
var _ insolar.SpecialLoggerFactory = &zerologAdapter{}

type zerologAdapter struct {
	logger       zerolog.Logger
	output       io.WriteCloser
	outputWraps  outputWrapFlag
	level        zerolog.Level
	callerConfig callerHookConfig
}

type outputWrapFlag uint32

const (
	outputWrappedWithBuffer outputWrapFlag = 1 << iota
	outputWrappedWithCritical
)

type loglevelChangeHandler struct {
}

func NewLoglevelChangeHandler() http.Handler {
	handler := &loglevelChangeHandler{}
	return handler
}

// ServeHTTP is an HTTP handler that changes the global minimum log level
func (h *loglevelChangeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	levelStr := "(nil)"
	if values["level"] != nil {
		levelStr = values["level"][0]
	}
	level, err := insolar.ParseLevel(levelStr)
	if err != nil {
		w.WriteHeader(500)
		_, _ = fmt.Fprintf(w, "Invalid level '%v': %v\n", levelStr, err)
		return
	}

	zlevel, err := InternalLevelToZerologLevel(level)
	if err != nil {
		w.WriteHeader(500)
		_, _ = fmt.Fprintf(w, "Invalid level '%v': %v\n", levelStr, err)
		return
	}

	zerolog.SetGlobalLevel(zlevel)

	w.WriteHeader(200)
	_, _ = fmt.Fprintf(w, "New log level: '%v'\n", levelStr)
}

func InternalLevelToZerologLevel(level insolar.LogLevel) (zerolog.Level, error) {
	switch level {
	case insolar.DebugLevel:
		return zerolog.DebugLevel, nil
	case insolar.InfoLevel:
		return zerolog.InfoLevel, nil
	case insolar.WarnLevel:
		return zerolog.WarnLevel, nil
	case insolar.ErrorLevel:
		return zerolog.ErrorLevel, nil
	case insolar.FatalLevel:
		return zerolog.FatalLevel, nil
	case insolar.PanicLevel:
		return zerolog.PanicLevel, nil
	}
	return zerolog.NoLevel, errors.New("Unknown internal level")
}

var _ io.WriteCloser = &closableConsoleWriter{}

type closableConsoleWriter struct {
	zerolog.ConsoleWriter
}

func (p *closableConsoleWriter) Close() error {
	if c, ok := p.Out.(io.Closer); ok {
		return c.Close()
	}
	return errors.New("unsupported: Close")
}

func (p *closableConsoleWriter) Sync() error {
	if c, ok := p.Out.(*os.File); ok {
		return c.Sync()
	}
	return errors.New("unsupported: Sync")
}

func newIsolatedOutput() io.WriteCloser {
	return os.NewFile(uintptr(syscall.Stderr), "/dev/stderr")
}

func newDefaultTextOutput() io.WriteCloser {
	return &closableConsoleWriter{zerolog.ConsoleWriter{
		Out:          newIsolatedOutput(),
		NoColor:      true,
		TimeFormat:   timestampFormat,
		PartsOrder:   fieldsOrder,
		FormatCaller: formatCaller(),
	}}
}

func selectFormatter(format insolar.LogFormat) (io.WriteCloser, error) {
	var output io.WriteCloser

	switch format {
	case insolar.TextFormat:
		output = newDefaultTextOutput()
	case insolar.JSONFormat:
		output = newIsolatedOutput()
	default:
		return nil, errors.New("unknown formatter " + format.String())
	}

	return output, nil
}

func newZerologAdapter(cfg configuration.Log) (*zerologAdapter, error) {
	format, err := insolar.ParseFormat(cfg.Formatter)
	if err != nil {
		return nil, err
	}

	za := &zerologAdapter{
		level: zerolog.InfoLevel,
		callerConfig: callerHookConfig{
			enabled:        true,
			skipFrameCount: defaultCallerSkipFrameCount,
		},
	}

	za.output, err = selectFormatter(format)
	if err != nil {
		return nil, err
	}

	if cfg.BufferSize > 0 {
		za.output = createBufferedLogger(za.output, cfg.BufferSize)
		za.outputWraps |= outputWrappedWithBuffer
	}
	za.output = critlog.FatalFlusher(za.output)

	logger := zerolog.New(za.output).Level(zerolog.InfoLevel).With().Timestamp().Logger()
	logger = logger.Hook(&metricsHook{})
	za.logger = logger

	return za, nil
}

// WithFields return copy of adapter with predefined fields.
func (z *zerologAdapter) WithFields(fields map[string]interface{}) insolar.Logger {
	zCtx := z.logger.With()
	for key, value := range fields {
		zCtx = zCtx.Interface(key, value)
	}

	zCopy := *z
	zCopy.logger = zCtx.Logger()
	return &zCopy
}

// WithField return copy of adapter with predefined single field.
func (z *zerologAdapter) WithField(key string, value interface{}) insolar.Logger {
	zCopy := *z
	zCopy.logger = z.logger.With().Interface(key, value).Logger()
	return &zCopy
}

// Debug logs a message at level Debug on the stdout.
func (z *zerologAdapter) Debug(args ...interface{}) {
	stats.Record(contextWithLogLevel(zerolog.DebugLevel), statLogCalls.M(1))
	z.loggerWithHooks().Debug().Msg(fmt.Sprint(args...))
}

// Debugf formatted logs a message at level Debug on the stdout.
func (z *zerologAdapter) Debugf(format string, args ...interface{}) {
	stats.Record(contextWithLogLevel(zerolog.DebugLevel), statLogCalls.M(1))
	z.loggerWithHooks().Debug().Msgf(format, args...)
}

// Info logs a message at level Info on the stdout.
func (z *zerologAdapter) Info(args ...interface{}) {
	stats.Record(contextWithLogLevel(zerolog.InfoLevel), statLogCalls.M(1))
	z.loggerWithHooks().Info().Msg(fmt.Sprint(args...))
}

// Infof formatted logs a message at level Info on the stdout.
func (z *zerologAdapter) Infof(format string, args ...interface{}) {
	stats.Record(contextWithLogLevel(zerolog.InfoLevel), statLogCalls.M(1))
	z.loggerWithHooks().Info().Msgf(format, args...)
}

// Warn logs a message at level Warn on the stdout.
func (z *zerologAdapter) Warn(args ...interface{}) {
	stats.Record(contextWithLogLevel(zerolog.WarnLevel), statLogCalls.M(1))
	z.loggerWithHooks().Warn().Msg(fmt.Sprint(args...))
}

// Warnf formatted logs a message at level Warn on the stdout.
func (z *zerologAdapter) Warnf(format string, args ...interface{}) {
	stats.Record(contextWithLogLevel(zerolog.WarnLevel), statLogCalls.M(1))
	z.loggerWithHooks().Warn().Msgf(format, args...)
}

// Error logs a message at level Error on the stdout.
func (z *zerologAdapter) Error(args ...interface{}) {
	stats.Record(contextWithLogLevel(zerolog.ErrorLevel), statLogCalls.M(1))
	z.loggerWithHooks().Error().Msg(fmt.Sprint(args...))
}

// Errorf formatted logs a message at level Error on the stdout.
func (z *zerologAdapter) Errorf(format string, args ...interface{}) {
	stats.Record(contextWithLogLevel(zerolog.ErrorLevel), statLogCalls.M(1))
	z.loggerWithHooks().Error().Msgf(format, args...)
}

// Fatal logs a message at level Fatal on the stdout.
func (z *zerologAdapter) Fatal(args ...interface{}) {
	stats.Record(contextWithLogLevel(zerolog.FatalLevel), statLogCalls.M(1))
	z.loggerWithHooks().Fatal().Msg(fmt.Sprint(args...))
}

// Fatalf formatted logs a message at level Fatal on the stdout.
func (z *zerologAdapter) Fatalf(format string, args ...interface{}) {
	stats.Record(contextWithLogLevel(zerolog.FatalLevel), statLogCalls.M(1))
	z.loggerWithHooks().Fatal().Msgf(format, args...)
}

// Panic logs a message at level Panic on the stdout.
func (z *zerologAdapter) Panic(args ...interface{}) {
	stats.Record(contextWithLogLevel(zerolog.PanicLevel), statLogCalls.M(1))
	z.loggerWithHooks().Panic().Msg(fmt.Sprint(args...))
}

// Panicf formatted logs a message at level Panic on the stdout.
func (z *zerologAdapter) Panicf(format string, args ...interface{}) {
	stats.Record(contextWithLogLevel(zerolog.PanicLevel), statLogCalls.M(1))
	z.loggerWithHooks().Panic().Msgf(format, args...)
}

// WithLevel sets log level.
func (z *zerologAdapter) WithLevel(level string) (insolar.Logger, error) {
	levelNumber, err := insolar.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	return z.WithLevelNumber(levelNumber)
}

// WithLevelNumber sets log level with constant.
func (z *zerologAdapter) WithLevelNumber(level insolar.LogLevel) (insolar.Logger, error) {
	if level == insolar.NoLevel {
		return z, nil
	}
	zerologLevel, err := InternalLevelToZerologLevel(level)
	if err != nil {
		return nil, err
	}
	zCopy := *z
	zCopy.level = zerologLevel
	zCopy.logger = z.logger.Level(zerologLevel)
	return &zCopy, nil
}

// SetOutput sets the output destination for the logger.
func (z *zerologAdapter) WithOutput(w io.Writer) insolar.Logger {
	zCopy := *z
	zCopy.logger = z.logger.Output(w)
	return &zCopy
}

// WithCaller switch on/off 'caller' field computation.
func (z *zerologAdapter) WithCaller(flag bool) insolar.Logger {
	zCopy := *z
	zCopy.callerConfig.enabled = flag
	// if caller disabled, probably we should avoid cost of call runtime.Caller, so disable func field
	if !flag {
		zCopy.callerConfig.funcname = flag
	}
	return &zCopy
}

// WithSkipFrameCount changes skipFrameCount by delta value (it can be negative).
func (z *zerologAdapter) WithSkipFrameCount(delta int) insolar.Logger {
	zCopy := *z
	zCopy.callerConfig.skipFrameCount += delta
	return &zCopy
}

// WithCaller switch on/off 'func' field computation.
func (z *zerologAdapter) WithFuncName(flag bool) insolar.Logger {
	zCopy := *z
	zCopy.callerConfig.funcname = flag
	return &zCopy
}

// WithFormat sets logger output format
func (z *zerologAdapter) WithFormat(format insolar.LogFormat) (insolar.Logger, error) {
	output, err := selectFormatter(format)
	if err != nil {
		return nil, err
	}

	return z.WithOutput(output), nil
}

func (z *zerologAdapter) loggerWithHooks() *zerolog.Logger {
	l := z.logger
	if z.callerConfig.funcname {
		l = l.Hook(newCallerHook(z.callerConfig.skipFrameCount + 2))
	} else if z.callerConfig.enabled {
		l = l.With().CallerWithSkipFrameCount(z.callerConfig.skipFrameCount).Logger()
	}
	return &l
}

func (z *zerologAdapter) Is(level insolar.LogLevel) bool {
	zerologLevel, err := InternalLevelToZerologLevel(level)
	if err != nil {
		panic(err)
	}

	return zerologLevel >= z.level && zerologLevel >= zerolog.GlobalLevel()
}

func (z *zerologAdapter) CreateCriticalLogger(bufSize int) insolar.Logger {
	if z.outputWraps&outputWrappedWithCritical != 0 {
		return z
	}
	zCopy := *z
	zCopy.output = critlog.NewCriticalWriter(z.output)

	// we disregard existence of a previous buffer for critical events
	zCopy.outputWraps = (z.outputWraps &^ outputWrappedWithBuffer) | outputWrappedWithCritical

	if bufSize > 0 {
		zCopy.output = createBufferedLogger(zCopy.output, bufSize)
		zCopy.outputWraps |= outputWrappedWithBuffer
	}

	zCopy.logger = z.logger.Output(zCopy.output)
	return &zCopy
}

func createBufferedLogger(output io.WriteCloser, bufSize int) io.WriteCloser {
	if bufSize <= 0 {
		panic("illegal value")
	}
	missedFunc := func(missed int) { panic(fmt.Errorf("logger dropped %d messages", missed)) }
	dw := diode.NewWriter(output, bufSize, 0, missedFunc)
	return &dw
}

func (z *zerologAdapter) CreateBufferedLogger(bufSize int) insolar.Logger {
	if bufSize <= 0 || z.outputWraps&outputWrappedWithBuffer != 0 {
		return z
	}
	zCopy := *z
	zCopy.output = createBufferedLogger(z.output, bufSize)
	zCopy.outputWraps |= outputWrappedWithBuffer
	zCopy.logger = z.logger.Output(zCopy.output)
	return &zCopy
}
