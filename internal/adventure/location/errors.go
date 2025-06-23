package location

import "errors"

var (
	ErrWrongLabel     = errors.New("wrong location label")
	ErrUndefLabel     = errors.New("undefined location label")
	ErrUndefDesc      = errors.New("undefined location description")
	ErrConnUndefLabel = errors.New("undefined connection label")
)
