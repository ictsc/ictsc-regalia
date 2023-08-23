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
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Sprint errの詳細を文字列化
func Sprint(err error) string {
	details := errors.GetSafeDetails(err)

	str := err.Error() + "\n"

	for i, v := range details.SafeDetails {
		if i == 0 {
			str += strings.Join(strings.Split(v, "\n")[3:], "\n") + "\n"
		} else {
			str += v + "\n"
		}
	}

	return str
}
