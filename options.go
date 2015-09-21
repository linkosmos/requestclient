package requestclient

import (
	"fmt"
	"net/http"
	"time"
)

// DefaultUserAgent - default user agent for this request client package
var DefaultUserAgent = fmt.Sprintf("Go RequestClient %s", VERSION)

// RequestClient options default values
const (
	DefaultDialerTimeout                = 30 * time.Second
	DefaultDialerDualStack              = false
	DefaultDialerKeepAlive              = 30 * time.Second
	DefaultTransportMaxTries            = 3
	DefaultTransportDisableKeepAlives   = false
	DefaultTransportDisableCompression  = false
	DefaultTransportMaxIdleConnsPerHost = http.DefaultMaxIdleConnsPerHost
	DefaultClientTimeout                = 1 * time.Minute // Default 3 min
	DefaultTLSInsecureSkipVerify        = false
)

// NewOptions - options struct initialized with default values
func NewOptions() (op *Options) {
	op = &Options{
		Headers:                      make(http.Header),
		DialerTimeout:                DefaultDialerTimeout,
		DialerDeadline:               time.Now().Add(DefaultDialerTimeout),
		DialerDualStack:              DefaultDialerDualStack,
		DialerKeepAlive:              DefaultDialerKeepAlive,
		TransportMaxTries:            DefaultTransportMaxTries,
		TransportDisableKeepAlives:   DefaultTransportDisableKeepAlives,
		TransportDisableCompression:  DefaultTransportDisableCompression,
		TransportMaxIdleConnsPerHost: DefaultTransportMaxIdleConnsPerHost,
		ClientTimeout:                DefaultClientTimeout,
		TLSInsecureSkipVerify:        DefaultTLSInsecureSkipVerify,
	}
	op.SetUserAgent(DefaultUserAgent)
	if !op.TransportDisableKeepAlives {
		op.Headers.Add("Connection", "Keep-Alive")
	}
	return op
}

// Options - keep options for continuous requests
type Options struct {

	////////////////////////////////
	// Dialer
	////////////////////////////////
	// Timeout is the maximum amount of time a dial will wait for
	// a connect to complete.
	//
	// With or without a timeout, the operating system may impose
	// its own earlier timeout. For instance, TCP timeouts are
	// often around 3 minutes.
	DialerTimeout time.Duration // The default is no timeout.

	// Deadline is the absolute point in time after which dials
	// will fail. If Timeout is set, it may fail earlier.
	// Zero means no deadline, or dependent on the operating system
	// as with the Timeout option.
	DialerDeadline time.Time

	// DualStack allows a single dial to attempt to establish
	// multiple IPv4 and IPv6 connections and to return the first
	// established connection when the network is "tcp" and the
	// destination is a host name that has multiple address family
	// DNS records.
	DialerDualStack bool

	// KeepAlive specifies the keep-alive period for an active
	// network connection.
	// If zero, keep-alives are not enabled. Network protocols
	// that do not support keep-alives ignore this field.
	DialerKeepAlive time.Duration

	//
	////////////////////////////////
	// Transport
	////////////////////////////////
	//

	// MaxTries, if non-zero, specifies the number of times we will retry on
	// failure. Retries are only attempted for temporary network errors or known
	// safe failures.
	TransportMaxTries uint

	// DisableKeepAlives, if true, prevents re-use of TCP connections
	// between different HTTP requests.
	TransportDisableKeepAlives bool

	// DisableCompression, if true, prevents the Transport from
	// requesting compression with an "Accept-Encoding: gzip"
	TransportDisableCompression bool

	// MaxIdleConnsPerHost, if non-zero, controls the maximum idle
	// (keep-alive) to keep per-host.
	TransportMaxIdleConnsPerHost int

	// RequestTimeout, if non-zero, specifies the amount of time for the entire
	// request. This includes dialing (if necessary), the response header as well
	// as the entire body.
	TransportRequestTimeout time.Duration

	// ResponseHeaderTimeout, if non-zero, specifies the amount of
	// time to wait for a server's response headers after fully
	// writing the request (including its body, if any). This
	// time does not include the time to read the response body.
	TransportResponseHeaderTimeout time.Duration

	//
	////////////////////////////////
	// Client
	////////////////////////////////
	//

	// Timeout specifies a time limit for requests made by this
	// Client. The timeout includes connection time, any
	// redirects, and reading the response body. The timer remains
	// running after Get, Head, Post, or Do return and will
	// interrupt reading of the Response.Body.
	//
	// A Timeout of zero means no timeout.
	//
	// The Client's Transport must support the CancelRequest
	// method or Client will return errors when attempting to make
	// a request with Get, Head, Post, or Do. Client's default
	// Transport (DefaultTransport) supports CancelRequest.
	ClientTimeout time.Duration

	//
	////////////////////////////////
	// TLS
	////////////////////////////////
	//

	// InsecureSkipVerify controls whether a client verifies the
	// server's certificate chain and host name.
	TLSInsecureSkipVerify bool

	//
	////////////////////////////////
	// Request Options
	////////////////////////////////
	//

	// Request headers
	Headers http.Header
}

//
////////////////////////////////
// Setters
////////////////////////////////
//

// SetUserAgent - sets request user agent
func (o *Options) SetUserAgent(ua string) {
	o.Headers.Add("User-Agent", ua)
}

//
////////////////////////////////
// Getters
////////////////////////////////
//

// GetUserAgent - sets request user agent
func (o *Options) GetUserAgent() (ua string) {
	ua = o.Headers.Get("User-Agent")
	if ua == "" {
		return DefaultUserAgent
	}
	return ua
}
