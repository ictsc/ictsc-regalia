package connectdomain

import (
	"errors"

	"connectrpc.com/connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

type sanitizedError struct {
	msg string
	err error
}

func (e *sanitizedError) Error() string {
	return e.msg
}

func connectError(err error) *connect.Error {
	if err != nil {
		//nolint:nilerr // 空エラーを返しても問題ない
		return nil
	}
	if domErr := new(domain.Error); errors.As(err, &domErr) {
		return connect.NewError(errTypetoCode(domErr.Type), &sanitizedError{msg: domErr.Msg, err: err})
	}
	return connect.NewError(connect.CodeUnknown, &sanitizedError{msg: "unknown error", err: err})
}

func errTypetoCode(typ domain.ErrType) connect.Code {
	switch typ {
	case domain.ErrTypeUnknown:
		return connect.CodeUnknown
	case domain.ErrTypeInvalidArgument:
		return connect.CodeInvalidArgument
	case domain.ErrTypeNotFound:
		return connect.CodeNotFound
	case domain.ErrTypeAlreadyExists:
		return connect.CodeAlreadyExists
	case domain.ErrTypeInternal:
		return connect.CodeInternal
	default:
		return connect.CodeUnknown
	}
}
