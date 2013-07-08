package cloudburst

import ()

func AggregateScoreboards(targets []*Target, duration int64) {

	globalScorecard := NewScorecard(-1)

	for i := 0; i < len(targets); i++ {
		// TODO DUMP TO SONAR

		scoreboard := targets[i].Scoreboard
		scorecard := scoreboard.Scorecard
		globalScorecard.merge(scorecard)

		//TODO ONLY ONE GENERATOR ?
	}

	globalScorecard.GetStatistics()
}
