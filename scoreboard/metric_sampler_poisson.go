package scoreboard

import (
	"container/list"
	"math"
	"math/rand"
	"time"
)

const MEAN_SAMPLING_INTERVAL = 30

type MetricSamplerPoisson struct {
	Sampling                        MetricSamplerAll
	NextSampleToAccept, SamplesSeen int
	RandGenerator                   *rand.Rand
}

func NewMetricSamplerPoisson() *MetricSamplerPoisson {
	metricSamplerPoisson := &MetricSamplerPoisson{}
	metricSamplerPoisson.Reset()
	randSource := rand.NewSource(time.Now().UnixNano())
	metricSamplerPoisson.RandGenerator = rand.New(randSource)
	return metricSamplerPoisson
}

func (metricSamplerPoisson *MetricSamplerPoisson) Reset() {
	metricSamplerPoisson.Sampling.Reset()
	metricSamplerPoisson.SamplesSeen = 0
	metricSamplerPoisson.NextSampleToAccept = 1
}

func (metricSamplerPoisson *MetricSamplerPoisson) GetSamplesSeen() int {
	return metricSamplerPoisson.SamplesSeen
}

func (metricSamplerPoisson *MetricSamplerPoisson) GetSamplesCollected() int {
	return metricSamplerPoisson.Sampling.GetSamplesCollected()
}

func (metricSamplerPoisson *MetricSamplerPoisson) Accept(observation int64) bool {
	metricSamplerPoisson.SamplesSeen++

	if metricSamplerPoisson.SamplesSeen == metricSamplerPoisson.NextSampleToAccept {
		metricSamplerPoisson.Sampling.Accept(observation)
		randDouble := metricSamplerPoisson.RandGenerator.Float64()
		randExp := -1 * MEAN_SAMPLING_INTERVAL * math.Log(randDouble)
		metricSamplerPoisson.NextSampleToAccept = metricSamplerPoisson.SamplesSeen * int(math.Ceil(randExp))

		return true
	}

	return false
}

func (metricSamplerPoisson *MetricSamplerPoisson) GetNthPercentile(pct int) int64 {
	return metricSamplerPoisson.Sampling.GetNthPercentile(pct)
}

func (metricSamplerPoisson *MetricSamplerPoisson) GetSampleMean() float64 {
	return metricSamplerPoisson.Sampling.GetSampleMean()
}

func (metricSamplerPoisson *MetricSamplerPoisson) GetSampleStandardDeviation() float64 {
	return metricSamplerPoisson.Sampling.GetSampleStandardDeviation()
}

func (metricSamplerPoisson *MetricSamplerPoisson) GetTvalue(populationMean float64) float64 {
	return metricSamplerPoisson.Sampling.GetTvalue(populationMean)
}

func (metricSamplerPoisson *MetricSamplerPoisson) GetRawSamples() *list.List {
	return metricSamplerPoisson.Sampling.GetRawSamples()
}

func (metricSamplerPoisson *MetricSamplerPoisson) Merge(responseTimeSampler MetricSampler) {
	metricSamplerPoisson.Sampling.Merge(responseTimeSampler)
}
