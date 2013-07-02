package cloudburst

import (
	"math"
	"math/rand"
	"time"
)

const meanSamplingInterval = 30.0

type PoissonSampling struct {
	TargetId                       int
	OperationName                  string
	NextSampleToAccept, SampleSeen int
	Random                         NegativeExponential
}

func NewPoissonSampling(targetId int, operationName string) PoissonSampling {
	poisson := PoissonSampling{}
	poisson.TargetId = targetId
	poisson.OperationName = operationName
	poisson.Random = NewNegativeExponential(meanSamplingInterval)
	poisson.Reset()

	//TODO SONAR RECORDER
	return poisson
}

func (poisson PoissonSampling) Reset() {
	// poisson.Sampling.Reset()
	poisson.SampleSeen = 0
	poisson.NextSampleToAccept = 1
}

func (poisson PoissonSampling) Accept(value int64) bool {
	poisson.SampleSeen++

	if poisson.SampleSeen == poisson.NextSampleToAccept {
		// TODO 

		rand.Seed(time.Now().UTC().UnixNano())
		random := rand.Float64()
		poisson.NextSampleToAccept = poisson.SampleSeen + int(math.Ceil(random))

		// Write to sonar
		return true
	}
	return false
}

type NegativeExponential struct {
	mean float64
}

func NewNegativeExponential(meanSamplingInterval float64) NegativeExponential {
	return NegativeExponential{}
}

func (ne NegativeExponential) NextDouble() float64 {
	if ne.mean == 0 {
		return 0.0
	}

	rand.Seed(time.Now().UTC().UnixNano())
	random := rand.Float64()
	return -1.0 * ne.mean * math.Log(random)
}
