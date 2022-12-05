package cartridge

//go:generate stringer -type Mirror

type Mirror byte

const (
	Horizontal Mirror = iota
	Vertical
	SingleLower
	SingleUpper
	FourScreen
)
