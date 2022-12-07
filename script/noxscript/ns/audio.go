package ns

import (
	"github.com/noxworld-dev/opennox-lib/script"
	"github.com/noxworld-dev/opennox-lib/script/noxscript/ns/audio"
)

// AudioEvent play an audio event at a location.
func AudioEvent(audio audio.Name, p script.Positioner) {
	if impl == nil {
		return
	}
	impl.AudioEvent(audio, p)
}

// Music plays music.
func Music(music int, volume int) {
	if impl == nil {
		return
	}
	impl.Music(music, volume)
}
func MusicPushEvent() {
	if impl == nil {
		return
	}
	impl.MusicPushEvent()
}
func MusicPopEvent() {
	if impl == nil {
		return
	}
	impl.MusicPopEvent()
}
func MusicEvent() {
	if impl == nil {
		return
	}
	impl.MusicEvent()
}
