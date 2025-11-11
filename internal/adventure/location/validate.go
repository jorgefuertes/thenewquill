package location

import (
	"errors"
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/adventure/id"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

func (l Location) Validate(allowNoID bool) error {
	if err := l.ID.Validate(false); err != nil && !allowNoID {
		if err == db.ErrInvalidLabelID && !allowNoID {
			return err
		}
	}

	if l.ID < id.Min && !allowNoID {
		return ErrWrongLabel
	}

	if l.Description == "" || l.Description == Undefined {
		return ErrUndefDesc
	}

	for _, conn := range l.Conns {
		if conn.WordID == id.Undefined {
			return ErrConnUndefLabel
		}
	}

	return nil
}

func (s *Service) ValidateAll() error {
	locations := s.db.Query(db.FilterByKind(kind.Location))
	defer locations.Close()

	var loc Location
	for locations.Next(&loc) {
		if err := loc.Validate(false); err != nil {
			return errors.Join(err, fmt.Errorf("label: %s", s.db.GetLabelName(loc.ID)))
		}

		for i, conn := range loc.Conns {
			if !s.db.Exists(db.FilterByID(conn.WordID)) {
				return errors.Join(
					ErrConnWordNotFound,
					fmt.Errorf(
						"conn %02d: %s(%d):%s(%d)->%s(%d)",
						i,
						s.db.GetLabelName(loc.ID),
						loc.ID,
						s.db.GetLabelName(conn.WordID),
						conn.WordID,
						s.db.GetLabelName(conn.LocationID),
						conn.LocationID,
					),
				)
			}

			if !s.db.Exists(db.FilterByID(conn.WordID)) {
				return errors.Join(
					ErrConnLocationNotFound,
					fmt.Errorf(
						"conn %02d: %s(%d):%s(%d)->%s(%d)",
						i,
						s.db.GetLabelName(loc.ID),
						loc.ID,
						s.db.GetLabelName(conn.WordID),
						conn.WordID,
						s.db.GetLabelName(conn.LocationID),
						conn.LocationID,
					),
				)
			}
		}
	}

	return nil
}
