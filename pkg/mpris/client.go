package mpris

import (
	"github.com/godbus/dbus/v5"
	"strings"
)

type Client struct {
	conn *dbus.Conn
}

func NewClient() (*Client, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	return &Client{conn}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Apps() (apps []App) {
	var list []string
	err := c.conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&list)
	if err != nil {
		return
	}
	for i := range list {
		if strings.HasPrefix(list[i], appInterface) {
			object := c.conn.Object(list[i], objectPath).(*dbus.Object)
			apps = append(apps, App{c.conn, object})
		}
	}
	return
}

func (c *Client) FindApp(name string) *App {
	object := c.conn.Object(appInterface+"."+name, objectPath).(*dbus.Object)
	return &App{c.conn, object}
}
