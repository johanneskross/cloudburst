package benchmark

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
}

type TargetSchedule struct {
	Delay, Rampup, Duration, Rampdown  int
	TargetFactory, WorkloadProfileName string
}

func GetConfig(path string) *Config {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("File not found: %v\n", err)
	}
	config := &Config{}
	json.Unmarshal(file, &config)
	return config
}
