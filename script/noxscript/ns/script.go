// Package ns exposes NoxScript API used in maps.
//
// This package is for documentation and type-checking only, all functions do nothing.
//
// Names of functions are synchronized with the latest NoxScript implementation:
// https://noxtools.github.io/noxscript/.
package ns

type Func = any

type Handle interface {
	// ScriptID returns internal script ID for the object.
	ScriptID() int
}

type Timer int

// SecondTimer creates a timer that calls the given script function after a delay given in seconds.
func SecondTimer(sec int, fnc Func) Timer {
	// header only
	return 0
}

// FrameTimer creates a timer that calls the given script function after a delay given in frames.
func FrameTimer(frames int, fnc Func) Timer {
	// header only
	return 0
}

// SecondTimerWithArg creates a timer that calls the given script function after a delay given in seconds.
// The given argument will be passed into the script function.
func SecondTimerWithArg(seconds int, arg any, fnc Func) Timer {
	// header only
	return 0
}

// FrameTimerWithArg creates a timer that calls the given script function after a delay given in frames.
// The given argument will be passed into the script function.
func FrameTimerWithArg(frames int, arg any, fnc Func) Timer {
	// header only
	return 0
}

// CancelTimer cancels a timer. Returns true if successful.
func CancelTimer(id Timer) bool {
	// header only
	return false
}

// RandomFloat generates random float.
func RandomFloat(min float32, max float32) float32 {
	// header only
	return 0
}

// Random generates random int.
func Random(min int, max int) int {
	// header only
	return 0
}

// IntToString converts an int to a string.
func IntToString(val int) string {
	// header only
	return ""
}

// FloatToString converts a float to a string.
func FloatToString(val float32) string {
	// header only
	return ""
}

// Distance between two locations.
func Distance(x1 float32, y1 float32, x2 float32, y2 float32) float32 {
	// header only
	return 0
}

// StopScript exits current script function.
func StopScript(value any) {
	// header only
}

// AutoSave triggers an autosave. Only solo games.
func AutoSave() {
	// header only
}

// StartupScreen shows startup screen to host.
func StartupScreen(which int) {
	// header only
}

// DeathScreen shows death screen to host.
func DeathScreen(which int) {
	// header only
}
