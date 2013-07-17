package cloudburst

import (
	"container/list"
	"math"
	"sort"
)

type MetricSampler interface {
	Reset()
	GetSamplesSeen() int
	GetSamplesCollected() int
	Accept(observation int64) bool
	GetNthPercentile(pct int) int64
	GetSampleMean() float64
	GetSampleStandardDeviation() float64
	GetTvalue(populationMean float64) float64
	GetRawSamples() *list.List
	Merge(responseTimeSampler MetricSampler)
}

func calcPercentile(array []float64, startPosition, length, percentile int) int64 {
	if length == 0 {
		return 0
	}

	values := make([]float64, length)
	copy(values, array[startPosition:startPosition+length])
	sort.Float64s(values)

	p := float64(percentile) / 100
	rank := (float64(length) - 1) * p

	if float64(int64(rank)) != rank {
		position1 := int64(rank)
		position2 := position1 + 1

		value1 := values[position1]
		value2 := values[position2]

		ratio := rank - float64(position1)
		result := value1 + value2 - value1*ratio
		return int64(result)
	} else {
		position1 := int64(rank)
		return int64(values[position1])
	}
}

func calcMean(array []float64, startPosition, length int) float64 {
	sum := 0.0
	for _, v := range array[startPosition : startPosition+length] {
		sum += v
	}
	return sum / float64(length)
}

func calcStandardDeviation(array []float64, startPosition, length int) float64 {
	mean := calcMean(array, startPosition, length)
	sum := 0.0
	for _, v := range array[startPosition : startPosition+length] {
		sum += math.Pow(v-mean, 2)
	}
	return math.Sqrt(sum / float64(length))
}
