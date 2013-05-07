package main

import (
	"./app"
	"./app/times"
)

func main() {
	timeSeries := *times.LoadTimeSeries()
	target := *app.NewTarget("target", timeSeries)
	target.RunTimeSeries()
}
