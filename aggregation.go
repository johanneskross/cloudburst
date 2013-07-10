package cloudburst

import (
	"container/list"
)

func AggregateScoreboards(targets *list.List, duration int64) {

	globalScorecard := NewScorecard(-1)

	for elem := targets.Front(); elem != nil; elem = elem.Next() {
		// TODO DUMP TO SONAR
		target := elem.Value.(*Target)
		scoreboard := target.Scoreboard
		scorecard := scoreboard.Scorecard
		globalScorecard.merge(scorecard)

		//TODO ONLY ONE GENERATOR ?
	}

	globalScorecard.GetStatistics()
}
