package dns

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/miekg/dns"
)

// DefaultNameServer - Google
const DefaultNameServer = "8.8.8.8:53"

// -
var (
	Timeout     = 20 * time.Second
	StepTimeout = 2 * time.Second
	RetryWait   = 2 * time.Second

	ErrEmptyIPS = errors.New("No IP's for a given host")
)

// Pool - for caching DNS records and easy retrieval
type Pool struct {
	NameServer string
	records    map[string]*net.TCPAddr
}

// New - returns new dns records pool
func New() *Pool {
	return &Pool{
		NameServer: DefaultNameServer,
		records:    make(map[string]*net.TCPAddr),
	}
}

// Get - retuns random IP from map for a given host
func (p *Pool) Get(hostport string) (*net.TCPAddr, error) {
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		return nil, err
	}
	if value, ok := p.records[hostport]; ok {
		return value, nil
	}
	ips, _, err := ResolveName(host, p.NameServer)
	if err != nil {
		return nil, err
	}
	if len(ips) == 0 {
		return nil, ErrEmptyIPS
	}
	po, _ := strconv.Atoi(port)
	p.records[hostport] = &net.TCPAddr{IP: ips[rand.Intn(len(ips))], Port: po}
	return p.Get(hostport)
}

// ResolveName - resolves name for given host and returns array of IP's
func ResolveName(name, nameserver string) (addrs []net.IP, dur time.Duration, err error) {
	dnsClient := &dns.Client{
		Net:          "tcp",
		ReadTimeout:  StepTimeout,
		WriteTimeout: StepTimeout,
	}
	dnsMessage := new(dns.Msg)
	dnsMessage.MsgHdr.RecursionDesired = true
	dnsMessage.SetQuestion(dns.Fqdn(name), dns.TypeA)
	addrs = make([]net.IP, 0)
	retryWait := RetryWait

Redo:
	var reply *dns.Msg
	var rtt time.Duration
	reply, rtt, err = dnsClient.Exchange(dnsMessage, nameserver)
	dur += rtt
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
			if dur+retryWait < Timeout {
				time.Sleep(retryWait)
				retryWait *= 2
				goto Redo
			}
		}
		return nil, dur, err
	}
	if reply.Rcode != dns.RcodeSuccess {
		err = fmt.Errorf(`ResolveName(%s, %s): %s`, name, nameserver, dns.RcodeToString[reply.Rcode])
		return nil, dur, err
	}
	for _, a := range reply.Answer {
		if rra, ok := a.(*dns.A); ok {
			addrs = append(addrs, rra.A)
		}
		if rra6, ok := a.(*dns.AAAA); ok {
			addrs = append(addrs, rra6.AAAA)
		}
	}
	if reply.MsgHdr.Truncated {
		goto Redo
	}
	return
}
