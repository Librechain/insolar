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
	"github.com/insolar/insolar/log/critlog"
	"github.com/insolar/insolar/log/inssyslog"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.opencensus.io/stats"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

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

type zerologAdapter struct {
	logger      zerolog.Logger
	output      io.Writer
	outputWraps outputWrapFlag
	//	bareOutput   io.WriteCloser
	level        zerolog.Level
	callerConfig callerHookConfig

	format     insolar.LogFormat
	bufferSize int
}

type outputWrapFlag uint32

const (
	outputWrappedWithBuffer outputWrapFlag = 1 << iota
	outputWrappedWithCritical
	outputWrappedWithFormatter
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

func newDefaultTextOutput(out io.Writer) io.WriteCloser {
	return &closableConsoleWriter{zerolog.ConsoleWriter{
		Out:          out,
		NoColor:      true,
		TimeFormat:   timestampFormat,
		PartsOrder:   fieldsOrder,
		FormatCaller: formatCaller(),
	}}
}

func selectOutput(output insolar.LogOutput) (io.WriteCloser, error) {
	switch output {
	case insolar.StdErrOutput:
		// we open a separate file handle as it will be closed, so it should not interfere with os.Stderr
		return os.NewFile(uintptr(syscall.Stderr), "/dev/stderr"), nil
	case insolar.SysLogOutput:
		return inssyslog.ConnectDefaultSyslog("insolar") // breaks dependency on windows
	default:
		return nil, errors.New("unknown output " + output.String())
	}
}

func selectFormatter(format insolar.LogFormat, output io.Writer) (io.Writer, error) {
	switch format {
	case insolar.TextFormat:
		return newDefaultTextOutput(output), nil
	case insolar.JSONFormat:
		return output, nil
	default:
		return nil, errors.New("unknown formatter " + format.String())
	}
}

const (
	defaultLogFormat = insolar.TextFormat
	defaultLogOutput = insolar.StdErrOutput
)

func newZerologAdapter(cfg configuration.Log) (*zerologAdapter, error) {
	outputType, err := insolar.ParseOutput(cfg.OutputType, defaultLogOutput)
	if err != nil {
		return nil, err
	}

	format, err := insolar.ParseFormat(cfg.Formatter, defaultLogFormat)
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

	za.output, err = selectOutput(outputType)
	if err != nil {
		return nil, err
	}

	za.format = format
	za.bufferSize = cfg.BufferSize

	err = za.prepareOutput()
	if err != nil {
		return nil, err
	}

	logger := zerolog.New(za.output).Level(zerolog.InfoLevel).With().Timestamp().Logger()
	logger = logger.Hook(&metricsHook{})
	za.logger = logger

	return za, nil
}

func (z *zerologAdapter) prepareOutput() error {
	var err error
	bareOutput := z.output
	z.output, err = selectFormatter(z.format, z.output)
	if err != nil {
		return err
	}

	if z.output != bareOutput {
		z.outputWraps |= outputWrappedWithFormatter
	}

	if z.bufferSize > 0 {
		dropBufOnFatal := z.bufferSize > 1000

		z.output = critlog.NewDiodeBufferedLevelWriter(z.output, z.bufferSize,
			10*time.Millisecond,
			dropBufOnFatal,
			func(missed int) []byte {
				return ([]byte)(fmt.Sprintf("logger dropped %d messages", missed))
			},
		)
		z.outputWraps |= outputWrappedWithBuffer | outputWrappedWithCritical
	} else {
		z.output = critlog.NewFatalDirectWriter(z.output)
	}
	return nil
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
	zCopy.output = w
	zCopy.outputWraps = 0

	err := zCopy.prepareOutput()
	if err != nil {
		panic(err)
	}

	zCopy.logger = z.logger.Output(zCopy.output)
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
// Deprecated: format change has no proper impact actually
func (z *zerologAdapter) WithFormat(format insolar.LogFormat) (insolar.Logger, error) {
	output, err := selectFormatter(format, z.output)
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
