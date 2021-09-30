package app

import "time"

const TICK_DURATION = time.Second / 20

type TickChecker struct {
	tick     uint64
	duration uint64
	cur      uint64
}

func NewTickChecker(tick uint64) *TickChecker {
	tc := new(TickChecker)

	tc.tick = tick
	tc.duration = tc.tick * uint64(TICK_DURATION)
	tc.cur = 0

	return tc
}

func (tc *TickChecker) Reset() {
	tc.cur = 0
}

func (tc *TickChecker) Next(t time.Duration) bool {
	tc.cur += uint64(t.Nanoseconds())

	// 到达
	if tc.cur >= tc.duration {
		tc.cur -= tc.duration
		return true
	}

	return false
}
