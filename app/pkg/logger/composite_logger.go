package logger

import (
	"IMP/app/internal/ports"
	"fmt"
	"runtime"
	"strings"
)

var instance CompositeLogger

type CompositeLogger struct {
	loggers []ports.Logger
}

func Init(loggers []ports.Logger) {
	instance = CompositeLogger{
		loggers: loggers,
	}
}

func Info(msg string, ctx map[string]interface{}) {
	for _, logger := range instance.loggers {
		logger.Info("[INFO] "+msg, ctx)
	}
}

func Warn(msg string, ctx map[string]interface{}) {
	for _, logger := range instance.loggers {
		logger.Warn("[WARNING] "+msg, ctx)
	}
}

func Error(msg string, ctx map[string]interface{}) {
	ctx = buildErrorContextWithStackTrace(ctx)

	for _, logger := range instance.loggers {
		logger.Error("[ERROR] "+msg, ctx)
	}
}

func buildErrorContextWithStackTrace(ctx map[string]interface{}) map[string]interface{} {
	context := cloneContext(ctx)

	if _, exists := context["stackTrace"]; exists {
		return context
	}

	stackTrace := extractStackTraceFromError(context["error"])
	if stackTrace == "" {
		stackTrace = buildFallbackStackTrace()
	}

	context["stackTrace"] = stackTrace

	return context
}

func cloneContext(ctx map[string]interface{}) map[string]interface{} {
	if ctx == nil {
		return map[string]interface{}{}
	}

	cloned := make(map[string]interface{}, len(ctx))
	for key, value := range ctx {
		cloned[key] = value
	}

	return cloned
}

func extractStackTraceFromError(errValue interface{}) string {
	err, ok := errValue.(error)
	if !ok || err == nil {
		return ""
	}

	expanded := fmt.Sprintf("%+v", err)
	if expanded != "" && expanded != err.Error() {
		return expanded
	}

	return ""
}

func buildFallbackStackTrace() string {
	const maxFrames = 48
	const skipFrames = 3

	pcs := make([]uintptr, maxFrames)
	count := runtime.Callers(skipFrames, pcs)
	if count == 0 {
		return ""
	}

	frames := runtime.CallersFrames(pcs[:count])
	lines := make([]string, 0, count)

	for {
		frame, more := frames.Next()
		if !more {
			if shouldIncludeFrame(frame.Function) {
				lines = append(lines, formatFrame(frame))
			}
			break
		}

		if shouldIncludeFrame(frame.Function) {
			lines = append(lines, formatFrame(frame))
		}
	}

	return strings.Join(lines, "\n")
}

func shouldIncludeFrame(function string) bool {
	return !strings.Contains(function, "/app/pkg/logger.")
}

func formatFrame(frame runtime.Frame) string {
	return fmt.Sprintf("%s\n\t%s:%d", frame.Function, frame.File, frame.Line)
}
