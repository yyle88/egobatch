package myassert_test

import (
	"testing"

	"github.com/yyle88/egobatch/internal/myassert"
	"github.com/yyle88/egobatch/internal/myerrors"
)

// TestNoError tests NoError function handles typed nil as expected
// TestNoError 测试 NoError 函数正确处理 typed nil
func TestNoError(t *testing.T) {
	var err *myerrors.Error = nil
	myassert.NoError(t, err)
}

// TestError tests Error function with non-nil error
// TestError 测试 Error 函数处理非 nil 错误
func TestError(t *testing.T) {
	err := myerrors.ErrorServiceError("test error")
	myassert.Error(t, err)
}
