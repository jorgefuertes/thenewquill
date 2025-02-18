package player

import "errors"

var ErrDuplicatedPlayerLabel = errors.New("duplicated player label")
var ErrOnlyOneHuman = errors.New("only one human can be defined")
