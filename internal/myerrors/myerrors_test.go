package myerrors_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/egobatch/internal/myerrors"
)

func TestError(t *testing.T) {
	err := myerrors.New("TEST_CODE", "test message: %s", "abc")
	require.NotNil(t, err)
	require.Equal(t, "[TEST_CODE] test message: abc", err.Error())
	require.Equal(t, "TEST_CODE", err.Code())
}

func TestErrorServiceError(t *testing.T) {
	err := myerrors.ErrorServiceError("service error: %d", 123)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "SERVICE_ERROR")
	require.Contains(t, err.Error(), "service error: 123")
	require.True(t, myerrors.IsServiceError(err))
}

func TestErrorWrongContext(t *testing.T) {
	err := myerrors.ErrorWrongContext("context error: %v", "timeout")
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "CONTEXT_ERROR")
	require.Contains(t, err.Error(), "context error: timeout")
	require.True(t, myerrors.IsWrongContext(err))
}

func TestIsServiceError(t *testing.T) {
	svcErr := myerrors.ErrorServiceError("service error")
	ctxErr := myerrors.ErrorWrongContext("context error")

	require.True(t, myerrors.IsServiceError(svcErr))
	require.False(t, myerrors.IsServiceError(ctxErr))
	require.False(t, myerrors.IsServiceError(nil))
}

func TestIsWrongContext(t *testing.T) {
	ctxErr := myerrors.ErrorWrongContext("context error")
	svcErr := myerrors.ErrorServiceError("service error")

	require.True(t, myerrors.IsWrongContext(ctxErr))
	require.False(t, myerrors.IsWrongContext(svcErr))
	require.False(t, myerrors.IsWrongContext(nil))
}
