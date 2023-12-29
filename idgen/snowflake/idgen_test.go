package snowflake

import (
	"testing"
	"time"
)

func Test_TimeOfID(t *testing.T) {
	var id int64 = 1365258938810892288
	t1, err := TimeOfID(id)
	if err != nil {
		t.Fatalf("expect time parse not error: %v", err)
	}
	t2 := time.Unix(1614338070, 954000000)
	if !t1.Equal(t2) {
		t.Errorf("expect time of id(%d): %v got: %v", id, t1, t2)
	}
}

func Test_GenerateID(t *testing.T) {
	g := New(0)
	id1 := g.MustGenIntID()
	id2 := g.MustGenIntID()

	if id1 == id2 {
		t.Errorf("expect id1: %d not equals to id2 %d", id1, id2)
	}
}

func Test_New_Panic(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Error("expect panic for node: -1")
		}
	}()
	New(-1)
}
