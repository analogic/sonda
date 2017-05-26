package main

import (
	"fmt"
	"runtime"
	"time"
	"sonda/src"
)

var speedPulsesCounter int
var directionPulsesCounter int
var direction int

func main() {
	runtime.GOMAXPROCS(5)

	fmt.Println("GPIO init")

	gpio := sonda.GPIO{SpeedPin: 17, DirectionPin: 25}
	defer gpio.Stop()

	fmt.Println("Starting listening")
	gpio.Init()

	// true - long pulse, false - short pulse

	filteredPulsesByTimes :=  make(chan sonda.Pulse)
	filteredPulsesByLogic :=  make(chan sonda.Pulse)

	go sonda.FilterPulsesByTimes(gpio.Channel, filteredPulsesByTimes)
	go sonda.FilterPulsesByLogic(filteredPulsesByTimes, filteredPulsesByLogic)

	go printResults()

	speedPulsesCounter = 0
	directionPulsesCounter = 0

	for p := range filteredPulsesByLogic {

		fmt.Print(p.String());

		if p.Invalid {
			continue
		}

		if p.Long {
			directionPulsesCounter++
		} else {
			speedPulsesCounter++
		}


		if speedPulsesCounter == 36 {
			direction = ((directionPulsesCounter * 10) + 70) % 360
		}

	}
}

func printResults() {
	for {
		time.Sleep(time.Second * 3)
		fmt.Println("\033[1;34m")

		speed := (float32(speedPulsesCounter) * (float32(30) / float32(1500))) / 3
		fmt.Printf("%v pulses, %v direction", speed, direction)

		fmt.Println("\033[0m")

		speedPulsesCounter = 0
		directionPulsesCounter = 0
	}
}