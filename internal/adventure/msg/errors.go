package msg

import "errors"

var ErrMsgAlreadyExists = errors.New("message already exists")
var ErrMsgNotFound = errors.New("message not found")
