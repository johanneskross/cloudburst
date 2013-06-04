package benchmark

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	TargetFactory  TargetFactory
	TargetSchedule []TargetSchedule
}

type TargetFactory struct {
	TargetFactoryClass  string
	TargetFactoryParams TargetFactoryParamsType
}

type TargetFactoryParamsType struct {
	TimesHost string
	Port      int
	Benchmark BenchmarkType
}

type BenchmarkType struct {
	Txrate, Maxitemsperloc, Itemspertxrate, Parallelism int
	Audit, Stopifauditfailed                            bool
	Plannedlineborrowpercent                            int
}

type TargetSchedule struct {
	Delay, Rampup, Duration, Rampdown int
	TargetFactory                     string
	Workloadprofile                   int
	WorkloadProfileName               string
	Workloadprofileoffset             int
}

func GetConfig(path string) *Config {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	config := &Config{}
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func SetBenchmarkVariables(config BenchmarkType) {
	txRate = config.Txrate
	dbSize = txRate * CUSTOMERS_PER_SCALE
	maxItemsPerLoc = config.Maxitemsperloc
	itemsPerTxRate = config.Itemspertxrate
	parallelism = config.Parallelism
	audit = config.Audit
	stopIfAuditFailed = config.Stopifauditfailed
	plannedLineBorrowPercent = config.Plannedlineborrowpercent
	numOfItems = txRate * itemsPerTxRate
	numItemsPerLoc = Min(maxItemsPerLoc, numOfItems)
}
