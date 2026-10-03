// Package font loads a default Pure-Go TTF for gg text drawing.
//
// Games and tools share one FontSource; call Apply or Face per size.
package font

import (
	"sync"

	"github.com/gogpu/gg"
	"github.com/gogpu/gg/text"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	once    sync.Once
	source  *text.FontSource
	errInit error
)

// Bundle owns a FontSource. Prefer Default() unless you need a custom TTF.
type Bundle struct {
	src *text.FontSource
}

// Default returns the shared Go Regular bundle (lazy-init, Pure Go).
func Default() (*Bundle, error) {
	once.Do(func() {
		source, errInit = text.NewFontSource(goregular.TTF)
	})
	if errInit != nil {
		return nil, errInit
	}
	return &Bundle{src: source}, nil
}

// MustDefault is Default that panics on failure (handy in main).
func MustDefault() *Bundle {
	b, err := Default()
	if err != nil {
		panic(err)
	}
	return b
}

// New creates a bundle from raw TTF/OTF bytes.
func New(ttf []byte) (*Bundle, error) {
	src, err := text.NewFontSource(ttf)
	if err != nil {
		return nil, err
	}
	return &Bundle{src: src}, nil
}

// Face returns a face at the given pixel size.
func (b *Bundle) Face(size float64) text.Face {
	if b == nil || b.src == nil {
		return nil
	}
	return b.src.Face(size)
}

// Apply sets dc's current font. No-op if bundle or face is unavailable.
func (b *Bundle) Apply(dc *gg.Context, size float64) {
	if dc == nil || b == nil {
		return
	}
	face := b.Face(size)
	if face == nil {
		return
	}
	dc.SetFont(face)
}

// Apply is a convenience on the default bundle.
func Apply(dc *gg.Context, size float64) {
	b, err := Default()
	if err != nil {
		return
	}
	b.Apply(dc, size)
}
