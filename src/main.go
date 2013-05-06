package main

import (
	"app"
	"fmt"
)

func main() {
	target1 := *app.NewTarget("target1", 5)
	fmt.Printf("%v with %v agents running\n", target1.Name, len(target1.AgentChannel))
	target1.UpdateAgents()	
	fmt.Printf("%v with %v agents running\n", target1.Name, len(target1.AgentChannel))
	
}
