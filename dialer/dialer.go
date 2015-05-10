package dialer

import (
	"log"
	"net"

	"github.com/linkosmosis/godns"
)

// Dialer -
type Dialer struct {
	*net.Dialer

	AddrsPool *godns.Pool
}

// New - initalize dial.Dialer wrapper
func New() *Dialer {
	return &Dialer{
		Dialer:    &net.Dialer{},
		AddrsPool: godns.New(),
	}
}

// Dial - lightweight version of dialer.Dial, this has cached
// dns hostport and shorter TCP connection setup
func (d *Dialer) Dial(network, address string) (net.Conn, error) {
	tcpAddr, err := d.AddrsPool.Get(address)
	if err != nil {
		log.Println("Failed to resolve:", err)
		return d.Dialer.Dial(network, address)
	}
	c, err := net.DialTCP(network, nil, tcpAddr)
	if err != nil {
		log.Println("Failed to setup DialTCP:", err)
		return d.Dialer.Dial(network, address)
	}
	if d.KeepAlive != 0 {
		c.SetKeepAlive(true)
		c.SetKeepAlivePeriod(d.KeepAlive)
		c.SetLinger(0)
		c.SetNoDelay(true)
	}
	return c, err
}
