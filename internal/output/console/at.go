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

func (c *console) toBounds(row, col int) (int, int) {
	if row < 0 {
		row = 0
	}

	if row >= c.Rows() {
		row = c.Rows() - 1
	}

	if col < 0 {
		col = 0
	}

	if col >= c.Cols() {
		col = c.Cols() - 1
	}

	return row, col
}

// At moves the cursor to the given position
func (c *console) At(row, col int) {
	c.at.Col = col
	c.at.Row = row
	c.at.Row, c.at.Col = c.toBounds(c.at.Row, c.at.Col)
}

/*
moveCursor moves the cursor n-cols to the right (positive) or left (negative):

  - Jumps to the next line if the cursor is at the end of the line and n is positive.
  - Jumps back to the previous line if the cursor is at the beginning of the line and n is negative.
  - Scrolls the screen if the cursor is at the end of the screen and n is positive.
*/
func (c *console) moveCursor(n int) {
	c.at.Col += n
	if c.at.Col < 0 {
		if c.at.Row > 0 {
			c.at.Row--
			c.at.Col = c.Cols() - 1
		} else {
			c.at.Col = 0
		}
	}

	if c.at.Col >= c.Cols() {
		c.incRow()
	}
}

func (c *console) incRow() {
	c.at.Row++
	c.at.Col = 0
	if c.at.Row >= c.Rows() {
		c.at.Row = c.Rows() - 1
		c.scroll()
	}
}

func (c *console) PushAt() {
	c.atStack = append(c.atStack, c.at)
}

func (c *console) PopAt() {
	if len(c.atStack) > 0 {
		c.at = c.atStack[len(c.atStack)-1]
		c.atStack = c.atStack[:len(c.atStack)-1]
	}
}
