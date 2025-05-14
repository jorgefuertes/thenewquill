package location

import (
	"errors"
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

func (l Location) Validate() error {
	if l.ID == db.UndefinedLabel.ID {
		return ErrUndefLabel
	}

	if l.ID < db.MinMeaningfulID {
		return ErrWrongLabel
	}

	if l.Description == "" || l.Description == Undefined {
		return ErrUndefDesc
	}

	return nil
}

func ValidateAll(d *db.DB) error {
	locations := make([]Location, 0)
	if err := d.GetByKindAs(db.Locations, db.NoSubKind, &locations); err != nil {
		return err
	}

	for _, loc := range locations {
		if err := loc.Validate(); err != nil {
			return errors.Join(err, fmt.Errorf("label: %s", d.GetLabelName(loc.ID)))
		}

		for i, conn := range loc.Conns {
			if !d.Exists(conn.WordID, db.Words, db.NoSubKind) {
				return errors.Join(
					ErrWrongLabel,
					fmt.Errorf("location: '%s', conn '%d', word '%d'", d.GetLabelName(loc.ID), i, conn.WordID),
				)
			}

			if !d.Exists(conn.LocationID) {
				return errors.Join(
					ErrWrongLabel,
					fmt.Errorf(
						"location: '%s', conn '%s' to location '%d'",
						d.GetLabelName(loc.ID),
						d.GetLabelName(conn.WordID),
						conn.LocationID,
					),
				)
			}
		}
	}

	return nil
}
