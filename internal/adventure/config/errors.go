package config

import "errors"

var (
	ErrCannotParseDate         = errors.New("cannot parse date")
	ErrUnrecognizedLanguage    = errors.New("unrecognized language")
	ErrUnrecognizedConfigLabel = errors.New("unrecognized config label")
)
