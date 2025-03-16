package config

import (
	"thenewquill/internal/compiler/db"
	"thenewquill/internal/compiler/section"
)

func (c Config) Export(d *db.DB) {
	db.NewRegister(section.Config, section.Config.String(),
		c.Title,
		c.Author,
		c.Description,
		c.Version,
		c.Date,
		c.Lang.Int(),
	)
}

func (c *Config) Import(d *db.DB) {
	regs := d.GetRegsForSection(section.Config)

	c.Title = regs[0].GetString()
	c.Author = regs[0].GetString()
	c.Description = regs[0].GetString()
	c.Version = regs[0].GetString()
	c.Date = regs[0].GetString()
	c.Lang = Lang(regs[0].GetInt())
}
