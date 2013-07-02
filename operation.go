package cloudburst

import ()

const NO_TRACE_LABEL = "[NONE]"
const STEADY_STATE_TRACE_LABEL = "[STEADY-STATE]"
const LATE_TRACE_LABEL = "[LATE]"
const RAMP_UP_TRACE_LABEL = "[RAMP-UP]"
const RAMP_DOWN_TRACE_LABEL = "[RAMP-DOWN]"

type Operation interface {
	Run()
	GetResults() OperationResult
}

type OperationResult struct {
	OperationIndex                                                int
	OperationName, OperationRequest                               string
	Failed                                                        bool
	LoadDefinition                                                int
	TimeStarted, TimeFinished, ProfileStartTime, ActionsPerformed int64
	TraceLabel                                                    string
}

func (or *OperationResult) ActionPerformed() {
	or.ActionsPerformed++
}

func (or *OperationResult) GetExecutionTime() int64 {
	return or.TimeFinished - or.TimeStarted
}
