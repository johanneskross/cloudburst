package app

import (
	"math/rand"
	"time"
)

type Generator struct {
	thinkTime, cycleTime int64
}

func NewGenerator () *Generator {
 return &Generator{2, 2}
}

func (g *Generator) NextOperation() Operate {
	r:= rand.New(rand.NewSource(time.Now().Unix()))
	i := r.Intn(2)
	return GetOperation(i)
}
