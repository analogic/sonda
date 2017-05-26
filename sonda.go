package main

import (
	//"fmt"
	"runtime"
	//"time"
	"sonda/src"
	"time"
	"fmt"
	"io/ioutil"
	"strconv"
)

var speedPulsesCounter int
var directionPulsesCounter int
var direction int

var speeds []float32
var directions []int

func main() {
	runtime.GOMAXPROCS(4)

	fmt.Println("Starting web server")

	webServer := sonda.WebServer{Port: 7777}
	go webServer.Init()

	fmt.Println("GPIO init")

	gpio := sonda.GPIO{SpeedPin: 17, DirectionPin: 25}
	gpio.Init()
	//defer gpio.Stop()

	filteredPulsesByTimes := make(chan sonda.Pulse)
	filteredPulsesByLogic := make(chan sonda.Pulse)

	go sonda.FilterPulsesByTimes(gpio.Channel, filteredPulsesByTimes)
	go sonda.FilterPulsesByLogic(filteredPulsesByTimes, filteredPulsesByLogic)

	go printResults(&webServer)

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
			newDirection := ((directionPulsesCounter * 10) + 70) % 360

			if (newDirection == 80 || newDirection == 60 || newDirection == 70) && direction > 140 {
				// do nothing, anemometr bug
			} else {
				direction = newDirection
			}
		}
	}
}

func printResults(w *sonda.WebServer) {
	counter := 0
	for {
		time.Sleep(time.Second * 1)

		speed := (float32(speedPulsesCounter) * (float32(30) / float32(1500)))
		fmt.Printf("\n\033[1;34m%vm/s, %v°\033[0m\n", speed, direction)

		speeds = append(speeds, speed)
		directions = append(directions, direction)

		w.WebSocket <- fmt.Sprintf("{\"direction_current\": %v, \"speed_current\": %v}", direction, speed)

		speedPulsesCounter = 0
		directionPulsesCounter = 0

		if(counter == 60) {
			printAverages(w)
			counter = 0
		}
		counter++
	}
}

func printAverages(w *sonda.WebServer) {
	tcpuRaw, _ := ioutil.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	tcpu, _ := strconv.ParseInt(string(tcpuRaw), 10, 64)
	loadRaw, _ := ioutil.ReadFile("/proc/loadavg")
	uptimeRaw, _ := ioutil.ReadFile("/proc/uptime")

	w.DataJson = fmt.Sprintf("{\"speed_average\": %v, \"speed_max\": %v, \"direction_average\": %v, \"temperature_cpu\": %v, \"temperature_gpu\": %v, \"load\": %v, \"uptime\": %v}",
		sonda.AverageSpeed(&speeds),
		sonda.MaxSpeed(&speeds),
		sonda.AverageDirection(&directions),
		tcpu,
		0,
		string(loadRaw),
		string(uptimeRaw))
	fmt.Print(w.DataJson)
}
