package location

import (
	"errors"
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/database"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

func (l Location) Validate(allowNoID bool) error {
	if err := l.ID.ValidateID(false); err != nil && !allowNoID {
		if err == primitive.ErrInvalidID && !allowNoID {
			return err
		}
	}

	if l.ID < primitive.MinID && !allowNoID {
		return ErrWrongLabel
	}

	if l.Description == "" || l.Description == Undefined {
		return ErrUndefDesc
	}

	for _, conn := range l.Conns {
		if conn.WordID == primitive.UndefinedID {
			return ErrConnWordUndefinedID
		}

		if conn.LocationID == primitive.UndefinedID {
			return ErrConnLocationUndefinedID
		}
	}

	return nil
}

func (s *Service) ValidateAll() error {
	locations := s.db.Query(database.FilterByKind(kind.Location))
	defer locations.Close()

	var loc Location
	for locations.Next(&loc) {
		if err := loc.Validate(false); err != nil {
			return errors.Join(err, fmt.Errorf("location %d: %s", loc.ID, s.db.GetLabelOrBlank(loc.LabelID)))
		}

		for i, conn := range loc.Conns {
			wordLabel, _ := s.db.GetLabelForStoreable(conn.WordID)
			dstLocLabel, _ := s.db.GetLabelForStoreable(conn.LocationID)
			connErr := fmt.Errorf(
				"conn %02d: %s(%d):%s(%d)->%s(%d)",
				i,
				s.db.GetLabelOrBlank(loc.LabelID),
				loc.ID,
				wordLabel,
				conn.WordID,
				dstLocLabel,
				conn.LocationID,
			)

			if !s.db.Exists(database.FilterByID(conn.WordID), database.FilterByKind(kind.Word)) {
				return errors.Join(ErrConnWordNotFound, connErr)
			}

			if !s.db.Exists(database.FilterByID(conn.LocationID), database.FilterByKind(kind.Location)) {
				return errors.Join(ErrConnLocationNotFound, connErr)
			}
		}
	}

	return nil
}
