package communication

import "time"

const (
	// DefaultRetryInterval for WaitForPing.
	DefaultRetryInterval = 250 * time.Millisecond

	// DefaultRetryTimeout for WaitForPing.
	DefaultRetryTimeout = 10 * time.Second
)

// WaitForPingAuto runs a ping function until no error is returned, or a timeout is reached. The error returned is the
// one of the last ping performed.
//
// It is a pre-configured wrapper for WaitForPing.
//
//	func ServerHealth() error {
//	  // Returns an error if server response is not 200.
//	  return nil
//	}
//
//	communication.WaitForPingAuto(ServerHealth)
func WaitForPingAuto(ping func() error) error {
	return WaitForPing(ping, DefaultRetryTimeout, DefaultRetryInterval)
}

// WaitForPing runs a ping function until no error is returned, or a timeout is reached. The error returned is the one
// of the last ping performed.
//
//	func ServerHealth() error {
//	  // Returns an error if server response is not 200.
//	  return nil
//	}
//
//	communication.WaitForPing(ServerHealth, time.Second, 100 * time.Millisecond)
//
// You can use WaitForPingAuto for pre-configured values.
func WaitForPing(ping func() error, timeout time.Duration, retryInterval time.Duration) error {
	err := ping()
	start := time.Now()

	for ; time.Since(start) < timeout && err != nil; err = ping() {
		time.Sleep(retryInterval)
	}

	return err
}
