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

func (waitTimeSummary *WaitTimeSummary) Run(quit chan bool) {
	for {
		select {
		case waitTime := <-waitTimeSummary.WaitTimeChannel:
			waitTimeSummary.Count++
			waitTimeSummary.TotalWaitTime += waitTime.WaitTime
			if waitTime.WaitTime > waitTimeSummary.MaxWaitTime {
				waitTimeSummary.MaxWaitTime = waitTime.WaitTime
			}
			if waitTime.WaitTime < waitTimeSummary.MinWaitTime {
				waitTimeSummary.MinWaitTime = waitTime.WaitTime
			}
			waitTimeSummary.WaitTimeSampler.Accept(waitTime.WaitTime)
		case <-quit:
			quit <- true
			return
		}
	}
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
