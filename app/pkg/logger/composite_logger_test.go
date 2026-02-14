package logger

import (
	"IMP/app/internal/ports"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stackAwareError struct {
	message string
	stack   string
}

func (s stackAwareError) Error() string {
	return s.message
}

func (s stackAwareError) Format(state fmt.State, verb rune) {
	if verb == 'v' {
		_, _ = state.Write([]byte(s.stack))
		return
	}

	_, _ = state.Write([]byte(s.message))
}

type logCall struct {
	message string
	context map[string]interface{}
}

type fakeLogger struct {
	infoCalls  []logCall
	warnCalls  []logCall
	errorCalls []logCall
}

func (f *fakeLogger) IsEnabled() bool {
	return true
}

func (f *fakeLogger) Info(message string, context map[string]interface{}) {
	f.infoCalls = append(f.infoCalls, logCall{message: message, context: context})
}

func (f *fakeLogger) Warn(message string, context map[string]interface{}) {
	f.warnCalls = append(f.warnCalls, logCall{message: message, context: context})
}

func (f *fakeLogger) Error(message string, context map[string]interface{}) {
	f.errorCalls = append(f.errorCalls, logCall{message: message, context: context})
}

func TestInfo_FanOutAndPrefix(t *testing.T) {
	l1 := &fakeLogger{}
	l2 := &fakeLogger{}
	Init([]ports.Logger{l1, l2})

	ctx := map[string]interface{}{"requestId": "abc-123"}
	Info("process started", ctx)

	require.Len(t, l1.infoCalls, 1)
	require.Len(t, l2.infoCalls, 1)
	assert.Equal(t, "[INFO] process started", l1.infoCalls[0].message)
	assert.Equal(t, "[INFO] process started", l2.infoCalls[0].message)
	assert.Equal(t, "abc-123", l1.infoCalls[0].context["requestId"])
}

func TestWarn_FanOutAndPrefix(t *testing.T) {
	l1 := &fakeLogger{}
	l2 := &fakeLogger{}
	Init([]ports.Logger{l1, l2})

	ctx := map[string]interface{}{"task": "poll"}
	Warn("slow response", ctx)

	require.Len(t, l1.warnCalls, 1)
	require.Len(t, l2.warnCalls, 1)
	assert.Equal(t, "[WARNING] slow response", l1.warnCalls[0].message)
	assert.Equal(t, "[WARNING] slow response", l2.warnCalls[0].message)
	assert.Equal(t, "poll", l2.warnCalls[0].context["task"])
}

func TestError_AddsFallbackStackTraceAndDoesNotMutateInputContext(t *testing.T) {
	l1 := &fakeLogger{}
	l2 := &fakeLogger{}
	Init([]ports.Logger{l1, l2})

	inputCtx := map[string]interface{}{"error": errors.New("plain error"), "taskType": "poll"}
	Error("failed", inputCtx)

	require.Len(t, l1.errorCalls, 1)
	require.Len(t, l2.errorCalls, 1)
	assert.Equal(t, "[ERROR] failed", l1.errorCalls[0].message)
	assert.Equal(t, "poll", l1.errorCalls[0].context["taskType"])
	assert.Contains(t, l1.errorCalls[0].context, "stackTrace")
	assert.NotEmpty(t, l1.errorCalls[0].context["stackTrace"])

	_, hasStackInOriginal := inputCtx["stackTrace"]
	assert.False(t, hasStackInOriginal, "original context must not be mutated")
}

func TestError_UsesExistingStackTraceFromContext(t *testing.T) {
	l := &fakeLogger{}
	Init([]ports.Logger{l})

	Error("failed", map[string]interface{}{
		"error":      errors.New("plain"),
		"stackTrace": "precomputed-stack",
	})

	require.Len(t, l.errorCalls, 1)
	assert.Equal(t, "precomputed-stack", l.errorCalls[0].context["stackTrace"])
}

func TestError_UsesEmbeddedStackTraceFromError(t *testing.T) {
	l := &fakeLogger{}
	Init([]ports.Logger{l})

	Error("failed", map[string]interface{}{
		"error": stackAwareError{
			message: "wrapped error",
			stack:   "embedded-stack-line-1\nembedded-stack-line-2",
		},
	})

	require.Len(t, l.errorCalls, 1)
	assert.Equal(t, "embedded-stack-line-1\nembedded-stack-line-2", l.errorCalls[0].context["stackTrace"])
}

func TestError_WorksWithNilContext(t *testing.T) {
	l := &fakeLogger{}
	Init([]ports.Logger{l})

	Error("failed", nil)

	require.Len(t, l.errorCalls, 1)
	assert.Equal(t, "[ERROR] failed", l.errorCalls[0].message)
	assert.Contains(t, l.errorCalls[0].context, "stackTrace")
}

func TestBuildErrorContextWithStackTrace_UsesFallbackWhenNoEmbeddedStack(t *testing.T) {
	ctx := map[string]interface{}{
		"error": errors.New("plain error"),
	}

	result := buildErrorContextWithStackTrace(ctx)

	assert.NotEmpty(t, result["stackTrace"])
	assert.Contains(t, result["stackTrace"], ".go:")
}

func TestBuildErrorContextWithStackTrace_UsesEmbeddedStackWhenPresent(t *testing.T) {
	ctx := map[string]interface{}{
		"error": stackAwareError{
			message: "wrapped error",
			stack:   "embedded-stack-line-1\nembedded-stack-line-2",
		},
	}

	result := buildErrorContextWithStackTrace(ctx)

	assert.Equal(t, "embedded-stack-line-1\nembedded-stack-line-2", result["stackTrace"])
}

func TestBuildErrorContextWithStackTrace_DoesNotOverrideExistingStack(t *testing.T) {
	ctx := map[string]interface{}{
		"error":      errors.New("plain"),
		"stackTrace": "precomputed-stack",
	}

	result := buildErrorContextWithStackTrace(ctx)

	assert.Equal(t, "precomputed-stack", result["stackTrace"])
}
