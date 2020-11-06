package mpris

import (
	"log"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

type Player struct {
	Name  string
	Owner string
	conn  *dbus.Conn
	obj   *dbus.Object
}

type PlaybackState string
type LoopState string

const (
	PlaybackPlaying     PlaybackState = "Playing"
	PlaybackPaused                    = "Paused"
	PlaybackStopped                   = "Stopped"
	PlaybackUnsupported               = ""

	LoopNone     LoopState = "None"
	LoopTrack              = "Track"
	LoopPlaylist           = "Playlist"
)

func (p *Player) Introspect() (*introspect.Node, error) {
	return introspect.Call(p.obj)
}

//
// Methods
//

// Methods on app: org.mpris.MediaPlayer2
func (p *Player) Raise() error {
	call := p.obj.Call(appInterface+".Raise", 0)
	return call.Err
}

func (p *Player) Quit() error {
	call := p.obj.Call(appInterface+".Quit", 0)
	return call.Err
}

// Methods on playback: org.mpris.MediaPlayer2.Player
func (p *Player) Play() error {
	call := p.obj.Call(playerInterface+".Play", 0)
	return call.Err
}

func (p *Player) Pause() error {
	call := p.obj.Call(playerInterface+".Pause", 0)
	return call.Err
}

func (p *Player) PlayPause() error {
	call := p.obj.Call(playerInterface+".PlayPause", 0)
	return call.Err
}

func (p *Player) Next() error {
	call := p.obj.Call(playerInterface+".Next", 0)
	return call.Err
}

func (p *Player) Previous() error {
	call := p.obj.Call(playerInterface+".Previous", 0)
	return call.Err
}

func (p *Player) Stop() error {
	call := p.obj.Call(playerInterface+".Stop", 0)
	return call.Err
}

func (p *Player) OpenUri(uri string) error {
	call := p.obj.Call(playerInterface+".OpenUri", 0, uri)
	return call.Err
}

func (p *Player) Seek(delta int64) error {
	call := p.obj.Call(playerInterface+".Seek", 0, delta)
	return call.Err
}

func (p *Player) SetPosition(trackId string, position int64) error {
	path := dbus.ObjectPath(trackId)
	call := p.obj.Call(playerInterface+".SetPosition", 0, &path, position)
	return call.Err
}

//
// Properties
//

// Properties on app: org.mpris.MediaPlayer2
func (p *Player) Identity() string {
	return getString(p.obj, appInterface, "Identity")
}

func (p *Player) DesktopEntry() string {
	return getString(p.obj, appInterface, "DesktopEntry")
}

func (p *Player) CanRaise() bool {
	return getBool(p.obj, appInterface, "CanRaise")
}

func (p *Player) CanQuit() bool {
	return getBool(p.obj, appInterface, "CanQuit")
}

func (p *Player) CanSetFullscreen() bool {
	return getBool(p.obj, appInterface, "CanSetFullscreen")
}

func (p *Player) Fullscreen(value ...bool) bool {
	if len(value) == 1 {
		setProp(p.obj, appInterface, "Fullscreen", value[0])
		return value[0]
	} else {
		return getBool(p.obj, appInterface, "Fullscreen")
	}
}

func (p *Player) HasTrackList() bool {
	return getBool(p.obj, appInterface, "HasTrackList")
}

func (p *Player) SupportedUriSchemes() []string {
	return getStringList(p.obj, appInterface, "SupportedUriSchemes")
}

func (p *Player) SupportedMimeTypes() []string {
	return getStringList(p.obj, appInterface, "SupportedMimeTypes")
}

// Properties on playback: org.mpris.MediaPlayer2.Player
func (p *Player) Shuffle(value ...bool) bool {
	if len(value) == 1 {
		setProp(p.obj, playerInterface, "Shuffle", value[0])
		return value[0]
	} else {
		return getBool(p.obj, playerInterface, "Shuffle")
	}
}

func (p *Player) Rate(value ...float64) float64 {
	if len(value) == 1 {
		setProp(p.obj, playerInterface, "Rate", value[0])
		return value[0]
	} else {
		return getDouble(p.obj, playerInterface, "Rate")
	}
}

func (p *Player) Volume(value ...float64) float64 {
	if len(value) == 1 {
		setProp(p.obj, playerInterface, "Volume", value[0])
		return value[0]
	} else {
		return getDouble(p.obj, playerInterface, "Volume")
	}
}

func (p *Player) MaximumRate() float64 {
	return getDouble(p.obj, playerInterface, "MaximumRate")
}

func (p *Player) MinimumRate() float64 {
	return getDouble(p.obj, playerInterface, "MinimumRate")
}

func (p *Player) Position() int64 {
	return getInt(p.obj, playerInterface, "Position")
}

func (p *Player) LoopStatus() LoopState {
	return LoopState(getString(p.obj, playerInterface, "LoopStatus"))
}

func (p *Player) PlaybackStatus() PlaybackState {
	return PlaybackState(getString(p.obj, playerInterface, "PlaybackStatus"))
}

func (p *Player) RawMetadata() map[string]dbus.Variant {
	result, err := getProp(p.obj, playerInterface, "Metadata")
	if err != nil {
		return make(map[string]dbus.Variant)
	} else {
		return result.Value().(map[string]dbus.Variant)
	}
}

func (p *Player) CanGoNext() bool {
	return getBool(p.obj, playerInterface, "CanGoNext")
}

func (p *Player) CanGoPrevious() bool {
	return getBool(p.obj, playerInterface, "CanGoPrevious")
}

func (p *Player) CanPlay() bool {
	return getBool(p.obj, playerInterface, "CanPlay")
}

func (p *Player) CanPause() bool {
	return getBool(p.obj, playerInterface, "CanPause")
}

func (p *Player) CanSeek() bool {
	return getBool(p.obj, playerInterface, "CanSeek")
}

func (p *Player) CanControl() bool {
	return getBool(p.obj, playerInterface, "CanControl")
}

//
// Signals
//
func (p *Player) OnSeek() (chan *dbus.Signal, error) {
	c := make(chan *dbus.Signal, 10)
	err := p.conn.AddMatchSignal(
		dbus.WithMatchInterface(playerInterface),
	)
	if err != nil {
		return c, err
	}
	p.conn.Signal(c)
	return c, nil
}

func (p *Player) OnChange() (chan *dbus.Signal, error) {
	log.Println("OnChange", p.obj.Destination())
	c := make(chan *dbus.Signal, 10)
	err := p.conn.AddMatchSignal(
		// dbus.WithMatchSender(p.app.Owner),
		dbus.WithMatchObjectPath(objectPath),
		dbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
	)
	if err != nil {
		return c, err
	}
	p.conn.Signal(c)
	return c, nil
}
