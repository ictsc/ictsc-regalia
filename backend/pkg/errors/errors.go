package errors

import (
	"strings"

	"github.com/cockroachdb/errors"
)

// New 新規エラー生成
//
//	fmt.Printf()同様、フォーマット記法に対応
func New(flag Flag, message string, args ...any) error {
	return newAppError(flag, errors.Newf(message, args...))
}

// Wrap エラーをラップする
func Wrap(flag Flag, err error) error {
	return newAppError(flag, errors.Wrap(err, ""))
}

// Is errors.Is()のラッパー
func Is(err, reference error) bool {
	return errors.Is(err, reference)
}

// IsFlag エラーにフラグがあるか判定する
func IsFlag(err error, flag Flag) bool {
	for {
		if err == nil {
			return false
		}

		// nolint:errorlint
		if e, ok := err.(*appError); ok {
			if e.flag == flag {
				return true
			}
		}

		err = errors.Unwrap(err)
	}
}

// As errors.As()のラッパー
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Sprint errの詳細を文字列化
func Sprint(err error) string {
	details := errors.GetSafeDetails(errors.Unwrap(err))

	str := err.Error() + "\n\t-- stack trace:\n\t| "
	// 最初の2行は errors.New() もしくは errors.Wrap() の情報なので除外
	str += strings.Join(strings.Split(details.SafeDetails[0], "\n")[3:], "\n\t| ")

	return str
}
