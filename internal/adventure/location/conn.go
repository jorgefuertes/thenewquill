package location

import "github.com/jorgefuertes/thenewquill/internal/adventure/db"

type Connection struct {
	WordID     db.ID
	LocationID db.ID
}
