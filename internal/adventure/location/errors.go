package location

import "errors"

var (
	ErrWrongLabel              = errors.New("wrong location label")
	ErrUndefDesc               = errors.New("undefined location description")
	ErrConnLocationNotFound    = errors.New("connection location not found")
	ErrConnWordNotFound        = errors.New("connection word not found")
	ErrConnWordUndefinedID     = errors.New("connection word undefined ID")
	ErrConnLocationUndefinedID = errors.New("connection location undefined ID")
)
