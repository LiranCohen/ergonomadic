package irc

import (
	"net"
)

type Client struct {
	conn       net.Conn
	send       chan<- string
	recv       <-chan string
	username   string
	realname   string
	nick       string
	registered bool
	invisible  bool
}

func NewClient(conn net.Conn) *Client {
	client := new(Client)
	client.conn = conn
	client.send = StringWriteChan(conn)
	client.recv = StringReadChan(conn)
	return client
}

// Adapt `chan string` to a `chan Message`.
func (c *Client) Communicate(server chan<- *ClientMessage) {
	for str := range c.recv {
		m := ParseMessage(str)
		if m != nil {
			server <- &ClientMessage{c, m}
		}
	}
}

func (c *Client) Nick() string {
	if c.nick != "" {
		return c.nick
	}
	return "<guest>"
}

func (c *Client) UModeString() string {
	if c.invisible {
		return "+i"
	}
	return ""
}