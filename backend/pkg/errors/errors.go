// Package errors エラーユーティリティー
package errors

import (
	"strings"

	"github.com/cockroachdb/errors"
)

// New 新規エラー生成
//
//	fmt.Printf()同様、フォーマット記法に対応
func New(message string, args ...any) error {
	return errors.Newf(message, args...)
}

// Wrap エラーをラップする
func Wrap(err error) error {
	return errors.Wrap(err, "")
}

// Is errors.Is()のラッパー
func Is(err, reference error) bool {
	return errors.Is(err, reference)
}

// As errors.As()のラッパー
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Sprint errの詳細を文字列化
func Sprint(err error) string {
	details := errors.GetSafeDetails(err)

	str := err.Error() + "\n  -- stack trace:\n  | "
	// 最初の2行は errors.New() もしくは errors.Wrap() の情報なので除外
	str += strings.Join(strings.Split(details.SafeDetails[0], "\n")[3:], "\n  | ")

	return str
}
