package constraint_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/egobatch/internal/constraint"
)

// customError is a simple error type used in tests
// customError 是测试用的简单错误类型
type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

// TestErrorType tests ErrorType constraint
// TestErrorType 测试 ErrorType 约束
func TestErrorType(t *testing.T) {
	erx := &customError{msg: "test error"}
	require.False(t, constraint.Pass(erx))
}

func TestPass(t *testing.T) {
	var erx *customError
	require.True(t, constraint.Pass(erx))
}
