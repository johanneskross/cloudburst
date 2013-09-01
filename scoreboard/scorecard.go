package scoreboard

import ()

type Scorecard struct {
	TargetId                                          int
	IntervalDuration, TotalOpsInitiated, TotalOpsLate int64
	OperationSummary                                  *OperationSummary
	OperationSummaryMap                               map[string]*OperationSummary
}

func NewScorecard(targetId int, timeActive int64) *Scorecard {
	operationSummary := NewOperationSummary(NewMetricSamplerPoisson())
	operationSummaryMap := make(map[string]*OperationSummary)
	return &Scorecard{targetId, timeActive, 0, 0, operationSummary, operationSummaryMap}
}

func (scorecard *Scorecard) ProcessLateResult(operationResult *OperationResult) {
	scorecard.TotalOpsInitiated++
	scorecard.TotalOpsLate++
}

func (scorecard *Scorecard) ProcessResult(operationResult *OperationResult) {
	operationSummary, exists := scorecard.OperationSummaryMap[operationResult.OperationName]
	if !exists {
		operationSummary = NewOperationSummary(NewMetricSamplerPoisson())
		scorecard.OperationSummaryMap[operationResult.OperationName] = operationSummary
	}
	operationSummary.processOperationSummary(operationResult)
	scorecard.OperationSummary.processOperationSummary(operationResult)
	scorecard.TotalOpsInitiated++
}

func (scorecard *Scorecard) Merge(source *Scorecard) {
	scorecard.TotalOpsInitiated += source.TotalOpsInitiated
	scorecard.TotalOpsLate += source.TotalOpsLate
	scorecard.OperationSummary.Merge(source.OperationSummary)

	for operationName, sourceOperationSummary := range source.OperationSummaryMap {

		operationSummary, exists := scorecard.OperationSummaryMap[operationName]
		if !exists {
			operationSummary = NewOperationSummary(NewMetricSamplerAll())
		}

		operationSummary.Merge(sourceOperationSummary)
		scorecard.OperationSummaryMap[operationName] = operationSummary
	}
}

func (scorecard *Scorecard) GetScorecardStatistics(duration int64) ScorecardStatistics {
	var offeredLoadOps float64
	if duration > 0 {
		offeredLoadOps = float64(scorecard.TotalOpsInitiated) / (float64(duration) / TO_NANO)
	}

	stats := ScorecardStatistics{}
	stats.RunDuration = duration
	stats.IntervalDuration = scorecard.IntervalDuration
	stats.TotalOpsInitiated = scorecard.TotalOpsInitiated
	stats.TotalOpsLate = scorecard.TotalOpsLate
	stats.OfferedLoadOps = offeredLoadOps
	stats.Summary = scorecard.OperationSummary.GetOperationSummaryStatistics(duration)
	stats.Operational = scorecard.GetScorecardOperationStatistics(duration)

	return stats
}

func (scorecard *Scorecard) GetScorecardOperationStatistics(duration int64) ScorecardOperationStatistics {
	stats := ScorecardOperationStatistics{}
	stats.Operations = make(map[string]OperationSummaryStatistics)
	for operationName, operationSummary := range scorecard.OperationSummaryMap {
		stats.Operations[operationName] = operationSummary.GetOperationSummaryStatistics(duration)
	}
	return stats
}
