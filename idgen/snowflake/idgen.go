package snowflake

import (
	"time"

	"github.com/bwmarrin/snowflake"

	"github.com/agztizoo/glue/idgen"
)

// New 创建雪花 ID 生成器.
func New(node int64) idgen.IDGenerator {
	gen, err := snowflake.NewNode(node)
	if err != nil {
		panic(err)
	}
	return idgen.New(func() (int64, error) {
		return gen.Generate().Int64(), nil
	}, TimeOfID)
}

// TimeOfID 解析 snowflake 生成 ID 的时间信息.
func TimeOfID(id int64) (time.Time, error) {
	// NodeBits + StepBits
	// t := snowflake.ID(id).Time()
	t := id>>(snowflake.NodeBits+snowflake.StepBits) + snowflake.Epoch
	s2m := int64(time.Second / time.Millisecond)
	sec := t / s2m
	nsec := (t % s2m) * int64(time.Millisecond)
	return time.Unix(sec, nsec), nil
}
