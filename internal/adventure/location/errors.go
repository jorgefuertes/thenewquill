package location

import "errors"

var (
	ErrWrongLabel           = errors.New("wrong location label")
	ErrUndefLabel           = errors.New("undefined location label")
	ErrUndefDesc            = errors.New("undefined location description")
	ErrConnUndefLabel       = errors.New("undefined connection label")
	ErrConnLocationNotFound = errors.New("connection location not found")
	ErrConnWordNotFound     = errors.New("connection word not found")
)
