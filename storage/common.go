package storage

import (
	"context"
	"errors"
)

var (
	ErrDbNotFound error = errors.New("not found in db")
)

func IsErrNotFound(err error) bool {
	return errors.Is(err, ErrDbNotFound)
}

func IsErrDeadline(err error) bool {
	return errors.Is(err, context.DeadlineExceeded)
}
