package myassert

import (
	"github.com/stretchr/testify/require"
	"github.com/yyle88/egobatch/internal/myerrors"
)

// NoError requires that error is nil
// Solves Go's (*T)(nil) != nil interface trap via adaptation
//
// NoError 要求错误为 nil
// 通过适配解决 Go 的 (*T)(nil) != nil 接口陷阱
func NoError(t require.TestingT, err *myerrors.Error, msgAndArgs ...interface{}) {
	require.NoError(t, adapt(err), msgAndArgs...)
}

// Error requires that error is not nil
// Provides robust error presence validation via internal adaptation
//
// Error 要求错误不为 nil
// 通过内部适配提供健壮的错误存在验证
func Error(t require.TestingT, err *myerrors.Error, msgAndArgs ...interface{}) {
	require.Error(t, adapt(err), msgAndArgs...)
}

// adapt converts custom error to standard error interface
// Solves Go's core nil interface problem where (*T)(nil) != nil
// Returns true nil when given nil pointer to prevent interface pollution
//
// adapt 将自定义错误转换为标准错误接口
// 解决 Go 的核心 nil 接口问题：(*T)(nil) != nil
// 给定 nil 指针时返回真正的 nil，防止接口污染
func adapt(err *myerrors.Error) error {
	if err != nil {
		return err
	}
	return nil
}
