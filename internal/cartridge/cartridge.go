package cartridge

import (
	"crypto/md5"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gabe565/gones/internal/consts"
	"github.com/gabe565/gones/internal/database"
	"github.com/gabe565/gones/internal/interrupt"
)

type Cartridge struct {
	hash string
	name string

	prg       []byte
	Chr       []byte
	Sram      []byte
	Mapper    byte `msgpack:"-"`
	Submapper byte `msgpack:"-"`
	Mirror    Mirror
	Battery   bool `msgpack:"-"`
}

func New() *Cartridge {
	return &Cartridge{
		Sram: make([]byte, 0x2000),
	}
}

func FromBytes(b []byte) *Cartridge {
	cart := New()
	cart.hash = fmt.Sprintf("%x", md5.Sum(b))
	if cart.hash != "" {
		cart.name, _ = database.FindNameByHash(cart.hash)
	}

	cart.prg = make([]byte, consts.PrgRomAddr, consts.PrgChunkSize*2)
	cart.prg = append(cart.prg, b...)
	cart.prg = cart.prg[:cap(cart.prg)]
	cart.prg[interrupt.ResetVector+1-consts.PrgChunkSize*2] = 0x86

	cart.Chr = make([]byte, consts.ChrChunkSize)

	return cart
}

func (c *Cartridge) Name() string {
	return c.name
}

func (c *Cartridge) SetName(path string) {
	c.name = strings.TrimSuffix(filepath.Base(path), ".nes")
}

func (c *Cartridge) Hash() string {
	return c.hash
}
