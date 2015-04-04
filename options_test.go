package requestclient

import "testing"

func TestNewOptions(t *testing.T) {
	op := NewOptions()
	if op.DialerKeepAlive != DefaultDialerKeepAlive {
		t.Error("Expected - NewOptions() to initialize with default values")
	}
}
