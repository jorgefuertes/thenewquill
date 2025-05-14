//go:build manual
// +build manual

package console_test

import (
	"testing"
	"time"

	"github.com/jorgefuertes/thenewquill/internal/output/console"

	lorem "github.com/drhodes/golorem"
	"github.com/stretchr/testify/require"
)

func TestConsole(t *testing.T) {
	c, err := console.New()
	require.NoError(t, err)
	c.Delay = time.Millisecond * 2
	go c.Run()

	t.Run("Scroll", func(t *testing.T) {
		for i := 0; i <= 100; i++ {
			c.Printf("[%04d] Hello, World!\n", i)
		}

		time.Sleep(time.Second * 2)

		c.Cls()
	})

	t.Run("Long text", func(t *testing.T) {
		for i := 0; i <= 25; i++ {
			c.WrapPrintf("%s\n\n", lorem.Sentence(20, 80))
		}

		time.Sleep(time.Second * 2)

		c.Cls()
	})

	t.Run("Quick text", func(t *testing.T) {
		c.Delay = time.Millisecond * 0

		for i := 0; i <= 25; i++ {
			c.WrapPrintf("%s\n\n", lorem.Sentence(20, 80))
		}

		time.Sleep(time.Second * 2)
	})

	c.Close()
}
