package main

import (
	"github.com/deluan/bring"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"github.com/veandco/go-sdl2/sdl"
)

type keyboardState struct {
	client *bring.Client
	shift  bool
	caps   bool
}

func hookKeyboard(win *sdlcanvas.Window, client *bring.Client) {
	ks := &keyboardState{client: client}
	win.KeyUp = ks.keyUp
	win.KeyDown = ks.keyDown
}

func (ks *keyboardState) sendKey(key bring.KeyCode, pressed bool) {
	k := int(key)
	if k >= int('A') && k <= int('Z') {
		if !ks.shift && !ks.caps {
			k = k + 32
		}
	}
	_ = ks.client.SendKey(bring.KeyCode(k), pressed)
}

func (ks *keyboardState) keyUp(scancode int, rn rune, name string) {
	if k := mapKey(scancode, rn); k >= 0 {
		if k == bring.KeyLeftShift || k == bring.KeyRightShift {
			ks.shift = false
		}
		if k == bring.KeyCapsLock {
			ks.caps = false
		}
		ks.sendKey(k, false)
	}
}

func (ks *keyboardState) keyDown(scancode int, rn rune, name string) {
	if k := mapKey(scancode, rn); k >= 0 {
		if k == bring.KeyLeftShift || k == bring.KeyRightShift {
			ks.shift = true
		}
		if k == bring.KeyCapsLock {
			ks.caps = true
		}
		ks.sendKey(k, true)
	}
}

func mapKey(scancode int, rn rune) bring.KeyCode {
	if k, ok := keyMap[scancode]; ok {
		return k
	}
	if rn != 0 {
		return bring.KeyCode(rn)
	}
	return -1
}

var keyMap = map[int]bring.KeyCode{
	sdl.SCANCODE_ESCAPE:     bring.KeyEscape,
	sdl.SCANCODE_RETURN:     bring.KeyEnter,
	sdl.SCANCODE_BACKSPACE:  bring.KeyBackspace,
	sdl.SCANCODE_TAB:        bring.KeyTab,
	sdl.SCANCODE_LCTRL:      bring.KeyLeftControl,
	sdl.SCANCODE_LSHIFT:     bring.KeyLeftShift,
	sdl.SCANCODE_RSHIFT:     bring.KeyRightShift,
	sdl.SCANCODE_LALT:       bring.KeyLeftAlt,
	sdl.SCANCODE_CAPSLOCK:   bring.KeyCapsLock,
	sdl.SCANCODE_F1:         bring.KeyF1,
	sdl.SCANCODE_F2:         bring.KeyF2,
	sdl.SCANCODE_F3:         bring.KeyF3,
	sdl.SCANCODE_F4:         bring.KeyF4,
	sdl.SCANCODE_F5:         bring.KeyF5,
	sdl.SCANCODE_F6:         bring.KeyF6,
	sdl.SCANCODE_F7:         bring.KeyF7,
	sdl.SCANCODE_F8:         bring.KeyF8,
	sdl.SCANCODE_F9:         bring.KeyF9,
	sdl.SCANCODE_F10:        bring.KeyF10,
	sdl.SCANCODE_PAUSE:      bring.KeyPause,
	sdl.SCANCODE_SCROLLLOCK: bring.KeyScroll,
	sdl.SCANCODE_F11:        bring.KeyF11,
	sdl.SCANCODE_F12:        bring.KeyF12,
	sdl.SCANCODE_F13:        bring.KeyF13,
	sdl.SCANCODE_F14:        bring.KeyF14,
	sdl.SCANCODE_F15:        bring.KeyF15,
	sdl.SCANCODE_F16:        bring.KeyF16,
	sdl.SCANCODE_F17:        bring.KeyF17,
	sdl.SCANCODE_F18:        bring.KeyF18,
	sdl.SCANCODE_F19:        bring.KeyF19,
	sdl.SCANCODE_RCTRL:      bring.KeyRightControl,
	sdl.SCANCODE_RALT:       bring.KeyRightAlt,
	sdl.SCANCODE_HOME:       bring.KeyHome,
	sdl.SCANCODE_UP:         bring.KeyUp,
	sdl.SCANCODE_PAGEUP:     bring.KeyPageUp,
	sdl.SCANCODE_LEFT:       bring.KeyLeft,
	sdl.SCANCODE_RIGHT:      bring.KeyRight,
	sdl.SCANCODE_END:        bring.KeyEnd,
	sdl.SCANCODE_DOWN:       bring.KeyDown,
	sdl.SCANCODE_INSERT:     bring.KeyInsert,
	sdl.SCANCODE_DELETE:     bring.KeyDelete,
}
