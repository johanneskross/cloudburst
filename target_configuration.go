package cloudburst

import (
	"github.com/johanneskross/cloudburst/times"
)

type TargetConfiguration struct {
	TargetId                          int
	TargetIp                          string
	Delay, RampUp, Duration, RampDown int64
	TimeSeries                        times.TimeSeries
}

func NewTargetConfiguration(targetId int, targetIp string, delay, rampUp, duration, rampDown int64, timeSeries times.TimeSeries) *TargetConfiguration {
	return &TargetConfiguration{targetId, targetIp, delay * TO_NANO, rampUp * TO_NANO, duration * TO_NANO, rampDown * TO_NANO, timeSeries}
}
