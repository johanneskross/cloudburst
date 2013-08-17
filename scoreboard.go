package cloudburst

import ()

type Scoreboard struct {
	TargetId               int
	Timing                 *Timing
	TotalDropoffs          int64
	WaitTimeSummaryMap     map[string]*WaitTimeSummary
	Scorecard              *Scorecard
	WaitTimeChannel        chan *WaitTime
	OperationResultChannel chan *OperationResult
}

func NewScoreboard(targetId int, timing *Timing) *Scoreboard {
	operationResultChannel := make(chan *OperationResult)
	waitTimeChannel := make(chan *WaitTime)
	waitTimeSummaryMap := make(map[string]*WaitTimeSummary)
	scorecard := NewScorecard(targetId, timing.SteadyStateDuration())
	return &Scoreboard{targetId, timing, 0, waitTimeSummaryMap, scorecard, waitTimeChannel, operationResultChannel}
}

func (scoreboard *Scoreboard) Run(quit chan bool) {
	for {
		select {
		case operationResult := <-scoreboard.OperationResultChannel:
			scoreboard.ModifyOperationResult(operationResult)
		case waitTime := <-scoreboard.WaitTimeChannel:
			scoreboard.ModifyWaitTime(waitTime)
		case <-quit:
			quit <- true
			return
		}
	}
}

func (scoreboard *Scoreboard) ModifyOperationResult(operationResult *OperationResult) {
	if scoreboard.Timing.InRampUp(operationResult.TimeStarted) {
		operationResult.TraceLabel = RAMP_UP_TRACE_LABEL
	} else if scoreboard.Timing.InSteadyState(operationResult.TimeFinished) {
		operationResult.TraceLabel = STEADY_STATE_TRACE_LABEL
	} else if scoreboard.Timing.InSteadyState(operationResult.TimeStarted) {
		operationResult.TraceLabel = LATE_TRACE_LABEL
	} else if scoreboard.Timing.InRampDown(operationResult.TimeStarted) {
		operationResult.TraceLabel = RAMP_DOWN_TRACE_LABEL
	}

	//	scoreboard.TotalDropOffWaitTime += dropOffWaitTime
	scoreboard.TotalDropoffs++
	//	if dropOffWaitTime > scoreboard.MaxDropOffWaitTime {
	//		scoreboard.MaxDropOffWaitTime = dropOffWaitTime
	//	}

	if operationResult.TraceLabel == STEADY_STATE_TRACE_LABEL {
		scoreboard.ProcessSteadyStateResult(operationResult)
	} else if operationResult.TraceLabel == LATE_TRACE_LABEL {
		scoreboard.ProcessLateStateResult(operationResult)
	}

}

func (scoreboard *Scoreboard) ModifyWaitTime(waitTime *WaitTime) {
	if !scoreboard.Timing.InSteadyState(waitTime.Time) {
		return
	}

	waitTimeSummary, exists := scoreboard.WaitTimeSummaryMap[waitTime.OperationName]
	if !exists {
		sampler := NewMetricSamplerDummy()
		waitTimeSummary = NewWaitTimeSummary(sampler)
		go waitTimeSummary.Run(make(chan bool))
		scoreboard.WaitTimeSummaryMap[waitTime.OperationName] = waitTimeSummary
	}
	waitTimeSummary.WaitTimeChannel <- waitTime
}

func (scoreboard *Scoreboard) ProcessLateStateResult(operationResult *OperationResult) {
	scoreboard.Scorecard.processLateResult(operationResult)
}

func (scoreboard *Scoreboard) ProcessSteadyStateResult(operationResult *OperationResult) {
	scoreboard.Scorecard.processResult(operationResult)

	if !operationResult.Failed {
		scoreboard.IssueMetricSnapshot(operationResult)
	}
}

func (scoreboard *Scoreboard) IssueMetricSnapshot(operationResult *OperationResult) {
	// write to sonar
}

func (scoreboard *Scoreboard) GetScorboardStatistics() ScoreboardStatistics {
	//	var averageDropOffQTime float64
	//	if scoreboard.TotalDropoffs > 0 {
	//		averageDropOffQTime = float64(scoreboard.TotalDropOffWaitTime) / float64(scoreboard.TotalDropoffs)
	//	}

	stats := ScoreboardStatistics{}
	stats.TargetId = scoreboard.TargetId
	stats.RunDuration = scoreboard.Timing.SteadyStateDuration()
	stats.StartTime = scoreboard.Timing.StartSteadyState
	stats.EndTime = scoreboard.Timing.EndSteadyState
	//	stats.TotalDropoffWaitTime = scoreboard.TotalDropOffWaitTime
	stats.TotalDropoffs = scoreboard.TotalDropoffs
	//	stats.AverageDropOffQTime = averageDropOffQTime
	//	stats.MaxDropOffQTime = scoreboard.MaxDropOffWaitTime
	stats.FinalScorecard = scoreboard.Scorecard.GetScorecardStatistics(scoreboard.Timing.SteadyStateDuration())
	stats.WaitsStats = scoreboard.GetWaitTimeStatistics()

	return stats
}

type ScoreboardStatistics struct {
	TargetId    int   `json:"target_id"`
	RunDuration int64 `json:"run_duration"`
	StartTime   int64 `json:"start_time"`
	EndTime     int64 `json:"end_time"`
	//	TotalDropoffWaitTime int64               `json:"total_dropoff_wait_time"`
	TotalDropoffs int64 `json:"totap_dropoffs"`
	//	AverageDropOffQTime float64             `json:"average_drop_off_q_time"`
	//	MaxDropOffQTime     int64               `json:"max_drop_off_q_time"`
	FinalScorecard ScorecardStatistics `json:"final_scorecard"`
	WaitsStats     WaitsStatistics     `json:"waits_stats"`
}

func (scoreboard *Scoreboard) GetWaitTimeStatistics() WaitsStatistics {
	waits := WaitsStatistics{}
	waits.Waits = make([]WaitTimeSummaryStatistics, len(scoreboard.Scorecard.operationSummaryMap))
	i := 0
	for operationName, _ := range scoreboard.Scorecard.operationSummaryMap {
		waitTimeSummaryStatistic := scoreboard.WaitTimeSummaryMap[operationName].GetStatistics()
		waitTimeSummaryStatistic.OperationName = operationName
		waits.Waits[i] = waitTimeSummaryStatistic
		i++
	}
	return waits
}

type WaitsStatistics struct {
	Waits []WaitTimeSummaryStatistics `json:"waits"`
}
