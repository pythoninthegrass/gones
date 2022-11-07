package cpu

import (
	"errors"
	"fmt"
	"github.com/gabe565/gones/internal/bitflags"
	"github.com/gabe565/gones/internal/bus"
	"github.com/gabe565/gones/internal/consts"
)

func New(b *bus.Bus) *CPU {
	return &CPU{
		Status:       DefaultStatus,
		StackPointer: StackReset,
		Bus:          b,
	}
}

// CPU implements the NES CPU.
//
// See [6502 Guide].
//
// [6502 Guide]: https://www.nesdev.org/obelisk-6502-guide/
type CPU struct {
	// ProgramCounter Program Counter
	ProgramCounter uint16

	// StackPointer Stack Pointer
	StackPointer byte

	// Status Processor Status
	Status bitflags.Flags

	// Accumulator Register A
	Accumulator byte

	// RegisterX Register X
	RegisterX byte

	// RegisterY Register Y
	RegisterY byte

	// Bus Main memory bus
	Bus *bus.Bus

	// Callback optional callback to Run before every tick
	Callback func(c *CPU) error

	// Debug enables opcode logging
	Debug bool
}

const (
	// StackAddr is the memory address of the stack
	StackAddr = 0x100

	// StackReset is the start value for the stack pointer
	StackReset = 0xFD
)

// Reset resets the CPU and sets ProgramCounter to the value of the [Reset] Vector.
func (c *CPU) Reset() {
	c.Accumulator = 0
	c.RegisterX = 0
	c.Status = DefaultStatus
	c.StackPointer = StackReset

	c.ProgramCounter = c.MemRead16(consts.ResetAddr)
}

// Load loads a program into PRG memory
func (c *CPU) Load(program []byte) {
	for k, v := range program {
		c.MemWrite(consts.PrgRomAddr+uint16(k), v)
	}
	c.MemWrite16(consts.ResetAddr, consts.PrgRomAddr)
}

// ErrUnsupportedOpcode indicates an unsupported opcode was evaluated.
var ErrUnsupportedOpcode = errors.New("unsupported opcode")

// Run is the main Run entrypoint.
func (c *CPU) Run() error {
	opcodes := OpCodeMap()

	for {
		if c.Callback != nil {
			if err := c.Callback(c); err != nil {
				if errors.Is(err, ErrBrk) {
					return nil
				}
				return err
			}
		}

		code := c.MemRead(c.ProgramCounter)
		c.ProgramCounter += 1
		prevPC := c.ProgramCounter

		op, ok := opcodes[code]
		if !ok {
			return fmt.Errorf("%w: $%02X", ErrUnsupportedOpcode, code)
		}

		if c.Debug {
			fmt.Println(op)
		}

		if err := op.Exec(c, op.Mode); err != nil {
			if errors.Is(err, ErrBrk) {
				return nil
			}
			return err
		}

		if prevPC == c.ProgramCounter {
			c.ProgramCounter += uint16(op.Len - 1)
		}
	}
}
