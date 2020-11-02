package mpris

import (
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

type App struct {
	conn  *dbus.Conn
	obj   *dbus.Object
	Name  string
	Owner string
}

func (a *App) Introspect() (*introspect.Node, error) {
	return introspect.Call(a.obj)
}

func (a *App) Player() Player {
	return Player{a.conn, a.obj, a}
}

func (a *App) Properties() Properties {
	return Properties{a.obj}
}

func (a *App) Identity() string {
	return getString(a.obj, appInterface, "Identity")
}

func (a *App) DesktopEntry() string {
	return getString(a.obj, appInterface, "DesktopEntry")
}

func (a *App) CanRaise() bool {
	return getBool(a.obj, appInterface, "CanRaise")
}

func (a *App) Raise() {
	a.obj.Call(appInterface+".Raise", 0)
}

func (a *App) CanQuit() bool {
	return getBool(a.obj, appInterface, "CanQuit")
}

func (a *App) Quit() {
	a.obj.Call(appInterface+".Quit", 0)
}

func (a *App) CanSetFullscreen() bool {
	return getBool(a.obj, appInterface, "CanSetFullscreen")
}

func (a *App) Fullscreen(value ...bool) bool {
	if len(value) == 1 {
		setProp(a.obj, appInterface, "Fullscreen", value[0])
		return value[0]
	} else {
		return getBool(a.obj, appInterface, "Fullscreen")
	}
}

func (a *App) HasTrackList() bool {
	return getBool(a.obj, appInterface, "HasTrackList")
}

func (a *App) SupportedUriSchemes() []string {
	return getStringList(a.obj, appInterface, "SupportedUriSchemes")
}

func (a *App) SupportedMimeTypes() []string {
	return getStringList(a.obj, appInterface, "SupportedMimeTypes")
}
