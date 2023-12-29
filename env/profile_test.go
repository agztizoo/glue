package env

import (
	"os"
	"testing"
)

func TestProfile(t *testing.T) {
	t.Run("OnlyShard", func(t *testing.T) {
		profile = ""
		shard := "OnlyShard"
		os.Setenv("TEST_SHARD_INDEX", shard)
		p := Profile()
		if p != "" {
			t.Errorf("expect:%s, got: %s", shard, p)
		}
	})
	t.Run("OnlyProfile", func(t *testing.T) {
		profile = "OnlyProfile"
		os.Setenv("TEST_SHARD_INDEX", "")
		p := Profile()
		if p != profile {
			t.Errorf("expect:%s, got: %s", profile, p)
		}
	})
	t.Run("ProfileAndShard", func(t *testing.T) {
		profile = "Profile"
		shard := "Shard"
		expect := "Profile"
		os.Setenv("TEST_SHARD_INDEX", shard)
		p := Profile()
		if p != expect {
			t.Errorf("expect:%s, got: %s", expect, p)
		}
	})
}

func TestShardProfile(t *testing.T) {
	t.Run("NoProfileAndNoShard", func(t *testing.T) {
		profile = ""
		os.Setenv("TEST_SHARD_INDEX", "")
		p := ShardProfile()
		if p != "" {
			t.Errorf("expect:%s, got: %s", profile, p)
		}
	})
	t.Run("OnlyProfile", func(t *testing.T) {
		profile = "OnlyProfile"
		os.Setenv("TEST_SHARD_INDEX", "")
		p := ShardProfile()
		if p != profile {
			t.Errorf("expect:%s, got: %s", profile, p)
		}
	})
	t.Run("OnlyShard", func(t *testing.T) {
		profile = ""
		shard := "OnlyShard"
		os.Setenv("TEST_SHARD_INDEX", shard)
		p := ShardProfile()
		if p != shard {
			t.Errorf("expect:%s, got: %s", shard, p)
		}
	})
	t.Run("ProfileAndShard", func(t *testing.T) {
		profile = "Profile"
		shard := "Shard"
		expect := "Profile-Shard"
		os.Setenv("TEST_SHARD_INDEX", shard)
		p := ShardProfile()
		if p != expect {
			t.Errorf("expect:%s, got: %s", expect, p)
		}
	})
}
