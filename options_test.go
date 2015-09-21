package requestclient

import "testing"

func TestNewOptions(t *testing.T) {
	op := NewOptions()
	if op.DialerKeepAlive != DefaultDialerKeepAlive {
		t.Error("Expected - NewOptions() to initialize with default values")
	}

	if got := op.Headers.Get("Connection"); got != "Keep-Alive" {
		t.Errorf("Expected {Connection => Keep-Alive} got %s", got)
	}
}
