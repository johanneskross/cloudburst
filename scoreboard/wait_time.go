package scoreboard

import ()

type WaitTime struct {
	Time, WaitTime int64
	OperationName  string
}

func NewWaitTime(time, waitTime int64, operationName string) *WaitTime {
	return &WaitTime{time, waitTime, operationName}
}

type WaitTimeSummary struct {
	Count, TotalWaitTime, MinWaitTime, MaxWaitTime int64
	WaitTimeSampler                                MetricSampler
	WaitTimeChannel                                chan *WaitTime
}

func NewWaitTimeSummary(waitTimeSampler MetricSampler) *WaitTimeSummary {
	return &WaitTimeSummary{0, 0, INT64_MAX_VALUE, INT64_MIN_VALUE, waitTimeSampler, make(chan *WaitTime)}
}

func (waitTimeSummary *WaitTimeSummary) Receive(waitTime int64) {
	waitTimeSummary.Count++
	waitTimeSummary.TotalWaitTime += waitTime
	if waitTime > waitTimeSummary.MaxWaitTime {
		waitTimeSummary.MaxWaitTime = waitTime
	}
	if waitTime < waitTimeSummary.MinWaitTime {
		waitTimeSummary.MinWaitTime = waitTime
	}
	waitTimeSummary.WaitTimeSampler.Accept(waitTime)

}

func (waitTimeSummary *WaitTimeSummary) GetStatistics() WaitTimeSummaryStatistics {
	minWaitTime := waitTimeSummary.MinWaitTime
	if minWaitTime == INT64_MAX_VALUE {
		minWaitTime = 0
	}

	maxWaitTime := waitTimeSummary.MaxWaitTime
	if maxWaitTime == INT64_MIN_VALUE {
		maxWaitTime = 0
	}

	var avgWaitTime float64
	if waitTimeSummary.Count > 0 {
		avgWaitTime = float64(waitTimeSummary.TotalWaitTime) / float64(waitTimeSummary.Count)
	}

	tvalue := waitTimeSummary.WaitTimeSampler.GetTvalue(avgWaitTime)

	stats := WaitTimeSummaryStatistics{}
	stats.AverageWaitTime = avgWaitTime
	stats.TotalWaitTime = waitTimeSummary.TotalWaitTime
	stats.MinWaitTime = minWaitTime
	stats.MaxWaitTime = maxWaitTime
	stats.PercentileWaitTime90th = waitTimeSummary.WaitTimeSampler.GetNthPercentile(90)
	stats.PercentileWaitTime99th = waitTimeSummary.WaitTimeSampler.GetNthPercentile(90)
	stats.SamplesCollected = waitTimeSummary.WaitTimeSampler.GetSamplesCollected()
	stats.SamplesSeen = waitTimeSummary.WaitTimeSampler.GetSamplesSeen()
	stats.SampleMean = waitTimeSummary.WaitTimeSampler.GetSampleMean()
	stats.SampleStandardDeviation = waitTimeSummary.WaitTimeSampler.GetSampleStandardDeviation()
	stats.TvalueAverageWaitTime = tvalue

	return stats
}
