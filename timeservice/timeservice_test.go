package timeservice

import (
	"testing"
	"time"
)

func TestTimeService(t *testing.T) {
	ts := NewTimeService()
	t1 := ts.Now()
	t2 := time.Now()

	d := t2.Sub(t1)

	if d < 0 {
		t.Errorf("expect duration great or equal 0, got: %v", d)
	}
	if d > time.Second {
		t.Errorf("expect duration less than 1 second, got: %v", d)
	}
}
