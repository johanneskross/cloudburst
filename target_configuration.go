package cloudburst

import ()

type TargetConfiguration struct {
	TargetIp                           string
	Offset, RampUp, Duration, RampDown int64
	WorkloadProfileIndex               int
	WorkloadProfileName                string
	WorkloadProfileOffset              int64
	TargetFactory                      TargetFactory
}

func NewTargetConfiguration(
	targetIp string,
	offset, rampUp, duration, rampDown int64,
	workloadProfileIndex int,
	workloadProfileName string,
	workloadProfileOffset int64,
	targetFactory TargetFactory) *TargetConfiguration {

	return &TargetConfiguration{
		targetIp,
		offset * TO_NANO,
		rampUp * TO_NANO,
		duration * TO_NANO,
		rampDown * TO_NANO,
		workloadProfileIndex,
		workloadProfileName,
		workloadProfileOffset * TO_NANO,
		targetFactory}
}
