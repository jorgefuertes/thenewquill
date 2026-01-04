package location

import (
	"errors"
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
)

func (s *Service) PostReplace() error {
	c := s.db.Query(database.FilterByKind(kind.Location))
	defer c.Close()

	var loc Location
	for c.Next(&loc) {
		dirty := false

		for i, conn := range loc.Conns {
			dst, err := s.Get().WithLabelID(conn.LocationID).First()
			if err != nil {
				return fmt.Errorf(
					"%w: connection error [SRCLOC:%d:%s] [CONN:#%d:%s] [DSTLOC:%d:%s]",
					err,
					loc.ID,
					s.db.GetLabelOrBlank(loc.LabelID),
					i,
					s.db.GetLabelOrBlank(conn.WordID),
					conn.LocationID,
					s.db.GetLabelOrBlank(conn.LocationID),
				)
			}

			loc.SetConn(conn.WordID, dst.ID)
			dirty = true
		}

		if !dirty {
			continue
		}

		if err := s.db.Update(&loc); err != nil {
			return errors.Join(
				err,
				fmt.Errorf("cannot update [SRCLOC:%d:%s]", loc.ID, s.db.GetLabelOrBlank(loc.LabelID)),
			)
		}
	}

	return nil
}
