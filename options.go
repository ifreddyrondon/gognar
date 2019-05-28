package bastion

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

const (
	defaultInternalErrMsg = "looks like something went wrong"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
	PanicLevel = "panic"
)

const defaultProfilerRoutePrefix = "/debug"

func defaultInternalErrCallbackFn(l zerolog.Logger) func(code int, r io.Reader) {
	return func(code int, r io.Reader) {
		var buf bytes.Buffer
		buf.ReadFrom(r)
		l.Info().
			Str("component", "internal error middleware").
			Int("status", code).
			Msg(buf.String())
	}
}

func defaultRecoveryCallbackFn(l zerolog.Logger) func(req *http.Request, err error) {
	return func(req *http.Request, err error) {
		l.Error().
			Str("component", "recovery middleware").
			Err(err).Dict("req", logreq(req)).
			Msg("Recovery middleware catch an error")
	}
}

// Options are used to define how the application should run.
type Options struct {
	// InternalErrMsg message returned to the user when catch a 500 status error.
	InternalErrMsg string
	// DisableInternalErrorMiddleware boolean flag to disable the internal error middleware.
	DisableInternalErrorMiddleware bool
	// DisableRecoveryMiddleware boolean flag to disable the recovery middleware.
	DisableRecoveryMiddleware bool
	// DisablePingRouter boolean flag to disable the ping router.
	DisablePingRouter bool
	// DisableLoggerMiddleware boolean flag to disable the logger middleware.
	DisableLoggerMiddleware bool
	// DisablePrettyLogging don't output a colored human readable version on the out writer.
	DisablePrettyLogging bool
	// LoggerOutput logger output writer. Default os.Stdout
	LoggerOutput io.Writer
	// LoggerLevel defines log levels. Default "debug".
	LoggerLevel string
	// ProfilerRoutePrefix is an optional path prefix for profiler subrouter. If left unspecified, `/debug/`
	// is used as the default path prefix.
	ProfilerRoutePrefix string
	// DisableProfiler boolean flag to disable the profiler router.
	DisableProfiler bool
	// RequestIDHeaderName header name to look out the request id.
	RequestIDHeaderName string

	enableProductionMode  *bool
	internalErrorCallback func(code int, reader io.Reader)
	recoveryCallback      func(req *http.Request, err error)
}

// IsProduction check if app is running in production mode
func (opts Options) IsProduction() bool {
	return *opts.enableProductionMode
}

func mode(isProdMode *bool) *bool {
	if isProdMode != nil {
		return isProdMode
	}
	isProdMode = new(bool)
	modeEnv := defaultString(os.Getenv("GO_ENV"), "")
	if modeEnv == "" {
		modeEnv = defaultString(os.Getenv("GO_ENVIRONMENT"), "")
	}
	if modeEnv == "production" || modeEnv == "prod" {
		*isProdMode = true
	}
	return isProdMode
}

func defaultString(s1, s2 string) string {
	if s1 == "" {
		return s2
	}
	return s1
}

// Opt helper type to create functional options
type Opt func(*Bastion)

// InternalErrMsg set the message returned to the user when catch a 500 status error.
func InternalErrMsg(msg string) Opt {
	return func(app *Bastion) {
		app.InternalErrMsg = msg
	}
}

// InternalErrCallback sets the callback function when internal error middleware catch a 500 error.
func InternalErrCallback(f func(int, io.Reader)) Opt {
	return func(app *Bastion) {
		app.internalErrorCallback = f
	}
}

// DisableInternalErrorMiddleware turn off internal error middleware.
func DisableInternalErrorMiddleware() Opt {
	return func(app *Bastion) {
		app.DisableInternalErrorMiddleware = true
	}
}

// RecoveryCallback sets the callback function to handler the request when recovers from panics.
func RecoveryCallback(f func(*http.Request, error)) Opt {
	return func(app *Bastion) {
		app.recoveryCallback = f
	}
}

// DisableRecoveryMiddleware turn off recovery middleware.
func DisableRecoveryMiddleware() Opt {
	return func(app *Bastion) {
		app.DisableRecoveryMiddleware = true
	}
}

// DisablePingRouter turn off ping route.
func DisablePingRouter() Opt {
	return func(app *Bastion) {
		app.DisablePingRouter = true
	}
}

func DisableLoggerMiddleware() Opt {
	return func(app *Bastion) {
		app.DisableLoggerMiddleware = true
	}
}

// DisablePrettyLogging turn off the pretty logging.
func DisablePrettyLogging() Opt {
	return func(app *Bastion) {
		app.DisablePrettyLogging = true
	}
}

// LoggerLevel set the logger level.
func LoggerLevel(lvl string) Opt {
	return func(app *Bastion) {
		app.LoggerLevel = lvl
	}
}

// LoggerOutput set the logger output writer
func LoggerOutput(w io.Writer) Opt {
	return func(app *Bastion) {
		app.LoggerOutput = w
	}
}

// ProductionMode set the app to production mode or force debug (false).
func ProductionMode(on ...bool) Opt {
	return func(app *Bastion) {
		var enable bool
		switch len(on) {
		case 0:
			enable = true
		case 1:
			enable = on[0]
		default:
			panic("too much parameters, ProductionMode only accepts one optional param.")
		}
		app.enableProductionMode = &enable
	}
}

// ProfilerRoutePrefix set the prefix path for the profile router.
func ProfilerRoutePrefix(prefix string) Opt {
	return func(app *Bastion) {
		if !strings.HasPrefix(prefix, "/") {
			app.ProfilerRoutePrefix = "/" + prefix
		} else {
			app.ProfilerRoutePrefix = prefix
		}
	}
}

// DisableProfiler turn on the profiler router.
func DisableProfiler() Opt {
	return func(app *Bastion) {
		app.DisableProfiler = true
	}
}

// RequestIDHeaderName set the header name to look out the request id.
func RequestIDHeaderName(headerName string) Opt {
	return func(app *Bastion) {
		app.RequestIDHeaderName = headerName
	}
}
