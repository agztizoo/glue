package idgen

import (
	"errors"
	"strconv"
	"testing"
	"time"
)

func Test_Normal(t *testing.T) {
	tm := time.Now()
	id := tm.Unix()<<32 + 123456
	sid := strconv.FormatInt(id, 10)
	idg := New(func() (int64, error) {
		return id, nil
	}, func(iid int64) (time.Time, error) {
		if iid == id {
			return tm, nil
		}
		return time.Time{}, errors.New("not generated id")
	})

	t.Run("GetID", func(t *testing.T) {
		rsid, err := idg.GenID()
		if err != nil {
			t.Error(err)
		}
		if rsid != sid {
			t.Errorf("expect id: %s, got: %s", sid, rsid)
		}
	})

	t.Run("MustGenID", func(t *testing.T) {
		rsid := idg.MustGenID()
		if rsid != sid {
			t.Errorf("expect id: %s, got: %s", sid, rsid)
		}
	})

	t.Run("GenIntID", func(t *testing.T) {
		rid, err := idg.GenIntID()
		if err != nil {
			t.Error(err)
		}
		if rid != id {
			t.Errorf("expect id: %d, got: %d", id, rid)
		}
	})

	t.Run("MustGenIntID", func(t *testing.T) {
		rid := idg.MustGenIntID()
		if rid != id {
			t.Errorf("expect id: %d, got: %d", id, rid)
		}
	})

	t.Run("TimeOfID", func(t *testing.T) {
		rt, err := idg.TimeOfID(sid)
		if err != nil {
			t.Error(err)
		}
		if rt != tm {
			t.Errorf("expect id: %v, got: %v", rt, tm)
		}
	})

	t.Run("TimeOfIntID", func(t *testing.T) {
		rt, err := idg.TimeOfIntID(id)
		if err != nil {
			t.Error(err)
		}
		if rt != tm {
			t.Errorf("expect id: %v, got: %v", rt, tm)
		}
	})

	t.Run("MustTimeOfID", func(t *testing.T) {
		rt := idg.MustTimeOfID(sid)
		if rt != tm {
			t.Errorf("expect id: %v, got: %v", rt, tm)
		}
	})

	t.Run("MustTimeOfIntID", func(t *testing.T) {
		rt := idg.MustTimeOfIntID(id)
		if rt != tm {
			t.Errorf("expect id: %v, got: %v", rt, tm)
		}
	})
}

func Test_Error(t *testing.T) {
	tm := time.Now()
	id := tm.Unix()<<32 + 123456
	sid := strconv.FormatInt(id, 10)
	idg := New(func() (int64, error) {
		return id, errors.New("generate id error")
	}, func(iid int64) (time.Time, error) {
		return time.Time{}, errors.New("parse time error")
	})

	t.Run("GetID", func(t *testing.T) {
		_, err := idg.GenID()
		if err == nil {
			t.Error("expect GenID error")
		}
	})

	t.Run("MustGenID", func(t *testing.T) {
		defer func() {
			if e := recover(); e == nil {
				t.Error("expect MustGenID panic")
			}
		}()
		idg.MustGenID()
	})

	t.Run("GenIntID", func(t *testing.T) {
		_, err := idg.GenIntID()
		if err == nil {
			t.Error("expect GenIntID error")
		}
	})

	t.Run("MustGenIntID", func(t *testing.T) {
		defer func() {
			if e := recover(); e == nil {
				t.Error("expect MustGenIntID panic")
			}
		}()
		idg.MustGenIntID()
	})

	t.Run("TimeOfID", func(t *testing.T) {
		_, err := idg.TimeOfID(sid)
		if err == nil {
			t.Error("expect TimeOfID error")
		}
	})

	t.Run("TimeOfIntID", func(t *testing.T) {
		_, err := idg.TimeOfIntID(id)
		if err == nil {
			t.Error("expect TimeOfIntID error")
		}
	})

	t.Run("MustTimeOfID", func(t *testing.T) {
		defer func() {
			if e := recover(); e == nil {
				t.Error("expect MustTimeOfID panic")
			}
		}()
		idg.MustTimeOfID(sid)
	})

	t.Run("MustTimeOfIntID", func(t *testing.T) {
		defer func() {
			if e := recover(); e == nil {
				t.Error("expect MustTimeOfIntID panic")
			}
		}()
		idg.MustTimeOfIntID(id)
	})
}

func Test_Error_String_ID(t *testing.T) {
	tm := time.Now()
	id := tm.Unix()<<32 + 123456
	idg := New(func() (int64, error) {
		return id, nil
	}, func(iid int64) (time.Time, error) {
		return tm, nil
	})

	t.Run("TimeOfID", func(t *testing.T) {
		_, err := idg.TimeOfID("error_id")
		if err == nil {
			t.Error("expect TimeOfID error")
		}
	})

	t.Run("MustTimeOfID", func(t *testing.T) {
		defer func() {
			if e := recover(); e == nil {
				t.Error("expect MustTimeOfID panic")
			}
		}()
		idg.MustTimeOfID("error_id")
	})
}
