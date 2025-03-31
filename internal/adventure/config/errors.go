package config

import "errors"

var (
	ErrCannotParseDate         = errors.New("cannot parse date")
	ErrUnrecognizedLanguage    = errors.New("unrecognized language")
	ErrUnrecognizedConfigField = errors.New("unrecognized config field")
	ErrMissingConfigField      = errors.New("missing config field")
	ErrMissingConfigRecord     = errors.New("missing config record")
)
