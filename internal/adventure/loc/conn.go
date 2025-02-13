package loc

import "thenewquill/internal/adventure/voc"

type Connection struct {
	Word *voc.Word
	To   *Location
}
