// Code generated by "stringer -type AddressingMode"; DO NOT EDIT.

package cpu

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Implied-0]
	_ = x[Accumulator-1]
	_ = x[Immediate-2]
	_ = x[ZeroPage-3]
	_ = x[ZeroPageX-4]
	_ = x[ZeroPageY-5]
	_ = x[Relative-6]
	_ = x[Absolute-7]
	_ = x[AbsoluteX-8]
	_ = x[AbsoluteY-9]
	_ = x[Indirect-10]
	_ = x[IndirectX-11]
	_ = x[IndirectY-12]
}

const _AddressingMode_name = "ImpliedAccumulatorImmediateZeroPageZeroPageXZeroPageYRelativeAbsoluteAbsoluteXAbsoluteYIndirectIndirectXIndirectY"

var _AddressingMode_index = [...]uint8{0, 7, 18, 27, 35, 44, 53, 61, 69, 78, 87, 95, 104, 113}

func (i AddressingMode) String() string {
	if i >= AddressingMode(len(_AddressingMode_index)-1) {
		return "AddressingMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _AddressingMode_name[_AddressingMode_index[i]:_AddressingMode_index[i+1]]
}
