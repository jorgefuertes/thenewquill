package loc

import "errors"

var ErrDuplicatedLocation = errors.New("duplicated location")
var ErrDuplicatedConnection = errors.New("duplicated connection")
var ErrLocationNotFound = errors.New("location not found")
