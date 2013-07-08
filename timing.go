package cloudburst

import (
	"time"
)

type Timing struct {
	Start, StartSteadyState, EndSteadyState, EndRun int64
	RampUp, Duration, RampDown                      int64
}

func NewTiming(rampUp, duration, rampDown int64) *Timing {
	start := time.Now().UnixNano()
	startSteadyState := start + rampUp
	endSteadyState := startSteadyState + duration
	endRun := endSteadyState + rampDown
	return &Timing{start, startSteadyState, endSteadyState, endRun, rampUp, duration, rampDown}
}

func (t *Timing) GetNewTiming() *Timing {
	return NewTiming(t.RampUp, t.Duration, t.RampDown)
}

func (t *Timing) SteadyStateDuration() int64 {
	return t.EndSteadyState - t.StartSteadyState
}

func (t *Timing) InRampUp(now int64) bool {
	return now < t.StartSteadyState
}

func (t *Timing) InSteadyState(time int64) bool {
	return time >= t.StartSteadyState && time <= t.EndSteadyState
}

func (t *Timing) InRampDown(time int64) bool {
	return time > t.EndSteadyState
}
