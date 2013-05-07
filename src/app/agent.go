package app

import (
	"fmt"
	"time"
)

type Agent struct {
	Name, Url	string
	Quit		chan bool
}

func NewAgent(name string, quit chan bool) *Agent {
	return &Agent{name, "http://www.example.com/", quit}
}

func (a *Agent) GenerateLoad() {
	pullUrl("http://www.example.com/")
}

func (a *Agent) Run(c chan bool) {
	c <- true
	for {
		select {
			case <- a.Quit:
				<- c
				fmt.Printf("Stopping agent%v ...\n", a.Name)
				return
			default:
				fmt.Printf("Running agent%v ...\n", a.Name)
				a.GenerateLoad()
				time.Sleep(3 * time.Second)
		}
	}
}