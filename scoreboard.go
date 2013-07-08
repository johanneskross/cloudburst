package cloudburst

import ()

type Scoreboard struct {
	TargetId, TotalDropoffs, TotalDropOffWaitTime, MaxDropOffWaitTime int
	Scorecard                                                         *Scorecard
	OperationResultChannel                                            chan OperationResult
}

func NewScoreboard(targetId int) *Scoreboard {
	osChan := make(chan OperationResult)
	return &Scoreboard{targetId, 0, 0, 0, NewScorecard(targetId), osChan}
}

func (scoreboard *Scoreboard) Run(quit chan bool) {
	for {
		select {
		case operationResult := <-scoreboard.OperationResultChannel:
			scoreboard.Scorecard.processResult(operationResult)
		case <-quit:
			quit <- true
			return
		}
	}
}
