package mpris

import (
	"strings"

	"github.com/godbus/dbus/v5"
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

func (c *Client) Players() (players []Player) {
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

			players = append(players, Player{
				conn:  c.conn,
				obj:   object,
				Name:  name,
				Owner: owner,
			})
		}
	}
	return
}

func (c *Client) FindPlayer(name string) *Player {
	for _, player := range c.Players() {
		if strings.HasSuffix(player.Name, name) {
			return &player
		}
	}
	return nil
}

func (c *Client) PlayerWithOwner(owner string) *Player {
	for _, player := range c.Players() {
		if player.Owner == owner {
			return &player
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
