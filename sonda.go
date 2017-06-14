package main

import (
	//"fmt"
//	"runtime"
	//"time"
	"sonda/src"
	"time"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	_ "net/http/pprof"
	"syscall"
	"os/exec"
	"log"
)

var speedPulsesCounter int
var speedPulsesOnlyDelay int
var directionPulsesCounter int

var direction int
var speed float32

var speeds []float32
var directions []int

func main() {
//	runtime.GOMAXPROCS(4)

	fmt.Println("Starting web server")

	webServer := sonda.WebServer{Port: 8080}
	go webServer.Init()

	fmt.Println("GPIO init")

	gpio := sonda.GPIO{SpeedPin: 17, DirectionPin: 25}
	gpio.Init()
	//defer gpio.Stop()

	//filteredPulsesByTimes := make(chan sonda.Pulse)
	//filteredPulsesByLogic := make(chan sonda.Pulse)


	//go sonda.FilterPulsesByTimes(gpio.Channel, filteredPulsesByTimes)
	//go sonda.FilterPulsesByLogic(filteredPulsesByTimes, filteredPulsesByLogic)

	go printResults(&webServer)

	speedPulsesCounter = 0
	directionPulsesCounter = 0

	speedPulsesOnlyDelay = 0

	previous := sonda.Pulse{Long: false}
	//for p := range filteredPulsesByLogic {
	for p := range gpio.Channel {

		//if p.Invalid {
		//	continue
		//}

		if p.Long {
			if(speedPulsesOnlyDelay > 1 || speedPulsesOnlyDelay > 35) {
				// we have got start
				direction = ((directionPulsesCounter * 10) + 70 + 180) % 360
				fmt.Printf(" \t\t\033[1;34m(s: %v, l: %v)direction %v°, speed %vm/s\033[0m\n\n", speedPulsesCounter, directionPulsesCounter, direction, speed)
				directionPulsesCounter = 0
			}
			speedPulsesOnlyDelay = 0

			directionPulsesCounter++
		} else {
			speedPulsesCounter++
			speedPulsesOnlyDelay++

			if previous.Long != true {
				fmt.Print(" ")
			}
		}

		fmt.Print(p.String());

		/*
		if speedPulsesCounter == 36 {
			newDirection := ((directionPulsesCounter * 10) + 70 + 180) % 360

			////if (newDirection == 80 || newDirection == 60 || newDirection == 70) && direction > 140 {
				// do nothing, anemometr bug
			////} else {
				direction = newDirection
			///}
		}*/

		previous = p
	}
}

func printResults(w *sonda.WebServer) {
	counter := 0
	for {
		time.Sleep(time.Second * 1)

		speed = (float32(speedPulsesCounter) * (float32(30) / float32(1500)))
		//fmt.Printf("\n\033[1;34m%vm/s, %v°\033[0m\n", speed, direction)

		speeds = append(speeds, speed)
		directions = append(directions, direction)

		w.WebSocket <- fmt.Sprintf("{\"direction_current\": %v, \"speed_current\": %v}", direction, speed)

		speedPulsesCounter = 0
		//directionPulsesCounter = 0

		if(counter == 60) {
			printAverages(w)
			counter = 0
		}
		counter++
	}
}

func printAverages(w *sonda.WebServer) {
	w.DataJson = fmt.Sprintf("{\"speed_average\": %v, \"speed_max\": %v, \"direction_average\": %v, \"temperature_cpu\": %v, \"temperature_gpu\": %v, \"load\": %v, \"uptime\": %v}",
		sonda.AverageSpeed(&speeds),
		sonda.MaxSpeed(&speeds),
		sonda.AverageDirection(&directions),
		raspiCpuTemp(),
		0,
		raspiLoad(),
		raspiUptime())
	fmt.Print(w.DataJson)

	speeds = []float32{}
	directions = []int{}
}

func raspiLoad() float64 {
	line, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return 0.0
	}

	fields := strings.Fields(string(line))
	one, _ := strconv.ParseFloat(fields[0], 32)

	return one
}

func raspiUptime() int64 {
	sysinfo := syscall.Sysinfo_t{}

	if err := syscall.Sysinfo(&sysinfo); err != nil {
		return 0.0
	}

	return int64(sysinfo.Uptime)
}

func raspiCpuTemp() float64 {
	out, err := exec.Command("/opt/vc/bin/vcgencmd", "measure_temp").Output()
	if err != nil {
		log.Fatal(err)
	}

	raw, _ := strconv.ParseFloat(string(out)[5:9], 32)
	return raw;
}
