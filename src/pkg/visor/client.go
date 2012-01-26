package visor

import (
	"fmt"
	"github.com/soundcloud/doozer"
	"net"
)

type Client struct {
	Addr *net.TCPAddr
	conn *doozer.Conn
	Root string
	rev  int64
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Del(path string) (err error) {
	rev, err := c.conn.Rev()
	if err != nil {
		return
	}

	c.rev = rev

	err = doozer.Walk(c.conn, rev, path, func(path string, f *doozer.FileInfo, e error) error {
		if e != nil {
			return e
		}

		if !f.IsDir {
			e = c.conn.Del(path, rev)
			if e != nil {
				return e
			}
		}

		return nil
	})

	return
}

func (c *Client) Exists(path string) (exists bool, err error) {
	_, rev, err := c.conn.Stat(path, nil)
	if err != nil {
		return
	}

	switch rev {
	case 0:
		exists = false
	default:
		exists = true
	}

	return exists, nil
}

func (c *Client) Get(path string) (value string, err error) {
	rev, err := c.conn.Rev()
	if err != nil {
		return
	}

	body, rev, err := c.conn.Get(path, &rev)
	if err != nil {
		return
	}
	if rev == 0 {
		err = ErrKeyNotFound

		return
	}

	c.rev = rev

	value = string(body)

	return
}

func (c *Client) Keys(path string) (keys []string, err error) {
	rev, err := c.conn.Rev()
	if err != nil {
		return
	}

	c.rev = rev

	keys, err = c.conn.Getdir(path, c.rev, 0, -1)
	if err != nil {
		return
	}

	return
}

func (c *Client) Set(path string, body string) (err error) {
	rev, err := c.conn.Set(path, c.rev, []byte(body))
	if err != nil {
		return
	}

	c.rev = rev

	return
}

func (c *Client) String() string {
	return fmt.Sprintf("%#v", c)
}

// INSTANCES

// TICKETS

func (c *Client) Tickets() ([]Ticket, error) {
	return nil, nil
}
func (c *Client) HostTickets(addr string) ([]Ticket, error) {
	return nil, nil
}

// EVENTS

func (c *Client) WatchEvent(listener chan *Event) error {
	rev, _ := c.conn.Rev()

	for {
		ev, _ := c.conn.Wait(c.Root+"*", rev)
		event := &Event{EV_APP_REG, string(ev.Body), &ev}
		rev = ev.Rev + 1
		listener <- event
	}
	return nil
}
func (c *Client) WatchTicket(listener chan *Ticket) error {
	return nil
}
