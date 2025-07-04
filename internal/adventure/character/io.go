package character

import (
	"fmt"
	"regexp"

	"github.com/jorgefuertes/thenewquill/internal/util"
)

func (c Character) Export() string {
	return fmt.Sprintf("%d|%d|%d|%d|%s|%d|%d|%d\n",
		c.GetKind().Byte(),
		c.ID,
		c.NounID,
		c.AdjectiveID,
		util.EscapeExportString(c.Description),
		c.LocationID,
		util.BoolToInt(c.Created),
		util.BoolToInt(c.Human),
	)
}

func (s *Service) Import(line string) error {
	rg := regexp.MustCompile(`^(\d+)\|(\d+)\|(\d+)\|(\d+)\|(.+)\|(\d+)\|(\d+)\|(\d+)$`)

	if !rg.MatchString(line) {
		return ErrCannotImport
	}

	return nil
}
