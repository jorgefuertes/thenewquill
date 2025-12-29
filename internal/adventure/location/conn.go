package location

type Connection struct {
	WordID     uint32 `valid:"required"`
	LocationID uint32 `valid:"required"`
}
