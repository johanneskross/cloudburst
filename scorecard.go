package cloudburst

import ()

type Scorecard struct {
	TargetId                                          int
	IntervalDuration, TotalOpsInitiated, TotalOpsLate int64
	operationSummary                                  *OperationSummary
	operationSummaryMap                               map[string]*OperationSummary
}

func NewScorecard(targetId int, timeActive int64) *Scorecard {
	operationSummary := NewOperationSummary(NewMetricSamplerPoisson())
	operationSummaryMap := make(map[string]*OperationSummary)
	return &Scorecard{targetId, timeActive, 0, 0, operationSummary, operationSummaryMap}
}

func (scorecard *Scorecard) processLateResult(operationResult *OperationResult) {
	scorecard.TotalOpsInitiated++
	scorecard.TotalOpsLate++
}

func (scorecard *Scorecard) processResult(operationResult *OperationResult) {
	operationSummary, exists := scorecard.operationSummaryMap[operationResult.OperationName]
	if !exists {
		operationSummary = NewOperationSummary(NewMetricSamplerPoisson())
		scorecard.operationSummaryMap[operationResult.OperationName] = operationSummary
	}
	operationSummary.processOperationSummary(operationResult)
	scorecard.operationSummary.processOperationSummary(operationResult)
	scorecard.TotalOpsInitiated++
}

func (scorecard *Scorecard) merge(source *Scorecard) {
	scorecard.TotalOpsInitiated += source.TotalOpsInitiated
	scorecard.TotalOpsLate += source.TotalOpsLate
	scorecard.operationSummary.Merge(source.operationSummary)

	for operationName, sourceOperationSummary := range source.operationSummaryMap {

		operationSummary, exists := scorecard.operationSummaryMap[operationName]
		if !exists {
			operationSummary = NewOperationSummary(NewMetricSamplerAll())
		}

		operationSummary.Merge(sourceOperationSummary)
		scorecard.operationSummaryMap[operationName] = operationSummary
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
	stats.Summary = scorecard.operationSummary.GetOperationSummaryStatistics(duration)
	stats.Operational = scorecard.GetScorecardOperationStatistics(duration)

	return stats
}

func (scorecard *Scorecard) GetScorecardOperationStatistics(duration int64) ScorecardOperationStatistics {
	stats := ScorecardOperationStatistics{}
	stats.Operations = make(map[string]OperationSummaryStatistics)
	for operationName, operationSummary := range scorecard.operationSummaryMap {
		stats.Operations[operationName] = operationSummary.GetOperationSummaryStatistics(duration)
	}
	return stats
}

type ScorecardStatistics struct {
	// aggreation_identifier
	RunDuration       int64                        `json:"run_duration"`
	IntervalDuration  int64                        `json:"interval_duration"`
	TotalOpsInitiated int64                        `json:"total_ops_initiated"`
	TotalOpsLate      int64                        `json:"total_ops_late"`
	OfferedLoadOps    float64                      `json:"offered_load_ops"`
	Summary           OperationSummaryStatistics   `json:"summary"`
	Operational       ScorecardOperationStatistics `json:"operational"`
}

type ScorecardOperationStatistics struct {
	Operations map[string]OperationSummaryStatistics `json:"operations"`
}
