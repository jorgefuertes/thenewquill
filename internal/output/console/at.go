package console

// Size returns the size of the console in columns and rows.
func (c *console) Size() (int, int) {
	return c.screen.Size()
}

// Rows returns the number of rows (height) in the console.
func (c *console) Rows() int {
	_, rows := c.screen.Size()

	return rows
}

// Cols returns the number of columns (width) in the console.
func (c *console) Cols() int {
	cols, _ := c.screen.Size()

	return cols
}
