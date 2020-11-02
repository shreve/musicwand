package mpris

import (
	"github.com/godbus/dbus/v5"
	"log"
)

const (
	objectPath = "/org/mpris/MediaPlayer2"

	appInterface    = "org.mpris.MediaPlayer2"
	playerInterface = "org.mpris.MediaPlayer2.Player"
)

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
