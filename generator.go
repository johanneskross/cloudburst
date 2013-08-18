package cloudburst

import ()

type Generator interface {
	GetThinkTime() int64
	GetCycleTime() int64
	NextRequest(url string) Operation
}
