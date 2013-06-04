package benchmark

import (
	"github.com/johanneskross/cloudburst"
	"math/rand"
	"time"
)

type FactoryImpl struct {
}

func (f FactoryImpl) CreateGenerator() cloudburst.Generator {
	generatorImpl := GeneratorImpl{2, 2}
	generator := cloudburst.Generator(generatorImpl)
	return generator
}

type GeneratorImpl struct {
	ThinkTime, CycleTime int64
}

func NewGeneratorImpl() *GeneratorImpl {
	return &GeneratorImpl{2, 2}
}

func (g GeneratorImpl) GetThinkTime() int64 {
	return g.ThinkTime
}

func (g GeneratorImpl) GetCycleTime() int64 {
	return g.CycleTime
}

func (g GeneratorImpl) NextRequest() cloudburst.Operation {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	i := r.Intn(3)
	return GetOperation(i)
}
