package scoreboard

import (
	"container/list"
)

const SIZE = 1000

type MetricSamplerAll struct {
	Index         int
	Samples       *list.List
	InvalidBuffer bool
	Buffer        []float64
}

func NewMetricSamplerAll() *MetricSamplerAll {
	metricSamplerAll := &MetricSamplerAll{}
	metricSamplerAll.Reset()
	return metricSamplerAll
}

func (metricSampler *MetricSamplerAll) Reset() {
	metricSampler.InvalidBuffer = true
	metricSampler.Samples = list.New()
	metricSampler.Index = 0
	metricSampler.NewBucket()
}

func (metricSampler *MetricSamplerAll) NewBucket() {
	bucket := make([]float64, SIZE)
	metricSampler.Samples.PushFront(bucket)
	metricSampler.Index = 0
}

func (metricSampler *MetricSamplerAll) Accept(observation int64) bool {
	bucket := metricSampler.Samples.Front().Value.([]float64)
	bucket[metricSampler.Index] = float64(observation)
	metricSampler.Index++
	metricSampler.InvalidBuffer = true
	if metricSampler.Index >= SIZE {
		metricSampler.NewBucket()
	}
	return true
}

func GetElementOfList(list *list.List, index int) *list.Element {
	i := 0
	for elem := list.Front(); elem != nil; elem = elem.Next() {
		if i == index {
			return elem
		}
		i++
	}
	return list.Front()
}

func (metricSampler *MetricSamplerAll) UpdateBuffer() {
	if !metricSampler.InvalidBuffer {
		return
	}

	metricSampler.Buffer = make([]float64, SIZE*metricSampler.Samples.Len())
	for i := 0; i < metricSampler.Samples.Len(); i++ {
		bucket := GetElementOfList(metricSampler.Samples, i).Value.([]float64)
		copy(metricSampler.Buffer[i*SIZE:], bucket[0:SIZE])
	}
}

func (metricSampler *MetricSamplerAll) GetSamplesSeen() int {
	return metricSampler.GetSamplesCollected()
}

func (metricSampler *MetricSamplerAll) GetSamplesCollected() int {
	return (metricSampler.Samples.Len()-1)*SIZE + metricSampler.Index
}

func (metricSampler *MetricSamplerAll) GetNthPercentile(pct int) int64 {
	metricSampler.UpdateBuffer()
	return calcPercentile(metricSampler.Buffer, 0, metricSampler.GetSamplesCollected(), pct)
}

func (metricSampler *MetricSamplerAll) GetSampleMean() float64 {
	metricSampler.UpdateBuffer()
	return calcMean(metricSampler.Buffer, 0, metricSampler.GetSamplesCollected())
}

func (metricSampler *MetricSamplerAll) GetSampleStandardDeviation() float64 {
	metricSampler.UpdateBuffer()
	return calcStandardDeviation(metricSampler.Buffer, 0, metricSampler.GetSamplesCollected())
}

func (metricSampler *MetricSamplerAll) GetTvalue(populationMean float64) float64 {
	return 0.0
}

func (metricSampler *MetricSamplerAll) GetRawSamples() *list.List {
	data := list.New()
	for i := 0; i < metricSampler.Samples.Len(); i++ {
		bucket := GetElementOfList(metricSampler.Samples, i).Value.([]float64)

		max := SIZE
		if i == (metricSampler.Samples.Len() - 1) {
			max = metricSampler.Index
		}

		for j := 0; j < max; j++ {
			data.PushBack(float64(bucket[j]))
		}
	}
	return data
}

func (metricSampler *MetricSamplerAll) Merge(responseTimeSampler MetricSampler) {
	for elem := responseTimeSampler.GetRawSamples().Front(); elem != nil; elem = elem.Next() {
		sample := elem.Value.(float64)
		metricSampler.Accept(int64(sample))
	}
}
