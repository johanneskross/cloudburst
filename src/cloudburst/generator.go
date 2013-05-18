package cloudburst

import ()

var factory Factory

type Factory interface {
	CreateGenerator() *Generator
}

type Generator interface {
	GetThinkTime() int
	GetCycleTime() int
	GetOperation()
	NextRequest() *Operation
}
