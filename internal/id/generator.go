package id

import (
	"errors"
	"sync"
	"time"
)

const (
	// 自定义纪元时间（例如：2020-01-01 UTC）
	epoch = int64(1577836800000) // 2020-01-01 00:00:00 UTC in milliseconds

	nodeBits = uint64(10) // 节点 ID 位数
	stepBits = uint64(12) // 序列号位数

	nodeMax = -1 ^ (-1 << nodeBits)
	stepMax = -1 ^ (-1 << stepBits)

	nodeMask = nodeMax << stepBits
	stepMask = stepMax
)

// Snowflake 雪花 ID 生成器
type Snowflake struct {
	mu        sync.Mutex
	timestamp int64
	node      int64
	step      int64
}

// NewSnowflake 创建一个新的雪花生成器
func NewSnowflake(node int64) (*Snowflake, error) {
	if node < 0 || node > nodeMax {
		return nil, errors.New("node ID out of range")
	}
	return &Snowflake{
		timestamp: 0,
		node:      node,
		step:      0,
	}, nil
}

// Generate 生成下一个雪花 ID
func (s *Snowflake) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()

	if s.timestamp == now {
		// 同一毫秒内，序列号递增
		s.step = (s.step + 1) & stepMax
		if s.step == 0 {
			// 序列号溢出，等待下一毫秒
			for now <= s.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// 新的毫秒，重置序列号
		s.step = 0
	}

	s.timestamp = now

	id := ((now - epoch) << (nodeBits + stepBits)) |
		((s.node & nodeMax) << stepBits) |
		(s.step & stepMask)

	return id
}
