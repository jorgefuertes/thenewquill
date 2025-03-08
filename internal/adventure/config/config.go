package config

import "thenewquill/internal/compiler/section"

type Config struct {
	Title       string
	Author      string
	Description string
	Version     string
	Date        string
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
		c.Date = value
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

func (c Config) Export() (section.Section, map[string]any) {
	return section.Config, map[string]any{
		"title":       c.Title,
		"author":      c.Author,
		"description": c.Description,
		"version":     c.Version,
		"date":        c.Date,
		"lang":        c.Lang,
	}
}
