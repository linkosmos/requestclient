package godns

import (
	"errors"
	"math/rand"
	"net"
)

// ErrEmptyIPS - no IP's were found for given hostport
var ErrEmptyIPS = errors.New("No IP's for a given hostport")

// TCPMap - map of host+port and resolved IP's, added some syntactic
// sugar and convenience methods
type TCPMap map[string][]*net.TCPAddr

// Add -
func (t *TCPMap) Add(hostport string, ip net.IP, port int) {
	if t.Exist(hostport) {
		t.append(hostport, ip, port)
	} else {
		var addresses []*net.TCPAddr
		(*t)[hostport] = addresses
		t.append(hostport, ip, port)
	}
}

// Get - returns random IP address if more than 1
func (t *TCPMap) Get(hostport string, randomize bool) (*net.TCPAddr, error) {
	addresses, _ := (*t)[hostport]
	size := len(addresses)
	if size == 0 {
		return nil, ErrEmptyIPS
	}
	if randomize {
		return addresses[rand.Intn(size)], nil
	}
	return addresses[0], nil
}

// BulkAdd - convenience method to iterate array of IP's and Add them
func (t *TCPMap) BulkAdd(hostport string, ips []net.IP, port int) {
	for _, ip := range ips {
		t.Add(hostport, ip, port)
	}
}

func (t *TCPMap) append(hostport string, ip net.IP, port int) {
	addresses, _ := (*t)[hostport]
	(*t)[hostport] = append(addresses, &net.TCPAddr{IP: ip, Port: port})
}

// Exist - checks whether hostport exist
func (t *TCPMap) Exist(hostport string) bool {
	_, ok := (*t)[hostport]
	return ok
}

// Len - returns size of IP addresses for hostport
func (t *TCPMap) Len(hostport string) int {
	if t.Exist(hostport) {
		addresses, _ := (*t)[hostport]
		return len(addresses)
	}
	return 0
}
