package cloudburst

import (
	"fmt"
)

const INT64_MAX_VALUE = 9223372036854775807
const INT64_MIN_VALUE = -9223372036854775808
const RTIME_T = 3000

type OperationSummary struct {
	ResponseTimeSampler PoissonSampling

	Merged            bool
	OpsSuccessful     int64
	OpsFailed         int64
	ActionsSuccessful int64
	Ops               int64

	MinResponseTime int64
	MaxResponseTime int64

	TotalResponseTime       int64
	OpsFailedRtimeThreshold int64
}

func NewOperationSummary(sampling PoissonSampling) *OperationSummary {
	os := &OperationSummary{}
	os.MinResponseTime = INT64_MAX_VALUE
	os.MaxResponseTime = INT64_MIN_VALUE
	return os
}

func (os *OperationSummary) processOperationSummary(operationResult OperationResult) {
	if operationResult.Failed {
		os.OpsFailed++
	} else {
		os.OpsSuccessful++
		os.ActionsSuccessful += operationResult.ActionsPerformed
		os.Ops++

		responseTime := operationResult.GetExecutionTime()
		os.ResponseTimeSampler.Accept(responseTime)
		os.TotalResponseTime += responseTime
		if responseTime > RTIME_T {
			os.OpsFailedRtimeThreshold++
		}

		// psquare

		if os.MaxResponseTime < responseTime {
			os.MaxResponseTime = responseTime
		}
		if os.MinResponseTime > responseTime {
			os.MinResponseTime = responseTime
		}
	}
}

func (operationSummary *OperationSummary) Merge(source *OperationSummary) {
	operationSummary.Merged = true

	operationSummary.Ops += source.Ops
	operationSummary.OpsSuccessful += source.OpsSuccessful
	operationSummary.OpsFailed += source.OpsFailed
	operationSummary.ActionsSuccessful += source.ActionsSuccessful

	if source.MinResponseTime < operationSummary.MinResponseTime {
		operationSummary.MinResponseTime = source.MinResponseTime
	}
	if source.MaxResponseTime > operationSummary.MaxResponseTime {
		operationSummary.MaxResponseTime = source.MaxResponseTime
	}

	operationSummary.TotalResponseTime += source.TotalResponseTime
	operationSummary.OpsFailedRtimeThreshold += source.OpsFailedRtimeThreshold

	// TODO merge time sampling
}

func (operationSummary *OperationSummary) GetStatistics() {
	fmt.Print("\n\tMerged: ")
	fmt.Println(operationSummary.Merged)
	fmt.Print("\tOps: ")
	fmt.Println(operationSummary.Ops)
	fmt.Print("\tOpsSuccessful: ")
	fmt.Println(operationSummary.OpsSuccessful)
	fmt.Print("\tOpsFailed: ")
	fmt.Println(operationSummary.OpsFailed)
	fmt.Print("\tActionsSuccessful: ")
	fmt.Println(operationSummary.ActionsSuccessful)
	fmt.Print("\tMinResponseTime[s]: ")
	fmt.Println(operationSummary.MinResponseTime / TO_NANO)
	fmt.Print("\tMaxResponseTime[s]: ")
	fmt.Println(operationSummary.MaxResponseTime / TO_NANO)
	fmt.Print("\tTotalResponseTime[s]: ")
	fmt.Println(operationSummary.TotalResponseTime / TO_NANO)
	fmt.Print("\tOpsFailedRtimeThreshold: ")
	fmt.Println(operationSummary.OpsFailedRtimeThreshold)

	//fmt.Printf("Struct: %+v", operationSummary)
}
