package cloudburst

import "fmt"

type Benchmark struct {
}

func NewBenchmark() *Benchmark {
	return &Benchmark{}
}

func (benchmark *Benchmark) Start(targetSchedule *TargetSchedule) {
	// create scenario
	scenario := NewScenario(targetSchedule)

	// launch scenario
	fmt.Println("Launch Scenario")
	scenario.Launch()

	// aggregateStatistics
	scenario.AggregateStatistics()
	fmt.Println("End Scenario")
}
