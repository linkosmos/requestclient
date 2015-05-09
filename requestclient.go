package requestclient

import (
	"crypto/tls"
	"net/http"

	"github.com/ernestas-poskus/httpcontrol"
	"github.com/linkosmosis/requestclient/dialer"
)

// -
const (
	RequestProto      = "HTTP/1.1"
	RequestProtoMinor = 0
	RequestProtoMajor = 1
)

// RequestClient - http.: Transport, Client wrapper
type RequestClient struct {

	// A Header represents the key-value pairs in an HTTP header.
	// Used for every http.Request formed by RequestClient
	Headers http.Header // requesting headers

	RequestProto string

	RequestProtoMinor, RequestProtoMajor int

	// A Config structure is used to configure a TLS client or server.
	// After one has been passed to a TLS function it must not be
	// modified. A Config may be reused; the tls package will also not
	// modify it.
	TLS *tls.Config

	// A Dialer contains options for connecting to an address.
	//
	// The zero value for each field is equivalent to dialing
	// without that option. Dialing with the zero value of Dialer
	// is therefore equivalent to just calling the Dial function.
	Dialer DialDialer

	// Transport is an implementation of RoundTripper that supports HTTP,
	// HTTPS, and HTTP proxies (for either HTTP or HTTPS with CONNECT).
	// Transport can also cache connections for future re-use.
	Transport RoundTripper

	// A Client is an HTTP client. Its zero value (DefaultClient) is a
	// usable client that uses DefaultTransport.
	//
	// The Client's Transport typically has internal state (cached TCP
	// connections), so Clients should be reused instead of created as
	// needed. Clients are safe for concurrent use by multiple goroutines.
	//
	// A Client is higher-level than a RoundTripper (such as Transport)
	// and additionally handles HTTP details such as cookies and
	// redirects.
	Client ClientRequester
}

// SetupNewDialer - setups new requestclient Dialer with specified options
func SetupNewDialer(op *Options) (dl *dialer.Dialer) {
	dl = dialer.New()
	dl.Timeout = op.DialerTimeout
	dl.Deadline = op.DialerDeadline
	dl.DualStack = op.DialerDualStack
	dl.KeepAlive = op.DialerKeepAlive
	dl.DualStack = op.DialerDualStack
	return dl
}

// New - returns Request Client
func New(op *Options) (r *RequestClient) {
	if op == nil {
		op = NewOptions()
	}
	r = &RequestClient{
		Headers:           op.Headers,
		RequestProto:      RequestProto,
		RequestProtoMinor: RequestProtoMinor,
		RequestProtoMajor: RequestProtoMajor,
		TLS: &tls.Config{
			InsecureSkipVerify: op.TLSInsecureSkipVerify,
		},
		Dialer: SetupNewDialer(op), // Setting Dialer
	}

	// Setting up TRANSPORT
	r.Transport = &httpcontrol.Transport{
		Dial:                  r.Dialer.Dial,
		TLSClientConfig:       r.TLS,
		MaxTries:              op.TransportMaxTries,
		DisableKeepAlives:     op.TransportDisableKeepAlives,
		DisableCompression:    op.TransportDisableCompression,
		MaxIdleConnsPerHost:   op.TransportMaxIdleConnsPerHost,
		RequestTimeout:        op.TransportRequestTimeout,
		ResponseHeaderTimeout: op.TransportResponseHeaderTimeout,
	}

	// Setting up CLIENT, higher level API of TRANSPORT
	r.Client = &http.Client{
		Transport: r.Transport,
		Timeout:   op.ClientTimeout,
	}
	return r
}

// Do - sends an HTTP request and returns an HTTP response, following
// policy (e.g. redirects, cookies, auth) as configured on the client.
//
// An error is returned if caused by client policy (such as
// CheckRedirect), or if there was an HTTP protocol error.
// A non-2xx response doesn't cause an error.
//
// When err is nil, resp always contains a non-nil resp.Body.
//
// Callers should close resp.Body when done reading from it. If
// resp.Body is not closed, the Client's underlying RoundTripper
// (typically Transport) may not be able to re-use a persistent TCP
// connection to the server for a subsequent "keep-alive" request.
//
// The request Body, if non-nil, will be closed by the underlying
// Transport, even on errors.
//
// Generally Get, Post, or PostForm will be used instead of Do.
func (r *RequestClient) Do(req *http.Request) (*http.Response, error) {
	return r.Client.Do(req)
}

// RoundTrip implements the RoundTripper interface.
//
// For higher-level HTTP client support (such as handling of cookies
// and redirects), see Get, Post, and the Client type.
func (r RequestClient) RoundTrip(req *http.Request) (*http.Response, error) {
	return r.Transport.RoundTrip(req)
}
