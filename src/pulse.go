package sonda

import (
	"strings"
	"bytes"
	"time"
)

type Pulse struct {
	Long    bool // long/short
	Invalid bool
	At      time.Time

	Reason  string
}

func (p *Pulse) String() string {
	if p.Invalid {
		if p.Long {
			return strings.ToUpper(p.Reason)
		} else {
			return p.Reason
		}
	}

	if p.Long {
		return "█"
	} else {
		return "▄"
	}
}

type Pulses struct {
	Pulses []Pulse
}

func (p *Pulses) Add(pulse Pulse) {
	p.Pulses = append(p.Pulses, pulse)
}

func (p *Pulses) String() string {
	var buffer bytes.Buffer
	for _, v := range p.Pulses {
		buffer.WriteString(v.String())
	}

	return buffer.String()
}
