//go:build manual

package console_test

import (
	"testing"
	"time"

	"github.com/jorgefuertes/thenewquill/internal/output/console"

	"github.com/stretchr/testify/require"
)

func TestConsoleInput(t *testing.T) {
	c, err := console.New()
	require.NoError(t, err)
	defer c.Close()

	c.Delay = time.Millisecond * 6
	go c.Run()

	t.Run("Input", func(t *testing.T) {
		for {
			prompt := "> Type something (x to exit):"
			res, err := c.Input(prompt, 0)
			if err != nil {
				t.Logf("Error: %s", err)
				break
			}

			if res == "x" {
				break
			}

			c.Printf("You typed: '%s'\n", res)
		}
	})
}
