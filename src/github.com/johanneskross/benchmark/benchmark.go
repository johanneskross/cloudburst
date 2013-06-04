package benchmark

import (
	"github.com/johanneskross/cloudburst"
	"github.com/johanneskross/times"
	"strconv"
)

func RunBenchmark() {
	scenario := cloudburst.NewScenario()
	schedule := *cloudburst.NewTargetSchedule()

	config := GetConfig("./config.json")
	params := config.TargetFactory.TargetFactoryParams
	SetBenchmarkVariables(config.TargetFactory.TargetFactoryParams.Benchmark)
	for i := 0; i < len(config.TargetSchedule); i++ {
		s := config.TargetSchedule[i]
		conf := cloudburst.NewTargetConfiguration(strconv.Itoa(i), s.Delay, s.Rampup, s.Duration, s.Rampdown, *times.LoadTimeSeries(params.TimesHost, params.Port, s.WorkloadProfileName))
		schedule.AddTargetConfiguration(conf)
	}

	factoryImpl := FactoryImpl{}
	factory := cloudburst.Factory(factoryImpl)
	scenario.Launch(schedule, factory)
}
