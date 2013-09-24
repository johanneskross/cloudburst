package cloudburst

import (
	s "github.com/johanneskross/cloudburst/scoreboard"
)

const NO_TRACE_LABEL = "[NONE]"
const STEADY_STATE_TRACE_LABEL = "[STEADY-STATE]"
const LATE_TRACE_LABEL = "[LATE]"
const RAMP_UP_TRACE_LABEL = "[RAMP-UP]"
const RAMP_DOWN_TRACE_LABEL = "[RAMP-DOWN]"

type Scoreboard struct {
	TargetId               int
	Timing                 *Timing
	TotalDropoffs          int64
	WaitTimeSummaryMap     map[string]*s.WaitTimeSummary
	Scorecard              *s.Scorecard
	WaitTimeChannel        chan *s.WaitTime
	OperationResultChannel chan *s.OperationResult
}

func NewScoreboard(targetId, channelSize int, timing *Timing) *Scoreboard {
	operationResultChannel := make(chan *s.OperationResult, channelSize)
	waitTimeChannel := make(chan *s.WaitTime, channelSize)
	waitTimeSummaryMap := make(map[string]*s.WaitTimeSummary)
	scorecard := s.NewScorecard(targetId, timing.SteadyStateDuration())
	return &Scoreboard{targetId, timing, 0, waitTimeSummaryMap, scorecard, waitTimeChannel, operationResultChannel}
}

func (scoreboard *Scoreboard) Run(scoreboardJoinChannel chan bool) {
	for {
		select {
		case operationResult := <-scoreboard.OperationResultChannel:
			scoreboard.ModifyOperationResult(operationResult)
		case <-scoreboardJoinChannel:
			scoreboardJoinChannel <- true
			return
		}
	}
}

func (scoreboard *Scoreboard) RunWaitTime(scoreboardJoinChannel2 chan bool) {
	for {
		select {
		case waitTime := <-scoreboard.WaitTimeChannel:
			scoreboard.ModifyWaitTime(waitTime)
		case <-scoreboardJoinChannel2:
			scoreboardJoinChannel2 <- true
			return
		}
	}
}

func (scoreboard *Scoreboard) ModifyOperationResult(operationResult *s.OperationResult) {
	if scoreboard.Timing.InRampUp(operationResult.TimeStarted) {
		operationResult.TraceLabel = RAMP_UP_TRACE_LABEL
	} else if scoreboard.Timing.InSteadyState(operationResult.TimeFinished) {
		operationResult.TraceLabel = STEADY_STATE_TRACE_LABEL
	} else if scoreboard.Timing.InSteadyState(operationResult.TimeStarted) {
		operationResult.TraceLabel = LATE_TRACE_LABEL
	} else if scoreboard.Timing.InRampDown(operationResult.TimeStarted) {
		operationResult.TraceLabel = RAMP_DOWN_TRACE_LABEL
	}

	scoreboard.TotalDropoffs++

	if operationResult.TraceLabel == STEADY_STATE_TRACE_LABEL {
		scoreboard.ProcessSteadyStateResult(operationResult)
	} else if operationResult.TraceLabel == LATE_TRACE_LABEL {
		scoreboard.ProcessLateStateResult(operationResult)
	}

}

func (scoreboard *Scoreboard) ModifyWaitTime(waitTime *s.WaitTime) {
	if !scoreboard.Timing.InSteadyState(waitTime.Time) {
		return
	}

	waitTimeSummary, exists := scoreboard.WaitTimeSummaryMap[waitTime.OperationName]
	if !exists {
		sampler := s.NewMetricSamplerDummy()
		waitTimeSummary = s.NewWaitTimeSummary(sampler)
		scoreboard.WaitTimeSummaryMap[waitTime.OperationName] = waitTimeSummary
	}
	waitTimeSummary.Receive(waitTime.WaitTime)
}

func (scoreboard *Scoreboard) ProcessLateStateResult(operationResult *s.OperationResult) {
	scoreboard.Scorecard.ProcessLateResult(operationResult)
}

func (scoreboard *Scoreboard) ProcessSteadyStateResult(operationResult *s.OperationResult) {
	scoreboard.Scorecard.ProcessResult(operationResult)
}

func (scoreboard *Scoreboard) GetScorboardStatistics() s.ScoreboardStatistics {
	stats := s.ScoreboardStatistics{}
	stats.TargetId = scoreboard.TargetId
	stats.RunDuration = scoreboard.Timing.SteadyStateDuration()
	stats.StartTime = scoreboard.Timing.StartSteadyState
	stats.EndTime = scoreboard.Timing.EndSteadyState
	stats.TotalDropoffs = scoreboard.TotalDropoffs
	stats.FinalScorecard = scoreboard.Scorecard.GetScorecardStatistics(scoreboard.Timing.SteadyStateDuration())
	stats.WaitTimeSummaries = scoreboard.GetWaitTimeStatistics()
	return stats
}

func (scoreboard *Scoreboard) GetWaitTimeStatistics() s.WaitTimeSummariesStatistics {
	waitTimeSummariesStatistics := s.WaitTimeSummariesStatistics{}
	waitTimeSummariesStatistics.WaitTimeSummaries = make([]s.WaitTimeSummaryStatistics, len(scoreboard.Scorecard.OperationSummaryMap))
	i := 0
	for operationName, _ := range scoreboard.Scorecard.OperationSummaryMap {
		waitTimeSummaryStatistic := scoreboard.WaitTimeSummaryMap[operationName].GetStatistics()
		waitTimeSummaryStatistic.OperationName = operationName
		waitTimeSummariesStatistics.WaitTimeSummaries[i] = waitTimeSummaryStatistic
		i++
	}
	return waitTimeSummariesStatistics
}
