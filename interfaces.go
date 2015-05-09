package requestclient

import (
	"net"
	"net/http"
)

// DialDialer - is an interface for connecting to an address.
type DialDialer interface {
	Dial(network, address string) (net.Conn, error)
}

// RoundTripper is an interface representing the ability to execute a
// single HTTP transaction, obtaining the Response for a given Request.
//
// A RoundTripper must be safe for concurrent use by multiple
// goroutines.
type RoundTripper interface {
	// RoundTrip executes a single HTTP transaction, returning
	// the Response for the request req.  RoundTrip should not
	// attempt to interpret the response.  In particular,
	// RoundTrip must return err == nil if it obtained a response,
	// regardless of the response's HTTP status code.  A non-nil
	// err should be reserved for failure to obtain a response.
	// Similarly, RoundTrip should not attempt to handle
	// higher-level protocol details such as redirects,
	// authentication, or cookies.
	//
	// RoundTrip should not modify the request, except for
	// consuming and closing the Body, including on errors. The
	// request's URL and Header fields are guaranteed to be
	// initialized.
	RoundTrip(*http.Request) (*http.Response, error)
}

// ClientRequester - higher level API that implements http.Client
type ClientRequester interface {
	Do(req *http.Request) (*http.Response, error)
}
