package console

import "errors"

var (
	ErrTimedOut        = errors.New("timed out")
	ErrCancelledByUser = errors.New("cancelled by user")
)
