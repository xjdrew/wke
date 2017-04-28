package webwindow

import (
	"github.com/lxn/win"
	"github.com/xjdrew/wke"
)

func lParamToKeyFlags(lParam uintptr) wke.KeyFlags {
	var flags wke.KeyFlags

	v := win.HIWORD(uint32(lParam))
	// KF_REPEAT
	if v&0x4000 != 0 {
		flags |= wke.KF_REPEAT
	}
	// KF_EXTENDED
	if v&0x0100 != 0 {
		flags |= wke.KF_EXTENDED
	}
	return flags
}

func wParamToMouseFlags(wParam uintptr) wke.MouseFlags {
	var flags wke.MouseFlags

	v := uint32(wParam)

	if v&win.MK_CONTROL != 0 {
		flags |= wke.MF_CONTROL
	}

	if v&win.MK_SHIFT != 0 {
		flags |= wke.MF_SHIFT
	}

	if v&win.MK_LBUTTON != 0 {
		flags |= wke.MF_LBUTTON
	}

	if v&win.MK_MBUTTON != 0 {
		flags |= wke.MF_MBUTTON
	}

	if v&win.MK_RBUTTON != 0 {
		flags |= wke.MF_RBUTTON
	}

	return flags
}

func wParamToWheelDelta(wParam uintptr) int {
	v := int16(win.HIWORD(uint32(wParam)))
	return int(v)
}
