package cloudburst

import ()

type Generator interface {
	GetWaitTime() int64
	NextRequest(url string) Operation
}
