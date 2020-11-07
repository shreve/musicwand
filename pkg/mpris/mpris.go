// Package mpris provides a client and server for the MPRIS D-Bus interface.
package mpris

import (
	"log"

	"github.com/godbus/dbus/v5"
)

const (
	objectPath      = "/org/mpris/MediaPlayer2"
	appInterface    = "org.mpris.MediaPlayer2"
	playerInterface = "org.mpris.MediaPlayer2.Player"
)

// Convert an error into a dbus failed error if the error exists.
func DbusError(err error) *dbus.Error {
	if err != nil {
		return dbus.MakeFailedError(err)
	} else {
		return nil
	}
}

func getProp(obj *dbus.Object, iface, prop string) (result dbus.Variant, err error) {
	err = obj.Call(getPropertyMethod, 0, iface, prop).Store(&result)
	if err != nil {
		log.Println("Error getting property: ", err)
	}
	return
}

func getBool(obj *dbus.Object, iface, prop string) bool {
	result, err := getProp(obj, iface, prop)
	if err != nil {
		return false
	} else {
		return result.Value().(bool)
	}
}

func getDouble(obj *dbus.Object, iface, prop string) float64 {
	result, err := getProp(obj, iface, prop)
	if err != nil {
		return 0.0
	} else {
		return result.Value().(float64)
	}
}

func getInt(obj *dbus.Object, iface, prop string) int64 {
	result, err := getProp(obj, iface, prop)
	if err != nil {
		return 0
	} else {
		return result.Value().(int64)
	}
}

func getString(obj *dbus.Object, iface, prop string) string {
	result, err := getProp(obj, iface, prop)
	if err != nil {
		return ""
	} else {
		return result.Value().(string)
	}
}

func getStringList(obj *dbus.Object, iface, prop string) []string {
	result, err := getProp(obj, iface, prop)
	if err != nil {
		return []string{}
	} else {
		return result.Value().([]string)
	}
}

func setProp(obj *dbus.Object, iface, prop string, value interface{}) error {
	call := obj.Call(setPropertyMethod, 0, iface, prop, value)
	return call.Err
}
