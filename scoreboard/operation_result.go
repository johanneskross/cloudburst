package scoreboard

import ()

type OperationResult struct {
	OperationIndex                                                int
	OperationName, OperationRequest                               string
	Failed                                                        bool
	LoadDefinition                                                int
	TimeStarted, TimeFinished, ProfileStartTime, ActionsPerformed int64
	TraceLabel                                                    string
}

func (or *OperationResult) ActionPerformed(amount int) {
	or.ActionsPerformed += int64(amount)
}

func (or *OperationResult) GetExecutionTime() int64 {
	return or.TimeFinished - or.TimeStarted
}
