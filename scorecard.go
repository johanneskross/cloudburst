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
	scorecard.TotalOpsInitiated++
}

func (scorecard *Scorecard) merge(source *Scorecard) {
	scorecard.TotalOpsInitiated += source.TotalOpsInitiated
	scorecard.operationSummary.Merge(source.operationSummary)

	for operationName, sourceOperationSummary := range source.operationSummaryMap {

		operationSummary, exists := scorecard.operationSummaryMap[operationName]
		if !exists {
			operationSummary = NewOperationSummary(NewPoissonSampling(-1, operationName)) // TODO not Poisson
		}

		operationSummary.Merge(sourceOperationSummary)
		scorecard.operationSummaryMap[operationName] = operationSummary
	}
}

func (scorecard *Scorecard) GetStatistics() {
	fmt.Println("------ STATISTICS -------")
	fmt.Print("TargetId: \t")
	fmt.Println(scorecard.TargetId)
	fmt.Print("TotalOpsInitiated: \t")
	fmt.Println(scorecard.TotalOpsInitiated)
	fmt.Print("OperationSummary: \t")
	scorecard.operationSummary.GetStatistics()
	fmt.Println("OperationSummaryMap:")
	for key, value := range scorecard.operationSummaryMap {
		fmt.Print("- OperationName: ")
		fmt.Println(key)
		fmt.Print("- OperationSummary: ")
		value.GetStatistics()
	}
	fmt.Println("---- END STATISTICS -----")
}
