package mpris

import (
	"encoding/xml"
	"errors"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"log"
	"strings"
)

const (
	introspectableInterface = "org.freedesktop.DBus.Introspectable"
	propertyInterface       = "org.freedesktop.DBus.Properties"
	peerInterface           = "org.freedesktop.DBus.Peer"

	mprisIntrospectXML = `<!DOCTYPE node PUBLIC "-//freedesktop//DTD D-BUS Object Introspection 1.0//EN" "http://www.freedesktop.org/standards/dbus/1.0/introspect.dtd">
<node name="/org/mpris/MediaPlayer2">
	<interface name="org.freedesktop.DBus.Peer">
		<method name="Ping"></method>
		<method name="GetMachineId">
			<arg name="machine_uuid" type="s" direction="out"></arg>
		</method>
	</interface>
	<interface name="org.freedesktop.DBus.Introspectable">
		<method name="Introspect">
			<arg name="data" type="s" direction="out"></arg>
		</method>
	</interface>
	<interface name="org.freedesktop.DBus.Properties">
		<method name="Get">
			<arg name="interface" type="s" direction="in"></arg>
			<arg name="property" type="s" direction="in"></arg>
			<arg name="value" type="v" direction="out"></arg>
		</method>
		<method name="GetAll">
			<arg name="interface" type="s" direction="in"></arg>
			<arg name="properties" type="a{sv}" direction="out"></arg>
		</method>
		<method name="Set">
			<arg name="interface" type="s" direction="in"></arg>
			<arg name="property" type="s" direction="in"></arg>
			<arg name="value" type="v" direction="in"></arg>
		</method>
		<signal name="PropertiesChanged">
			<arg name="interface" type="s"></arg>
			<arg name="changed_properties" type="a{sv}"></arg>
			<arg name="invalidated_properties" type="as"></arg>
		</signal>
	</interface>
	<interface name="org.mpris.MediaPlayer2.Player">
		<method name="Next"></method>
		<method name="Previous"></method>
		<method name="Pause"></method>
		<method name="PlayPause"></method>
		<method name="Stop"></method>
		<method name="Play"></method>
		<method name="Seek">
			<arg type="x" direction="in"></arg>
		</method>
		<method name="SetPosition">
			<arg type="o" direction="in"></arg>
			<arg type="x" direction="in"></arg>
		</method>
		<method name="OpenUri">
			<arg type="s" direction="in"></arg>
		</method>
		<signal name="Seeked">
			<arg type="x"></arg>
		</signal>
		<property name="PlaybackStatus" type="s" access="read"></property>
		<property name="LoopStatus" type="s" access="readwrite"></property>
		<property name="Rate" type="d" access="readwrite"></property>
		<property name="Shuffle" type="b" access="readwrite"></property>
		<property name="Volume" type="d" access="readwrite"></property>
		<property name="Position" type="x" access="read">
			<annotation name="org.freedesktop.DBus.Property.EmitsChangedSignal" value="false"></annotation>
		</property>
		<property name="MinimumRate" type="d" access="read"></property>
		<property name="MaximumRate" type="d" access="read"></property>
		<property name="CanGoNext" type="b" access="read"></property>
		<property name="CanGoPrevious" type="b" access="read"></property>
		<property name="CanPlay" type="b" access="read"></property>
		<property name="CanPause" type="b" access="read"></property>
		<property name="CanSeek" type="b" access="read"></property>
		<property name="CanControl" type="b" access="read">
			<annotation name="org.freedesktop.DBus.Property.EmitsChangedSignal" value="false"></annotation>
		</property>
		<property name="Metadata" type="a{sv}" access="read"></property>
	</interface>
	<interface name="org.mpris.MediaPlayer2">
		<method name="Raise"></method>
		<method name="Quit"></method>
		<property name="CanQuit" type="b" access="read"></property>
		<property name="Fullscreen" type="b" access="readwrite"></property>
		<property name="CanSetFullscreen" type="b" access="read"></property>
		<property name="CanRaise" type="b" access="read"></property>
		<property name="HasTrackList" type="b" access="read"></property>
		<property name="Identity" type="s" access="read"></property>
		<property name="SupportedUriSchemes" type="as" access="read"></property>
		<property name="SupportedMimeTypes" type="as" access="read"></property>
	</interface>
</node>
`
)

var mprisIntrospect = func() (node introspect.Node) {
	xml.NewDecoder(strings.NewReader(mprisIntrospectXML)).Decode(&node)
	return
}()

////
//// Peer Server
////
//// Information about the machine running this DBus service.
//// This seems to be implemented automatically.
////
//type peerServer struct{}

//func (p peerServer) GetMachineId() *dbus.Error {
//	return nil
//}

//func (p peerServer) Ping() *dbus.Error {
//	return nil
//}

//
// MPRIS Server
//
// Create and listen to start application.
//
type Server struct {
	Conn            *dbus.Conn
	Name            string
	AppServer       IsApp
	PlayerServer    IsPlayer
	PropertyHandler HandlesProperties

	def    introspect.Node
	custom map[string]interface{}
}

type HandlesProperties interface {
	Get(iface, prop string) (dbus.Variant, *dbus.Error)
	GetAll(iface string) (map[string]dbus.Variant, *dbus.Error)
	Set(iface, prop string, value dbus.Variant) *dbus.Error
}

type IsApp interface {
	Quit() *dbus.Error
	Raise() *dbus.Error
}

type IsPlayer interface {
	Next() *dbus.Error
	OpenUri(uri string) *dbus.Error
	Pause() *dbus.Error
	Play() *dbus.Error
	PlayPause() *dbus.Error
	Previous() *dbus.Error
	Seek(delta int64) *dbus.Error
	SetPosition(trackId string, position int64) *dbus.Error
	Stop() *dbus.Error
}

func NewServer(name string) (*Server, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	server := Server{def: mprisIntrospect, Conn: conn, Name: name}
	server.custom = make(map[string]interface{})
	return &server, nil
}

func (s *Server) Close() {
	s.Conn.ReleaseName(appInterface + "." + s.Name)
	s.Conn.Close()
}

func (s *Server) AddInterface(name string, handler interface{}) {
	s.custom[name] = handler
	customInt := introspect.Interface{
		Name:    name,
		Methods: introspect.Methods(handler),
	}
	s.def.Interfaces = append(s.def.Interfaces, customInt)
}

func (s *Server) Listen() error {
	// First, publish the introspection of the whole server. This is static.
	s.Conn.Export(
		introspect.NewIntrospectable(&s.def),
		objectPath,
		introspectableInterface)

	s.Conn.Export(s.PropertyHandler, objectPath, propertyInterface)
	s.Conn.Export(s.AppServer, objectPath, appInterface)
	s.Conn.Export(s.PlayerServer, objectPath, playerInterface)

	for name, handler := range s.custom {
		s.Conn.Export(handler, objectPath, name)
	}

	serverName := appInterface + "." + s.Name

	// Now let's name our server
	reply, err := s.Conn.RequestName(serverName, dbus.NameFlagReplaceExisting)
	if err != nil || reply != dbus.RequestNameReplyPrimaryOwner {
		return errors.New("Unable to claim " + serverName)
	}

	log.Println("Started DBus server on " + serverName)

	select {}
}
