package mpris

import (
	"log"

	"github.com/godbus/dbus/v5"
)

const (
	getPropertyMethod    = "org.freedesktop.DBus.Properties.Get"
	getAllPropertyMethod = "org.freedesktop.DBus.Properties.GetAll"
	setPropertyMethod    = "org.freedesktop.DBus.Properties.Set"
	introspectMethod     = "org.freedesktop.DBus.Introspectable.Introspect"
)

func (p *Player) Get(iface, prop string) (result dbus.Variant, err error) {
	err = p.obj.Call(getPropertyMethod, 0, iface, prop).Store(&result)
	if err != nil {
		log.Println("Error getting property: ", err)
	}
	return
}

func (p *Player) GetAll(iface string) (result map[string]dbus.Variant, err error) {
	err = p.obj.Call(getAllPropertyMethod, 0, iface).Store(&result)
	if err != nil {
		log.Println("Error getting all properties: ", err, p.obj)
	}
	return
}

func (p *Player) Set(iface, prop string, value interface{}) (err error) {
	call := p.obj.Call(setPropertyMethod, 0, iface, prop, value)
	err = call.Err
	if err != nil {
		log.Println("Error setting property: ", err)
	}
	return
}
