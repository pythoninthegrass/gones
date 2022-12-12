package apu

import "github.com/gabe565/gones/internal/memory"

type CPU interface {
	memory.Read8
	AddStall(uint8)
}

var dmcTicks = [...]byte{
	214, 190, 170, 160, 143, 127, 113, 107, 95, 80, 71, 64, 53, 42, 36, 27,
}

type DMC struct {
	Enabled bool
	Value   byte

	IrqEnabled bool
	Loop       bool
	LoadCount  byte

	TickPeriod byte
	TickValue  byte

	SampleAddr uint16
	SampleLen  uint16

	cpu           CPU
	CurrAddr      uint16
	CurrLen       uint16
	ShiftRegister byte
	BitCount      byte
}

func (d *DMC) Write(addr uint16, data byte) {
	switch addr {
	case 0x4010:
		d.IrqEnabled = data>>7&1 == 1
		d.Loop = data>>6&1 == 1
		d.TickPeriod = dmcTicks[data&0xF]
	case 0x4011:
		d.LoadCount = data & 0x7F
	case 0x4012:
		d.SampleAddr = 0xC000 | uint16(data)<<6
	case 0x4013:
		d.SampleLen = uint16(data)<<4 | 1
	}
}

func (d *DMC) SetEnabled(v bool) {
	d.Enabled = v
	if v {
		if d.CurrLen == 0 {
			d.restart()
		}
	} else {
		d.CurrLen = 0
	}
}

func (d *DMC) restart() {
	d.CurrAddr = d.SampleAddr
	d.CurrLen = d.SampleLen
}

func (d *DMC) stepTimer() {
	if !d.Enabled {
		return
	}
	d.stepReader()
	if d.TickValue == 0 {
		d.TickValue = d.TickPeriod
		d.stepShifter()
	} else {
		d.TickValue -= 1
	}
}

func (d *DMC) stepReader() {
	if d.CurrLen > 0 && d.BitCount == 0 {
		d.cpu.AddStall(4)
		d.ShiftRegister = d.cpu.ReadMem(d.CurrAddr)
		d.BitCount = 8
		d.CurrAddr += 1
		if d.CurrAddr == 0 {
			d.CurrAddr = 0x8000
		}
		d.CurrLen -= 1
		if d.CurrLen == 0 && d.Loop {
			d.restart()
		}
	}
}

func (d *DMC) stepShifter() {
	if d.BitCount == 0 {
		return
	}
	if d.ShiftRegister&1 == 1 {
		if d.Value <= 125 {
			d.Value += 2
		}
	} else {
		if d.Value >= 2 {
			d.Value -= 2
		}
	}
	d.ShiftRegister >>= 1
	d.BitCount -= 1
}

func (d *DMC) output() byte {
	return d.Value
}
