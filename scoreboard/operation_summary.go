package scoreboard

import (
	"math"
)

const TO_NANO = 1000000000
const INT64_MAX_VALUE = 9223372036854775807
const INT64_MIN_VALUE = -9223372036854775808
const RTIME_T = 3000000000
const meanResponseTimeSamplingInterval = 30

type OperationSummary struct {
	ResponseTimeSampler MetricSampler

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

func NewOperationSummary(metricSampler MetricSampler) *OperationSummary {
	os := &OperationSummary{}
	os.ResponseTimeSampler = metricSampler
	os.MinResponseTime = INT64_MAX_VALUE
	os.MaxResponseTime = INT64_MIN_VALUE
	return os
}

func (os *OperationSummary) ResetSamples() {
	os.ResponseTimeSampler.Reset()
}

func (os *OperationSummary) processOperationSummary(operationResult *OperationResult) {
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

	operationSummary.ResponseTimeSampler.Merge(source.ResponseTimeSampler)
}

func (operationSummary *OperationSummary) GetOperationSummaryStatistics(duration int64) OperationSummaryStatistics {
	var effectiveLoadOperations float64
	var effectiveLoadRequests float64
	var averageRTime float64

	if duration > 0 {
		effectiveLoadOperations = float64(operationSummary.OpsSuccessful) / (float64(duration) / TO_NANO)
		effectiveLoadRequests = float64(operationSummary.ActionsSuccessful) / (float64(duration) / TO_NANO)
	}

	if operationSummary.OpsSuccessful > 0 {
		averageRTime = float64(operationSummary.TotalResponseTime) / float64(operationSummary.OpsSuccessful)
	}

	stats := OperationSummaryStatistics{}
	stats.OpsSuccessful = operationSummary.OpsSuccessful
	stats.OpsFailed = operationSummary.OpsFailed
	stats.OpsSeen = operationSummary.Ops
	stats.ActionsSuccessful = operationSummary.ActionsSuccessful
	stats.Ops = operationSummary.Ops

	stats.EffectiveLoadOps = effectiveLoadOperations
	stats.EffectiveLoadReq = effectiveLoadRequests

	stats.RtimeTotal = operationSummary.TotalResponseTime
	stats.RtimeThrFailed = operationSummary.OpsFailedRtimeThreshold
	stats.RtimeAverage = operationSummary.nNaN(averageRTime)
	stats.RtimeMax = operationSummary.MaxResponseTime
	stats.RtimeMin = operationSummary.MinResponseTime

	stats.SamplerSamplesCollected = operationSummary.ResponseTimeSampler.GetSamplesCollected()
	stats.SamplerSamplesSeen = operationSummary.ResponseTimeSampler.GetSamplesSeen()
	stats.SamplerRtime50th = operationSummary.nNaN(float64(operationSummary.ResponseTimeSampler.GetNthPercentile(50)))
	stats.SamplerRtime90th = operationSummary.nNaN(float64(operationSummary.ResponseTimeSampler.GetNthPercentile(90)))
	stats.SamplerRtime95th = operationSummary.nNaN(float64(operationSummary.ResponseTimeSampler.GetNthPercentile(95)))
	stats.SamplerRtime99th = operationSummary.nNaN(float64(operationSummary.ResponseTimeSampler.GetNthPercentile(99)))
	stats.SamplerRtimeMean = operationSummary.nNaN(float64(operationSummary.ResponseTimeSampler.GetSampleMean()))
	stats.SamplerRtimeStdev = operationSummary.nNaN(float64(operationSummary.ResponseTimeSampler.GetSampleStandardDeviation()))
	stats.SamplerRtimeTvalue = operationSummary.nNaN(float64(operationSummary.ResponseTimeSampler.GetTvalue(averageRTime)))

	return stats
}

func (operationSummary *OperationSummary) nNaN(value float64) float64 {
	if math.IsNaN(value) {
		return 0
	} else if math.IsInf(value, 0) {
		return 0
	}
	return value
}
