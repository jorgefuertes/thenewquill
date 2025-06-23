package location

import (
	"errors"
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

func (l Location) Validate(allowNoID db.Allow) error {
	if err := l.ID.Validate(db.DontAllowSpecial); err != nil && !allowNoID {
		return err
	}

	if l.ID < db.MinMeaningfulID {
		return ErrWrongLabel
	}

	if l.Description == "" || l.Description == Undefined {
		return ErrUndefDesc
	}

	for _, conn := range l.Conns {
		if conn.WordID == db.UndefinedLabel.ID {
			return ErrConnUndefLabel
		}
	}

	return nil
}

func (s *Service) ValidateAll() error {
	locations := s.db.Query(db.Locations)
	defer locations.Close()

	var loc Location
	for locations.Next(&loc) {
		if err := loc.Validate(db.DontAllowNoID); err != nil {
			return errors.Join(err, fmt.Errorf("label: %s", s.db.GetLabelName(loc.ID)))
		}

		for i, conn := range loc.Conns {
			if !s.db.Exists(conn.WordID, db.Words) {
				return errors.Join(
					ErrWrongLabel,
					fmt.Errorf("location: '%s', conn '%d', word '%d'", s.db.GetLabelName(loc.ID), i, conn.WordID),
				)
			}

			if !s.db.Exists(conn.LocationID, db.Locations) {
				return errors.Join(
					ErrWrongLabel,
					fmt.Errorf(
						"location: '%s', conn '%s' to location '%d'",
						s.db.GetLabelName(loc.ID),
						s.db.GetLabelName(conn.WordID),
						conn.LocationID,
					),
				)
			}
		}
	}

	return nil
}
