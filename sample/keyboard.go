package main

import (
	"github.com/deluan/bring"
	"github.com/faiface/pixel/pixelgl"
)

var (
	keys map[pixelgl.Button]bring.KeyCode
)

// Rant: why pixelgl keyboard events handling is so messy?!?
func collectKeyStrokes(win *pixelgl.Window) (pressed []bring.KeyCode, released []bring.KeyCode) {
	for k, v := range keys {
		key := v
		if win.JustPressed(k) || win.Repeated(k) {
			pressed = append(pressed, key)
		}
		if win.JustReleased(k) {
			released = append(released, key)
		}
	}
	controlPressed := win.Pressed(pixelgl.KeyLeftControl) || win.Pressed(pixelgl.KeyRightControl) ||
		win.Pressed(pixelgl.KeyLeftAlt) || win.Pressed(pixelgl.KeyRightAlt)
	if controlPressed {
		shiftPressed := win.Pressed(pixelgl.KeyLeftShift) || win.Pressed(pixelgl.KeyRightShift)
		for ch := 32; ch < 127; ch++ {
			isLetter := ch >= int('A') && ch <= int('Z')
			key := ch
			if isLetter && !shiftPressed {
				key = ch + 32
			}
			if win.JustPressed(pixelgl.Button(ch)) || win.Repeated(pixelgl.Button(ch)) {
				pressed = append(pressed, bring.KeyCode(key))
			}
			if win.JustReleased(pixelgl.Button(ch)) {
				released = append(released, bring.KeyCode(key))
			}
		}
	} else {
		for _, ch := range win.Typed() {
			pressed = append(pressed, bring.KeyCode(int(ch)))
			released = append(released, bring.KeyCode(int(ch)))
		}
	}
	return
}

func init() {
	keys = map[pixelgl.Button]bring.KeyCode{
		pixelgl.KeyLeftAlt:      bring.KeyLeftAlt,
		pixelgl.KeyRightAlt:     bring.KeyRightAlt,
		pixelgl.KeyLeftControl:  bring.KeyLeftControl,
		pixelgl.KeyRightControl: bring.KeyRightControl,
		pixelgl.KeyLeftShift:    bring.KeyLeftShift,
		pixelgl.KeyRightShift:   bring.KeyRightShift,
		pixelgl.KeyBackspace:    bring.KeyBackspace,
		pixelgl.KeyCapsLock:     bring.KeyCapsLock,
		pixelgl.KeyDelete:       bring.KeyDelete,
		pixelgl.KeyDown:         bring.KeyDown,
		pixelgl.KeyEnd:          bring.KeyEnd,
		pixelgl.KeyEnter:        bring.KeyEnter,
		pixelgl.KeyEscape:       bring.KeyEscape,
		pixelgl.KeyF1:           bring.KeyF1,
		pixelgl.KeyF2:           bring.KeyF2,
		pixelgl.KeyF3:           bring.KeyF3,
		pixelgl.KeyF4:           bring.KeyF4,
		pixelgl.KeyF5:           bring.KeyF5,
		pixelgl.KeyF6:           bring.KeyF6,
		pixelgl.KeyF7:           bring.KeyF7,
		pixelgl.KeyF8:           bring.KeyF8,
		pixelgl.KeyF9:           bring.KeyF9,
		pixelgl.KeyF10:          bring.KeyF10,
		pixelgl.KeyF11:          bring.KeyF11,
		pixelgl.KeyF12:          bring.KeyF12,
		pixelgl.KeyF13:          bring.KeyF13,
		pixelgl.KeyF14:          bring.KeyF14,
		pixelgl.KeyF15:          bring.KeyF15,
		pixelgl.KeyF16:          bring.KeyF16,
		pixelgl.KeyF17:          bring.KeyF17,
		pixelgl.KeyF18:          bring.KeyF18,
		pixelgl.KeyF19:          bring.KeyF19,
		pixelgl.KeyF20:          bring.KeyF20,
		pixelgl.KeyF21:          bring.KeyF21,
		pixelgl.KeyF22:          bring.KeyF22,
		pixelgl.KeyF23:          bring.KeyF23,
		pixelgl.KeyF24:          bring.KeyF24,
		pixelgl.KeyHome:         bring.KeyHome,
		pixelgl.KeyInsert:       bring.KeyInsert,
		pixelgl.KeyLeft:         bring.KeyLeft,
		pixelgl.KeyNumLock:      bring.KeyNumLock,
		pixelgl.KeyPageDown:     bring.KeyPageDown,
		pixelgl.KeyPageUp:       bring.KeyPageUp,
		pixelgl.KeyPause:        bring.KeyPause,
		pixelgl.KeyPrintScreen:  bring.KeyPrintScreen,
		pixelgl.KeyRight:        bring.KeyRight,
		pixelgl.KeyTab:          bring.KeyTab,
		pixelgl.KeyUp:           bring.KeyUp,
		// pixelgl.KeyMeta:         bring.KeyMeta,
		// pixelgl.KeySuper:        bring.KeySuper,
		// pixelgl.KeyWin:          bring.KeyWin,
	}
}
