package cloudburst

import (
	"github.com/johanneskross/cloudburst/scoreboard"
)

type Operation interface {
	Run(timing *Timing) *scoreboard.OperationResult
}
