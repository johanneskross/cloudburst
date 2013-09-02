package cloudburst

import (
	"encoding/json"
	s "github.com/johanneskross/cloudburst/scoreboard"
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
	// use all available cpu cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	// start target manager
	targetManagerJoinChannel := make(chan bool)
	go scenario.TargetManager.processSchedule(targetManagerJoinChannel)
	<-targetManagerJoinChannel
}

func (scenario *Scenario) AggregateStatistics() map[string][]byte {
	results := make(map[string][]byte)
	globalScorecard := s.NewScorecard(-1, scenario.TargetSchedule.Duration)

	// get scoreboard for each target
	for elem := scenario.TargetManager.Targets.Front(); elem != nil; elem = elem.Next() {
		target := elem.Value.(*Target)
		scoreboard := target.Scoreboard
		scorecard := scoreboard.Scorecard
		globalScorecard.Merge(scorecard)

		stats, err := json.Marshal(scoreboard.GetScorboardStatistics())
		if err == nil {
			description := "target" + strconv.Itoa(target.TargetId)
			results[description] = stats
		}
	}

	// get global scoreboard for all targets
	stats, err := json.Marshal(globalScorecard.GetScorecardStatistics(
		scenario.TargetSchedule.Duration))
	if err == nil {
		description := "summary"
		results[description] = stats
	}

	return results
}
