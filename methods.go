package requestclient

import (
	"net/http"
	"net/url"
)

// Request methods
const (
	GET    = "GET"
	HEAD   = "HEAD"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

// GET - request
func (r *RequestClient) GET(u *url.URL) *http.Request {
	return r.NewRequest(GET, u, nil)
}

// HEAD - request
func (r *RequestClient) HEAD(u *url.URL) *http.Request {
	return r.NewRequest(HEAD, u, nil)
}
