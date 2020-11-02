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
	for _, name := range list {
		if strings.HasPrefix(name, appInterface) {
			object := c.conn.Object(name, objectPath).(*dbus.Object)
			var owner string
			c.conn.BusObject().Call("org.freedesktop.DBus.GetNameOwner", 0, name).Store(&owner)
			apps = append(apps, App{c.conn, object, name, owner})
		}
	}
	return
}

func (c *Client) FindApp(name string) *App {
	for _, app := range c.Apps() {
		if strings.HasSuffix(app.Name, name) {
			return &app
		}
	}
	return nil
}

func (c *Client) AppWithOwner(owner string) *App {
	for _, app := range c.Apps() {
		if app.Owner == owner {
			return &app
		}
	}
	return nil
}

func (c *Client) OnAnyPlayerChange() (chan *dbus.Signal, error) {
	signals := make(chan *dbus.Signal, 10)
	err := c.conn.AddMatchSignal(
		dbus.WithMatchObjectPath(objectPath),
		dbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
	)
	if err != nil {
		return signals, err
	}
	c.conn.Signal(signals)
	return signals, nil
}
