package lib

import (
	"errors"
	"sync"
	"time"
)

const (
	workerIDBits   uint8 = 10
	numberBits     uint8 = 12
	maxWorkerID          = -1 ^ (-1 << workerIDBits)
	maxNumber            = -1 ^ (-1 << numberBits)
	workerIDShift        = numberBits
	timestampShift       = numberBits + workerIDBits
	epoch                = 1745741259124
)

type Snowflake struct {
	mtx           *sync.Mutex
	workerID      uint64
	number        uint64
	lastTimestamp uint64
	sleepDuration time.Duration
}

var (
	instance *Snowflake
	once     sync.Once
)

func GetSnowflakeInstance(workerID uint64, sleepDuration time.Duration) (*Snowflake, error) {
	var err error
	once.Do(func() {
		instance, err = newSnowflake(workerID, sleepDuration)
	})
	if err != nil {
		return nil, err
	}
	return instance, nil

}

// NewSnowflake 创建一个雪花算法的生成器实例
// 传入工作节点ID和发生时钟回拨时的休眠时间
func newSnowflake(workerID uint64, sleepDuration time.Duration) (*Snowflake, error) {
	if workerID > maxWorkerID {
		return nil, errors.New("worker ID exceeds max value")
	}
	return &Snowflake{
		mtx:           new(sync.Mutex),
		workerID:      workerID,
		number:        0,
		lastTimestamp: 0,
		sleepDuration: sleepDuration,
	}, nil
}

func (s *Snowflake) GetID() (uint64, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	timestamp := uint64(time.Now().UnixMilli()) - epoch
	if timestamp < s.lastTimestamp {
		time.Sleep(s.sleepDuration)
		timestamp = uint64(time.Now().UnixMilli()) - epoch
		if timestamp < s.lastTimestamp {
			return 0, errors.New("clock moved backwards")
		}
	}

	if s.lastTimestamp == timestamp {
		s.number = (s.number + 1) & maxNumber
		if s.number == 0 {
			for timestamp <= s.lastTimestamp {
				timestamp = uint64(time.Now().UnixMilli()) - epoch
			}
		}
	} else {
		s.number = 0
	}

	s.lastTimestamp = timestamp

	id := ((timestamp << timestampShift) | (s.workerID << workerIDShift) | s.number)
	return id, nil
}
