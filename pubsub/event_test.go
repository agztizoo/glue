package pubsub

import (
	"reflect"
	"testing"
	"time"
)

func TestPredefinedHeaderGetter(t *testing.T) {
	tm := time.Now()
	ev := NewEvent("1", "payload", tm)

	t.Run("get id", func(t *testing.T) {
		if ev.GetID() != "1" {
			t.Errorf("expect id: %s, got: %s", "1", ev.GetID())
		}
	})

	t.Run("get payload", func(t *testing.T) {
		if ev.GetPayload() != "payload" {
			t.Errorf("expect payload: %v, got: %v", "2", ev.GetPayload())
		}
	})

	t.Run("get create time", func(t *testing.T) {
		if !ev.GetCreateTime().Equal(tm) {
			t.Errorf("expect create time: %v, got: %v", tm, ev.GetCreateTime())
		}
	})

	t.Run("get aggregate type", func(t *testing.T) {
		t.Run("empty aggregate type", func(t *testing.T) {
			if ev.GetAggregateType() != "" {
				t.Errorf("expect empty aggregate type, got: %s", ev.GetAggregateType())
			}
		})

		t.Run("none string aggregate type", func(t *testing.T) {
			ev.SetHeader(EventHeaderAggregateType, 1)
			if ev.GetAggregateType() != "" {
				t.Errorf("expect empty aggregate type, got: %s", ev.GetAggregateType())
			}
		})

		t.Run("none empty aggregate type", func(t *testing.T) {
			ev.SetHeader(EventHeaderAggregateType, "Order")
			if ev.GetAggregateType() != "Order" {
				t.Errorf("expect aggregate type: %s, got: %s", "Order", ev.GetAggregateType())
			}
		})
	})

	t.Run("get aggregate id", func(t *testing.T) {
		t.Run("empty aggregate id", func(t *testing.T) {
			if ev.GetAggregateID() != "" {
				t.Errorf("expect empty aggregate id, got: %s", ev.GetAggregateID())
			}
		})

		t.Run("none string aggregate id", func(t *testing.T) {
			ev.SetHeader(EventHeaderAggregateID, 1)
			if ev.GetAggregateID() != "" {
				t.Errorf("expect empty aggregate id, got: %s", ev.GetAggregateID())
			}
		})

		t.Run("none empty aggregate id", func(t *testing.T) {
			ev.SetHeader(EventHeaderAggregateID, "OrderID_1")
			if ev.GetAggregateID() != "OrderID_1" {
				t.Errorf("expect aggregate id: %s, got: %s", "OrderID_1", ev.GetAggregateID())
			}
		})
	})

	t.Run("get tenant", func(t *testing.T) {
		t.Run("empty tenant", func(t *testing.T) {
			if ev.GetTenantID() != "" {
				t.Errorf("expect empty tenant, got: %s", ev.GetTenantID())
			}
		})

		t.Run("none string tenant", func(t *testing.T) {
			ev.SetHeader(EventHeaderTenantID, 1)
			if ev.GetTenantID() != "" {
				t.Errorf("expect empty tenant, got: %s", ev.GetTenantID())
			}
		})

		t.Run("none empty tenant", func(t *testing.T) {
			ev.SetHeader(EventHeaderTenantID, "tenant_1")
			if ev.GetTenantID() != "tenant_1" {
				t.Errorf("expect tenant: %s, got: %s", "tenant_1", ev.GetTenantID())
			}
		})
	})

	t.Run("get user", func(t *testing.T) {
		t.Run("empty user", func(t *testing.T) {
			if ev.GetUserID() != "" {
				t.Errorf("expect empty user, got: %s", ev.GetUserID())
			}
		})

		t.Run("none string user", func(t *testing.T) {
			ev.SetHeader(EventHeaderUserID, 1)
			if ev.GetUserID() != "" {
				t.Errorf("expect empty user, got: %s", ev.GetUserID())
			}
		})

		t.Run("none empty user", func(t *testing.T) {
			ev.SetHeader(EventHeaderUserID, "user_1")
			if ev.GetUserID() != "user_1" {
				t.Errorf("expect user: %s, got: %s", "user_1", ev.GetUserID())
			}
		})
	})
}

func TestGetPayload(t *testing.T) {
	ev := NewEvent("1", "payload", time.Now())

	t.Run("get payload", func(t *testing.T) {
		if ev.GetPayload() != "payload" {
			t.Errorf("expect payload: %v, got: %v", "payload", ev.GetPayload())
		}
	})

	t.Run("reset payload", func(t *testing.T) {
		ev.SetPayload("reset payload")
		if ev.GetPayload() != "reset payload" {
			t.Errorf("expect payload: %v, got: %v", "reset payload", ev.GetPayload())
		}
	})
}

func TestGetHeader(t *testing.T) {
	ev := NewEvent("1", "payload", time.Now())
	ev.SetHeader("1", 1)
	ev.SetHeader("2", "2")
	ev.SetHeader("2", 3)

	t.Run("get header", func(t *testing.T) {
		if ev.GetHeader("2") != 3 {
			t.Errorf("expect: %v, got: %v", 3, ev.GetHeader("2"))
		}
	})

	t.Run("get headers", func(t *testing.T) {
		hs := ev.GetHeaders()
		ex := map[string]interface{}{
			"1": 1,
			"2": 3,
		}
		if !reflect.DeepEqual(hs, ex) {
			t.Errorf("expect: %v, got: %v", ex, hs)
		}
	})
}
