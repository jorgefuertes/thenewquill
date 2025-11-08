package adventure

import (
	"io"

	"github.com/jorgefuertes/thenewquill/internal/bin"
)

func (a *Adventure) Export(path string) (int64, error) {
	if err := a.Validate(); err != nil {
		return 0, err
	}

	b := bin.New()
	for _, s := range a.DB.Data {
		if err := b.WriteStoreable(s); err != nil {
			return 0, err
		}
	}

	return b.Save(path)
}

func (a *Adventure) Import(r io.Reader) error {
	a.DB.Reset()

	return nil
}
