package cloudburst

import (
	"fmt"
)

type Scorecard struct {
	TargetId            int
	TotalOpsInitiated   int64
	operationSummary    *OperationSummary
	operationSummaryMap map[string]*OperationSummary
}

func NewScorecard(targetId int) *Scorecard {
	operationSummary := NewOperationSummary(NewPoissonSampling(targetId, "all"))
	operationSummaryMap := make(map[string]*OperationSummary)
	return &Scorecard{targetId, 0, operationSummary, operationSummaryMap}
}

func (scorecard *Scorecard) processResult(operationResult OperationResult) {
	operationSummary, exists := scorecard.operationSummaryMap[operationResult.OperationName]
	if !exists {
		operationSummary = NewOperationSummary(NewPoissonSampling(scorecard.TargetId, operationResult.OperationName))
		scorecard.operationSummaryMap[operationResult.OperationName] = operationSummary
	}
	operationSummary.processOperationSummary(operationResult)
	scorecard.operationSummary.processOperationSummary(operationResult)
	fmt.Println(scorecard.operationSummary)
	scorecard.TotalOpsInitiated++
}
