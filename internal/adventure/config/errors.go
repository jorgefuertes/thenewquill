package config

import "errors"

var (
	ErrCannotParseDate         = errors.New("cannot parse date")
	ErrUnrecognizedLanguage    = errors.New("unrecognized language")
	ErrUnrecognizedConfigLabel = errors.New("unrecognized config label")
	ErrMissingConfigField      = errors.New("missing config field")
	ErrMissingConfigRecord     = errors.New("missing config record")
)
