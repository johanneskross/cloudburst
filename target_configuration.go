package cloudburst

import (
	"github.com/johanneskross/times"
)

type TargetConfiguration struct {
	TargetId                          int
	TargetIp                          string
	Delay, RampUp, Duration, RampDown int
	TimeSeries                        times.TimeSeries
}

func NewTargetConfiguration(targetId int, targetIp string, delay, rampUp, duration, rampDown int, timeSeries times.TimeSeries) *TargetConfiguration {
	return &TargetConfiguration{targetId, targetIp, delay, rampUp, duration, rampDown, timeSeries}
}
