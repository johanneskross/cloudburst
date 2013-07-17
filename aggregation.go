package cloudburst

import (
	"container/list"
	"encoding/json"
	"fmt"
)

func AggregateScoreboards(targets *list.List, duration int64) {

	globalScorecard := NewScorecard(-1, duration)

	for elem := targets.Front(); elem != nil; elem = elem.Next() {
		// TODO DUMP TO SONAR
		target := elem.Value.(*Target)
		scoreboard := target.Scoreboard
		scorecard := scoreboard.Scorecard
		globalScorecard.merge(scorecard)

		stats, _ := json.Marshal(scoreboard.GetScorboardStatistics())
		fmt.Println(string(stats))
		fmt.Println("---")

		//TODO ONLY ONE GENERATOR ?
	}

	stats, _ := json.Marshal(globalScorecard.GetScorecardStatistics(duration))
	fmt.Println("----------")
	fmt.Println(string(stats))
}
