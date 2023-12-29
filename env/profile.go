package env

import (
	"fmt"
	"os"
)

// 链接时设定 Profile.
//
// 用于程序 link 阶段设置不同 Profile.
//
// 不要改变 profile 的命名,改变命名会影响程序链接.
var profile string

// Profile 返回原始 profile 值.
//
// 返回链接阶段设置的 profile 值.
func Profile() string {
	return profile
}

// ShardProfile 返回分片 profile 值.
//
// 同时包含分片和 profile 时,返回分片 profile.
// 不包含分片时, 返回原始 profile.
func ShardProfile() string {
	shard := getShard()
	if shard == "" {
		return profile
	}
	if profile == "" {
		return shard
	}
	return fmt.Sprintf("%s-%s", profile, shard)
}

func getShard() string {
	return os.Getenv("TEST_SHARD_INDEX")
}
