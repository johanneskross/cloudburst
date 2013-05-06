package app

import (
	"fmt"
	"time"
)

type Agent struct {
	Name, Url			string
}

func NewAgent(name string) *Agent {
	return &Agent{name, "http://www.example.com/"}
}

func (a *Agent) GenerateLoad() {
	pullUrl("http://www.example.com/")
}

func (a *Agent) Run(c, quit chan bool) {
	fmt.Println("yes")
	c <- true
	for {
		select {
			case <- quit:
				<- c
				fmt.Printf("Stopping agent%v ...\n", a.Name)
				return
			default:
				fmt.Printf("Running agent%v ...\n", a.Name)
				time.Sleep(3 * time.Second)
		}
	}
}