package config

import "errors"

var (
	ErrUnrecognizedLanguage    = errors.New("unrecognized language")
	ErrUnrecognizedConfigField = errors.New("unrecognized config field")
	ErrMissingConfigField      = errors.New("missing config field")
)
