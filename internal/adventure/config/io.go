package config

import (
	"thenewquill/internal/compiler/db"
	"thenewquill/internal/compiler/section"
)

func (c Config) Export(d *db.DB) {
	d.Append(section.Config, section.Config.String(),
		c.Title,
		c.Author,
		c.Description,
		c.Version,
		c.Date,
		c.Lang.Byte(),
	)
}

func (c *Config) Import(d *db.DB) error {
	r := d.GetRecord(section.Config, section.Config.String())

	if r == nil {
		return ErrMissingConfigRecord
	}

	c.Title = r.FieldAsString(0)
	c.Author = r.FieldAsString(1)
	c.Description = r.FieldAsString(2)
	c.Version = r.FieldAsString(3)
	c.Date = r.FieldAsString(4)
	c.Lang = Lang(r.FieldAsByte(5))

	return nil
}
