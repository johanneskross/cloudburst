package cloudburst

import (
	"encoding/json"
	s "github.com/johanneskross/cloudburst/scoreboard"
	"io/ioutil"
	"runtime"
	"strconv"
)

type Scenario struct {
	TargetManager  *TargetManager
	TargetSchedule *TargetSchedule
}

func NewScenario(targetSchedule *TargetSchedule) *Scenario {
	return &Scenario{NewTargetManager(targetSchedule), targetSchedule}
}

func (scenario *Scenario) Launch() {
	// Use all available cpu cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	targetManagerJoinChannel := make(chan bool)
	go scenario.TargetManager.processSchedule(targetManagerJoinChannel)
	<-targetManagerJoinChannel
}

func (scenario *Scenario) AggregateStatistics() {
	globalScorecard := s.NewScorecard(-1, scenario.TargetSchedule.Duration)

	for elem := scenario.TargetManager.Targets.Front(); elem != nil; elem = elem.Next() {
		target := elem.Value.(*Target)
		scoreboard := target.Scoreboard
		scorecard := scoreboard.Scorecard
		globalScorecard.Merge(scorecard)

		stats, err := json.Marshal(scoreboard.GetScorboardStatistics())
		if err == nil {
			filename := "target" + strconv.Itoa(target.TargetId) + ".json"
			ioutil.WriteFile(filename, stats, 0644)
		}
	}

	stats, err := json.Marshal(globalScorecard.GetScorecardStatistics(scenario.TargetSchedule.Duration))
	if err == nil {
		filename := "summary.json"
		ioutil.WriteFile(filename, stats, 0644)
	}
}
