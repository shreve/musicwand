package main

import (
	"github.com/godbus/dbus/v5"
	"github.com/shreve/musicwand/pkg/mpris"
	"log"
)

//
// App Server
//
// Controls and information about the media player application.
//
type appServer struct {
	client *mpris.App
}

func (a *appServer) Quit() *dbus.Error {
	a.client.Quit()
	return nil
}

func (a *appServer) Raise() *dbus.Error {
	a.client.Raise()
	return nil
}

//
// Player Server
//
// Controls for the playback of media.
//
type playerServer struct {
	client *mpris.Player
}

func (p playerServer) Next() *dbus.Error {
	p.client.Next()
	return nil
}

func (p playerServer) OpenUri(uri string) *dbus.Error {
	p.client.OpenUri(uri)
	return nil
}

func (p playerServer) Pause() *dbus.Error {
	p.client.Pause()
	return nil
}

func (p playerServer) Play() *dbus.Error {
	p.client.Play()
	return nil
}

func (p playerServer) PlayPause() *dbus.Error {
	p.client.PlayPause()
	return nil
}

func (p playerServer) Previous() *dbus.Error {
	p.client.Previous()
	return nil
}

func (p playerServer) Seek(delta int64) *dbus.Error {
	p.client.Seek(delta)
	return nil
}

func (p playerServer) SetPosition(trackId string, position int64) *dbus.Error {
	p.client.SetPosition(trackId, position)
	return nil
}

func (p playerServer) Stop() *dbus.Error {
	p.client.Stop()
	return nil
}

//
// Property Server
//
// Forwards properties from destination object.
//
type propertyHandler struct {
	client *mpris.App
}

func (p propertyHandler) Get(iface, prop string) (dbus.Variant, *dbus.Error) {
	result, _ := p.client.Properties().Get(iface, prop)
	return result, nil
}

func (p propertyHandler) GetAll(iface string) (map[string]dbus.Variant, *dbus.Error) {
	// Prevent recursion
	if iface == "org.freedesktop.DBus.Properties" || iface == "com.github.shreve.musicwand" {
		return nil, nil
	}
	result, _ := p.client.Properties().GetAll(iface)
	return result, nil
}

func (p propertyHandler) Set(iface, prop string, value dbus.Variant) *dbus.Error {
	p.client.Properties().Set(iface, prop, value)
	return nil
}

type State struct {
	client     mpris.Client
	CurrentApp mpris.App
}

func (s *State) SetCurrentApp(name string) *dbus.Error {
	newApp := *s.client.FindApp(name)
	if newApp.Identity() != "" {
		s.CurrentApp = newApp
	} else {
	}
	return nil
}

func (s *State) selectApp() {
	apps := s.client.Apps()
	if len(apps) == 0 {
		log.Fatal("Unable to connect to any music players")
	}

	s.CurrentApp = apps[0]
	for _, app := range apps {
		player := app.Player()
		if player.PlaybackStatus() == mpris.PlaybackPlaying {
			s.CurrentApp = app
			break
		}
	}
}

//
// Daemon
//
// The process to perform the work of music wand
//
func RunDaemon() {
	server, err := mpris.NewServer("musicwand")
	if err != nil {
		log.Fatal(err)
	}

	client, err := mpris.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	state := State{client: *client}
	state.selectApp()

	app := &state.CurrentApp
	player := app.Player()

	server.PropertyHandler = &propertyHandler{app}
	server.AppServer = &appServer{app}
	server.PlayerServer = &playerServer{&player}

	server.AddInterface("com.github.shreve.musicwand", &state)

	if err := server.Listen(); err != nil {
		log.Fatal(err)
	}
}
