package redirect

import (
	"fmt"
	"net/http"
)

// Redirect - stops following redirects after specified number of times
// compatible with http.Client
type Redirect struct {
	StopAfter int
}

// DefaultCheckRedirect - default check redirect function from http.Client
func (r *Redirect) DefaultCheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= r.StopAfter {
		return fmt.Errorf("Stopped after %d redirects")
	}
	return nil
}
