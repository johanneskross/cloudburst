package cloudburst

import (
	"github.com/johanneskross/times"
)

type TargetConfiguration struct {
	Name                              string
	Delay, RampUp, Duration, RampDown int
	TimeSeries                        times.TimeSeries
}

func NewTargetConfiguration(name string, delay, rampUp, duration, rampDown int, timeSeries times.TimeSeries) *TargetConfiguration {
	return &TargetConfiguration{name, delay, rampUp, duration, rampDown, timeSeries}
}
