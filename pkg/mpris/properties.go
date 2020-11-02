package mpris

import (
	"github.com/godbus/dbus/v5"
	"log"
)

const (
	getPropertyMethod    = "org.freedesktop.DBus.Properties.Get"
	getAllPropertyMethod = "org.freedesktop.DBus.Properties.GetAll"
	setPropertyMethod    = "org.freedesktop.DBus.Properties.Set"
	introspectMethod     = "org.freedesktop.DBus.Introspectable.Introspect"
)

type Properties struct {
	obj *dbus.Object
}

func (p Properties) Get(iface, prop string) (result dbus.Variant, err error) {
	err = p.obj.Call(getPropertyMethod, 0, iface, prop).Store(&result)
	if err != nil {
		log.Println("Error getting property: ", err)
	}
	return
}

func (p Properties) GetAll(iface string) (result map[string]dbus.Variant, err error) {
	err = p.obj.Call(getAllPropertyMethod, 0, iface).Store(&result)
	if err != nil {
		log.Println("Error getting all properties: ", err, p.obj)
	}
	return
}

func (p Properties) Set(iface, prop string, value interface{}) (err error) {
	call := p.obj.Call(setPropertyMethod, 0, iface, prop, value)
	err = call.Err
	if err != nil {
		log.Println("Error setting property: ", err)
	}
	return
}
