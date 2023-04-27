package user

import (
	"math/rand"
	"time"
)

type Sign struct {
	Points int
	LastDo time.Time
	Signed bool
}

func (s *Sign) isSigned() bool {
	if s.isTimeOk() {
		s.Signed = true
	} else {
		s.Signed = false
	}
	return s.Signed
}

func (s *Sign) pointAdd() int {
	s.LastDo = time.Now()
	s.Signed = true
	n := rand.Intn(3)
	s.Points += n
	return n
}

func (s *Sign) showPoints() int {
	return s.Points
}

func (s *Sign) buy(point int) {
	s.Points -= point
}

func (s *Sign) isTimeOk() bool {
	//获得上次操作时间 转化为整点
	LastDoTimeStamp := s.LastDo.Truncate(24 * time.Hour)
	//获取现在时间 转化为整点
	now := time.Now()
	zeroClock := now.Truncate(24 * time.Hour)
	if LastDoTimeStamp == zeroClock {
		return false
	}
	return true
}
