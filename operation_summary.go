package cloudburst

import ()

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
