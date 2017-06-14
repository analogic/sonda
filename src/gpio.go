package sonda

import (
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"

	"time"
)

type GPIO struct {
	SpeedPin int
	DirectionPin int
	digitalSpeedPin embd.DigitalPin
	digitalDirectionPin embd.DigitalPin
	Channel chan Pulse
}

func (g *GPIO) Init() {
	g.Channel = make(chan Pulse)
	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}
	//g.initSpeed()
	//g.initDirection()
	g.initTogether()
}
//
//func (g *GPIO) initSpeed() {
//	var err error
//	g.digitalSpeedPin, err = embd.NewDigitalPin(g.SpeedPin)
//	if err != nil {
//		panic(err)
//	}
//
//	if err := g.digitalSpeedPin.SetDirection(embd.In); err != nil {
//		panic(err)
//	}
//	g.digitalSpeedPin.ActiveLow(false)
//
//	err = g.digitalSpeedPin.Watch(embd.EdgeRising, func(speed embd.DigitalPin) {
//		time.Sleep(5 * time.Millisecond) // we need speed pulse come always after direction pulse
//		g.Channel <- Pulse{Long: false, At: time.Now()}
//	})
//	if err != nil {
//		panic(err)
//	}
//}
//
//func (g *GPIO) initDirection() {
//	var err error
//	g.digitalDirectionPin, err = embd.NewDigitalPin(g.DirectionPin)
//	if err != nil {
//		panic(err)
//	}
//
//	if err := g.digitalDirectionPin.SetDirection(embd.In); err != nil {
//		panic(err)
//	}
//	g.digitalDirectionPin.ActiveLow(false)
//
//	err = g.digitalDirectionPin.Watch(embd.EdgeRising, func(direction embd.DigitalPin) {
//		g.Channel <- Pulse{Long: true, At: time.Now()}
//	})
//	if err != nil {
//		panic(err)
//	}
//}

func (g *GPIO) initTogether() {
	var err error
	g.digitalSpeedPin, err = embd.NewDigitalPin(g.SpeedPin)
	if err != nil {
		panic(err)
	}

	if err := g.digitalSpeedPin.SetDirection(embd.In); err != nil {
		panic(err)
	}
	g.digitalSpeedPin.ActiveLow(false)


	g.digitalDirectionPin, err = embd.NewDigitalPin(g.DirectionPin)
	if err != nil {
		panic(err)
	}

	if err := g.digitalDirectionPin.SetDirection(embd.In); err != nil {
		panic(err)
	}
	g.digitalDirectionPin.ActiveLow(false)

	err = g.digitalSpeedPin.Watch(embd.EdgeRising, func(direction embd.DigitalPin) {
		now := time.Now()
		g.Channel <- Pulse{Long: false, At: now}
		time.Sleep(1 * time.Millisecond) // we need speed pulse come always after direction pulse
		dir, _ := g.digitalDirectionPin.Read()
		if  dir == 1 {
			g.Channel <- Pulse{Long: true, At: now}
		}
	})
	if err != nil {
		panic(err)
	}

}

func (g *GPIO) Stop() {
	close(g.Channel)

	g.digitalDirectionPin.Close()
	g.digitalSpeedPin.Close()

	embd.CloseGPIO()
}

