package config

type Config struct {
	Title       string
	Author      string
	Description string
	Version     string
	Date        string
	Lang        Lang
}

func New() *Config {
	c := &Config{}
	c.Reset()

	return c
}

func (c *Config) Reset() {
	c.Title = "No Title"
	c.Author = "Unknown Author"
	c.Description = "No Description"
	c.Version = "1.0.0"
	c.Date = "2025-01-01"
	c.Lang = ES
}

func (c *Config) Set(field Field, value string) error {
	switch field {
	case TitleField:
		c.Title = value
	case AuthorField:
		c.Author = value
	case DescriptionField, DescField:
		c.Description = value
	case VersionField:
		c.Version = value
	case DateField:
		c.Date = value
	case LangField:
		lang := LangFromString(value)
		if lang == Undefined {
			return ErrUnrecognizedLanguage
		}

		c.Lang = lang
	default:
		return ErrUnrecognizedConfigField
	}

	return nil
}
