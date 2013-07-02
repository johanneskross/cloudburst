package cloudburst

import ()

type Factory interface {
	CreateGenerator() Generator
}

type Generator interface {
	GetThinkTime() int64
	GetCycleTime() int64
	NextRequest(url string) Operation
}
