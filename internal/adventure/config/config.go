package config

type Config struct {
	Title       string
	Author      string
	Description string
	Version     string
	Date        string
	Lang        Lang
}

func New() Config {
	c := Config{}
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
