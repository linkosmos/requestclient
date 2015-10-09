package godns

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/miekg/dns"
)

// DefaultNameServer - Google
const DefaultNameServer = "8.8.8.8:53"

// Pool - for caching DNS records and easy retrieval
type Pool struct {
	NameServer                      string
	Randomize                       bool
	records                         TCPMap
	Timeout, StepTimeout, RetryWait time.Duration
}

// New - returns new dns records pool
func New() *Pool {
	return &Pool{
		NameServer:  DefaultNameServer,
		Randomize:   true,
		Timeout:     5 * time.Second,
		StepTimeout: 2 * time.Second,
		RetryWait:   2 * time.Second,
		records:     make(TCPMap, 5),
	}
}

// Get - retuns first or random IP assigned to hostport
func (p *Pool) Get(hostport string) (*net.TCPAddr, error) {
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		return nil, err
	}
	if p.records.Exist(hostport) {
		return p.records.Get(hostport, p.Randomize)
	}
	ips, duration, err := p.ResolveName(host, p.NameServer)
	if err != nil {
		return nil, err
	}
	logrus.Warningf("%s took %f, to resolve %s", hostport, duration.Seconds(), ips)
	if len(ips) == 0 {
		return nil, ErrEmptyIPS
	}
	portNum, _ := strconv.Atoi(port)
	p.records.BulkAdd(hostport, ips, portNum)
	return p.records.Get(hostport, p.Randomize)
}

// ResolveName - resolves name for given host and returns array of IP's
func (p *Pool) ResolveName(name, nameserver string) (addrs []net.IP, dur time.Duration, err error) {
	dnsClient := &dns.Client{
		Net:          "tcp",
		ReadTimeout:  p.StepTimeout,
		WriteTimeout: p.StepTimeout,
	}
	dnsMessage := new(dns.Msg)
	dnsMessage.MsgHdr.RecursionDesired = true
	dnsMessage.SetQuestion(dns.Fqdn(name), dns.TypeA)
	addrs = make([]net.IP, 0, 5)
	retryWait := p.RetryWait

Redo:
	var reply *dns.Msg
	var rtt time.Duration
	reply, rtt, err = dnsClient.Exchange(dnsMessage, nameserver)
	dur += rtt
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
			if dur+retryWait < p.Timeout {
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
