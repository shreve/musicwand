package main

import (
	"log"

	"github.com/godbus/dbus/v5"
	"github.com/shreve/musicwand/pkg/mpris"
)

//
// App Server
//
// Controls and information about the media player application.
//
type appServer struct {
	client *mpris.Player
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
	client *mpris.Player
}

func (p propertyHandler) Get(iface, prop string) (dbus.Variant, *dbus.Error) {
	result, _ := p.client.Get(iface, prop)
	return result, nil
}

func (p propertyHandler) GetAll(iface string) (map[string]dbus.Variant, *dbus.Error) {
	// Prevent recursion
	if iface == "org.freedesktop.DBus.Properties" || iface == "com.github.shreve.musicwand" {
		return nil, nil
	}
	result, _ := p.client.GetAll(iface)
	return result, nil
}

func (p propertyHandler) Set(iface, prop string, value dbus.Variant) *dbus.Error {
	p.client.Set(iface, prop, value)
	return nil
}

type State struct {
	client        mpris.Client
	CurrentPlayer mpris.Player
}

func (s *State) SetCurrentPlayer(name string) *dbus.Error {
	newPlayer := s.client.FindPlayer(name)
	if newPlayer != nil {
		s.CurrentPlayer = *newPlayer
	} else {
		return dbus.NewError("Unable to find that app", []interface{}{})
	}
	return nil
}

func (s *State) setPlayer(player mpris.Player) {
	log.Println("Setting current app to:", player.Name)
	s.CurrentPlayer = player
}

func (s *State) selectPlayer() {
	players := s.client.Players()
	if len(players) == 0 {
		log.Fatal("Unable to connect to any music players")
	}

	for _, player := range players {
		if player.PlaybackStatus() == mpris.PlaybackPlaying {
			s.setPlayer(player)
			return
		}
	}
	s.setPlayer(players[0])
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
	state.selectPlayer()

	player := &state.CurrentPlayer

	server.PropertyHandler = &propertyHandler{player}
	server.AppServer = &appServer{player}
	server.PlayerServer = &playerServer{player}

	server.AddInterface("com.github.shreve.musicwand", &state)

	go func() {
		events, _ := client.OnAnyPlayerChange()
		for {
			event := <-events
			player := client.PlayerWithOwner(event.Sender)
			if player != nil {
				state.setPlayer(*player)
			}
		}
	}()

	if err := server.Listen(); err != nil {
		log.Fatal(err)
	}
}
