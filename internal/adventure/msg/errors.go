package msg

import "errors"

var (
	ErrMsgAlreadyExists = errors.New("message already exists")
	ErrMsgNotFound      = errors.New("message not found")
	ErrMsgEmpty         = errors.New("message is empty")
	ErrMsgPluralEmpty   = errors.New("message plural is not complete, it should have 'zero', 'one' and 'many' texts")
)
