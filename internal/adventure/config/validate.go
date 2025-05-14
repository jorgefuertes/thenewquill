package config

import "errors"

func (c *Config) Validate() error {
	if c.Title == "" {
		return errors.Join(ErrMissingConfigField, errors.New("title"))
	}

	if c.Author == "" {
		return errors.Join(ErrMissingConfigField, errors.New("author"))
	}

	if c.Description == "" {
		return errors.Join(ErrMissingConfigField, errors.New("description"))
	}

	if c.Version == "" {
		return errors.Join(ErrMissingConfigField, errors.New("version"))
	}

	if c.Date == "" {
		return errors.Join(ErrMissingConfigField, errors.New("date"))
	}

	if c.Lang == Undefined {
		return errors.Join(ErrMissingConfigField, errors.New("lang"))
	}

	return nil
}
