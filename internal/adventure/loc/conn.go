package loc

import "thenewquill/internal/adventure/words"

type Connection struct {
	Word *words.Word
	To   *Location
}
