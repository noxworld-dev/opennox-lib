package sdl

import (
	"image"
	"strings"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/noxworld-dev/opennox-lib/client/seat"
	"github.com/noxworld-dev/opennox-lib/types"
)

func (win *Window) ReplaceInputs(cfg seat.InputConfig) seat.InputConfig {
	oldCfg := win.onInput
	win.onInput = cfg
	return oldCfg
}

func (win *Window) OnInput(fnc func(ev seat.InputEvent)) {
	win.onInput = append(win.onInput, fnc)
}

func (win *Window) SetTextInput(enable bool) {
	if win.textInp == enable {
		return
	}
	win.textInp = enable
	if enable {
		sdl.StartTextInput()
	} else {
		sdl.StopTextInput()
	}
}

func (win *Window) InputTick() {
	for {
		switch ev := sdl.PollEvent().(type) {
		case nil:
			// no more events
			return
		case sdl.TextEditingEvent:
			win.processTextEditingEvent(&ev)
		case sdl.TextInputEvent:
			win.processTextInputEvent(&ev)
		case sdl.KeyboardEvent:
			win.processKeyboardEvent(&ev)
		case sdl.MouseButtonEvent:
			win.processMouseButtonEvent(&ev)
		case sdl.MouseMotionEvent:
			win.processMotionEvent(&ev)
		case sdl.MouseWheelEvent:
			win.processWheelEvent(&ev)
		case sdl.ControllerAxisEvent:
			if debugGpad {
				Log.Printf("SDL event: SDL_CONTROLLERAXISMOTION (%x): joy=%d, axis=%d, val=%d\n",
					ev.GetType(), ev.Which, ev.Axis, ev.Value)
			}
			win.processGamepadAxisEvent(&ev)
		case sdl.ControllerButtonEvent:
			if debugGpad {
				Log.Printf("SDL event: SDL_CONTROLLERBUTTON (%x): joy=%d, btn=%d, state=%d\n",
					ev.GetType(), ev.Which, ev.Button, ev.State)
			}
			win.processGamepadButtonEvent(&ev)
		case *sdl.ControllerDeviceEvent:
			switch ev.GetType() {
			case sdl.CONTROLLERDEVICEADDED:
				if debugGpad {
					Log.Printf("SDL event: SDL_CONTROLLERDEVICEADDED (%x): joy=%d\n", ev.GetType(), ev.Which)
				}
				win.processGamepadDeviceEvent(ev)
			case sdl.CONTROLLERDEVICEREMOVED:
				if debugGpad {
					Log.Printf("SDL event: SDL_CONTROLLERDEVICEREMOVED (%x): joy=%d\n", ev.GetType(), ev.Which)
				}
				win.processGamepadDeviceEvent(ev)
			case sdl.CONTROLLERDEVICEREMAPPED:
				if debugGpad {
					Log.Printf("SDL event: SDL_CONTROLLERDEVICEREMAPPED (%x)\n", ev.GetType())
				}
			}
		case sdl.WindowEvent:
			win.processWindowEvent(&ev)
		case sdl.QuitEvent:
			win.processQuitEvent(&ev)
		}
		// TODO: touch events for WASM
	}
}

func (win *Window) inputEvent(ev seat.InputEvent) {
	for _, fnc := range win.onInput {
		fnc(ev)
	}
}

func (win *Window) processQuitEvent(ev *sdl.QuitEvent) {
	win.inputEvent(seat.WindowClosed)
}

func (win *Window) processWindowEvent(ev *sdl.WindowEvent) {
	switch ev.Event {
	case sdl.WINDOWEVENT_FOCUS_LOST:
		win.inputEvent(seat.WindowUnfocused)
	case sdl.WINDOWEVENT_FOCUS_GAINED:
		win.inputEvent(seat.WindowFocused)
	}
}

func (win *Window) processTextEditingEvent(ev *sdl.TextEditingEvent) {
	win.inputEvent(&seat.TextEditEvent{
		Text: ev.GetText(),
	})
}

func (win *Window) processTextInputEvent(ev *sdl.TextInputEvent) {
	text := ev.GetText()
	if sdl.GetModState()&sdl.KMOD_CTRL != 0 && len(text) == 1 && strings.ToLower(text) == "v" {
		return // ignore "V" from Ctrl-V
	}
	win.inputEvent(&seat.TextInputEvent{
		Text: text,
	})
}

func (win *Window) processKeyboardEvent(ev *sdl.KeyboardEvent) {
	if win.textInp && ev.State == sdl.PRESSED && sdl.GetModState()&sdl.KMOD_CTRL != 0 && ev.Keysym.Scancode == sdl.SCANCODE_V {
		text, err := sdl.GetClipboardText()
		if err != nil {
			Log.Printf("cannot get clipboard text: %v", err)
			return
		}
		win.inputEvent(&seat.TextInputEvent{
			Text: text,
		})
		return
	}
	key := scanCodeToKeyNum[ev.Keysym.Scancode]
	win.inputEvent(&seat.KeyboardEvent{
		Key:     key,
		Pressed: ev.State == sdl.PRESSED,
	})
}

func (win *Window) processMouseButtonEvent(ev *sdl.MouseButtonEvent) {
	pressed := ev.State == sdl.PRESSED
	// TODO: handle focus, or move to other place
	//if pressed {
	//	h.iface.WindowEvent(WindowFocus)
	//}

	var button seat.MouseButton
	switch ev.Button {
	case sdl.BUTTON_LEFT:
		button = seat.MouseButtonLeft
	case sdl.BUTTON_RIGHT:
		button = seat.MouseButtonRight
	case sdl.BUTTON_MIDDLE:
		button = seat.MouseButtonMiddle
	default:
		return
	}
	win.inputEvent(&seat.MouseButtonEvent{
		Button:  button,
		Pressed: pressed,
	})
}

func (win *Window) processMotionEvent(ev *sdl.MouseMotionEvent) {
	win.inputEvent(&seat.MouseMoveEvent{
		Relative: win.rel,
		Pos:      image.Point{X: int(ev.X), Y: int(ev.Y)},
		Rel:      types.Pointf{X: float32(ev.XRel), Y: float32(ev.YRel)},
	})
}

func (win *Window) processWheelEvent(ev *sdl.MouseWheelEvent) {
	win.inputEvent(&seat.MouseWheelEvent{
		Wheel: int(ev.Y),
	})
}

func (win *Window) processGamepadButtonEvent(ev *sdl.ControllerButtonEvent) {
	// TODO: handle gamepads (again)
}

func (win *Window) processGamepadAxisEvent(ev *sdl.ControllerAxisEvent) {
	// TODO: handle gamepads (again)
}

func (win *Window) processGamepadDeviceEvent(ev *sdl.ControllerDeviceEvent) {
	// TODO: handle gamepads (again)
}
