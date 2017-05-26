package sonda

func FilterPulsesByLogic(inputPulses chan Pulse, outputPulses chan Pulse) {
	defer close(outputPulses)
	last := Pulse{Long: true}

	for p := range inputPulses {
		if p.Long && last.Long {
			p.Invalid = true
			p.Reason = "l"
		}

		outputPulses <- p
		last = p
	}
}

func FilterPulsesByTimes(inputPulses chan Pulse, outputPulses chan Pulse) {
	defer close(outputPulses)

	var lastDirectionPulse Pulse
	var lastSpeedPulse Pulse

	var diffs []int
	var min int

	for current := range inputPulses {

		if len(diffs) > 30 &&
			!current.Long &&
			(current.At.Nanosecond() - lastSpeedPulse.At.Nanosecond() < (min * 7)/10) {

			current.Invalid = true
			current.Reason = "t"
		}

		outputPulses <- current

		if current.Long {
			if len(diffs) > 30 {
				diffs = append(diffs[1:], current.At.Nanosecond() - lastDirectionPulse.At.Nanosecond())
			} else {
				diffs = append(diffs, current.At.Nanosecond() - lastDirectionPulse.At.Nanosecond())
			}

			n := diffs[0]
			// find smallest
			for _, v := range diffs {
				if v < n {
					n = v
				}
			}
			min = n
			lastDirectionPulse = current
		} else {
			if(!current.Invalid) {
				lastSpeedPulse = current
			}
		}
	}
}
