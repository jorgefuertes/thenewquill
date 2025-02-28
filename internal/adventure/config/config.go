package config

import (
	"errors"
	"time"
)

type Config struct {
	Title       string
	Author      string
	Description string
	Version     string
	Date        time.Time
	Lang        Lang
}

func New() Config {
	return Config{}
}

func (c *Config) Set(label Label, value string) error {
	switch label {
	case TitleLabel:
		c.Title = value
	case AuthorLabel:
		c.Author = value
	case DescriptionLabel, DescLabel:
		c.Description = value
	case VersionLabel:
		c.Version = value
	case DateLabel:
		var err error
		c.Date, err = time.Parse("02-01-2006", value)
		if err != nil {
			c.Date, err = time.Parse("2006-01-02", value)
			if err != nil {
				return errors.Join(ErrCannotParseDate, err)
			}
		}
	case LangLabel:
		lang := LangFromString(value)
		if lang == Undefined {
			return ErrUnrecognizedLanguage
		}

		c.Lang = lang
	default:
		return ErrUnrecognizedConfigLabel
	}

	return nil
}
