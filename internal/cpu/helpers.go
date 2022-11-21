package cpu

import "github.com/gabe565/gones/internal/bitflags"

// updateZeroAndNegFlags updates zero and negative flags
func (c *CPU) updateZeroAndNegFlags(result byte) {
	c.Status.Set(Zero, result == 0)
	c.Status.Set(Negative, bitflags.Flags(result).Intersects(Negative))
}

func (c *CPU) branch(condition bool) {
	if condition {
		c.Cycles += 1

		jump := int8(c.ReadMem(c.ProgramCounter))
		jumpAddr := c.ProgramCounter + 1 + uint16(jump)

		if (c.ProgramCounter+1)&0xFF0 != jumpAddr {
			c.Cycles += 1
		}

		c.ProgramCounter = jumpAddr
	}
}

func (c *CPU) compare(mode AddressingMode, rhs byte) {
	addr, pageCrossed := c.getOperandAddress(mode)
	if pageCrossed {
		defer func() {
			c.Cycles += 1
		}()
	}
	data := c.ReadMem(addr)
	c.Status.Set(Carry, data <= rhs)
	c.updateZeroAndNegFlags(rhs - data)
}

func (c *CPU) setAccumulator(v byte) {
	c.Accumulator = v
	c.updateZeroAndNegFlags(c.Accumulator)
}

func (c *CPU) addAccumulator(data byte) {
	sum := uint16(c.Accumulator) + uint16(data)
	if c.Status.Intersects(Carry) {
		sum += 1
	}

	carry := sum > 0xFF
	c.Status.Set(Carry, carry)

	result := byte(sum)
	c.Status.Set(Overflow, (data^result)&(result^c.Accumulator)&0x80 != 0)

	c.setAccumulator(result)
}
