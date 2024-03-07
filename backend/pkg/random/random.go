package random

import (
	"crypto/rand"

	"github.com/ictsc/ictsc-outlands/backend/pkg/errors"
)

// NewString 指定された長さのランダム文字列を生成
func NewString(digit uint32) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$"

	if digit == 0 {
		return "", errors.New(errors.ErrBadArgument, "digit must be greater than 0")
	}

	byteSlice := make([]byte, digit)
	if _, err := rand.Read(byteSlice); err != nil {
		return "", errors.Wrap(errors.ErrUnknown, err)
	}

	var result string
	for _, v := range byteSlice {
		result += string(letters[int(v)%len(letters)])
	}

	return result, nil
}
