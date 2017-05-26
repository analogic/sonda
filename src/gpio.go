package sonda

import (
	"github.com/kidoman/embd"
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
	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}

	go g.initSpeed()
	go g.initDirection()
}

func (g *GPIO) initSpeed() {
	var err error
	g.digitalSpeedPin, err = embd.NewDigitalPin(g.SpeedPin)
	if err != nil {
		panic(err)
	}

	if err := g.digitalSpeedPin.SetDirection(embd.In); err != nil {
		panic(err)
	}
	g.digitalSpeedPin.ActiveLow(false)

	err = g.digitalSpeedPin.Watch(embd.EdgeRising, func(speed embd.DigitalPin) {
		g.Channel <- Pulse{Long: false, At: time.Now()}
	})
	if err != nil {
		panic(err)
	}
}

func (g *GPIO) initDirection() {
	var err error
	g.digitalDirectionPin, err = embd.NewDigitalPin(g.DirectionPin)
	if err != nil {
		panic(err)
	}

	if err := g.digitalDirectionPin.SetDirection(embd.In); err != nil {
		panic(err)
	}
	g.digitalDirectionPin.ActiveLow(false)

	err = g.digitalDirectionPin.Watch(embd.EdgeRising, func(direction embd.DigitalPin) {
		g.Channel <- Pulse{Long: false, At: time.Now()}
	})
	if err != nil {
		panic(err)
	}
}

func (g *GPIO) Stop() {
	close(g.Channel)

	g.digitalDirectionPin.Close()
	g.digitalDirectionPin.Close()

	embd.CloseGPIO()
}

