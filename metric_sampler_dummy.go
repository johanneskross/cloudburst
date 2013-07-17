package cloudburst

import (
	"container/list"
)

type MetricSamplerDummy struct {
}

func NewMetricSamplerDummy() *MetricSamplerDummy {
	return &MetricSamplerDummy{}
}

func (metricSamplerDummy *MetricSamplerDummy) Reset() {
}

func (metricSamplerDummy *MetricSamplerDummy) GetSamplesSeen() int {
	return 0
}

func (metricSamplerDummy *MetricSamplerDummy) GetSamplesCollected() int {
	return 0
}

func (metricSamplerDummy *MetricSamplerDummy) Accept(observation int64) bool {
	return false
}

func (metricSamplerDummy *MetricSamplerDummy) GetNthPercentile(pct int) int64 {
	return 0
}

func (metricSamplerDummy *MetricSamplerDummy) GetSampleMean() float64 {
	return 0
}

func (metricSamplerDummy *MetricSamplerDummy) GetSampleStandardDeviation() float64 {
	return 0
}

func (metricSamplerDummy *MetricSamplerDummy) GetTvalue(populationMean float64) float64 {
	return 0.0
}

func (metricSamplerDummy *MetricSamplerDummy) GetRawSamples() *list.List {
	return list.New()
}

func (metricSamplerDummy *MetricSamplerDummy) Merge(responseTimeSampler MetricSampler) {
}
