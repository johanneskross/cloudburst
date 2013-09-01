package cloudburst

import ()

type Benchmark struct {
}

func NewBenchmark() *Benchmark {
	return &Benchmark{}
}

func (benchmark *Benchmark) Start(targetSchedule *TargetSchedule) {
	// create scenario
	scenario := NewScenario(targetSchedule)

	// launch scenario
	scenario.Launch()

	// aggregateStatistics
	scenario.AggregateStatistics()
}
