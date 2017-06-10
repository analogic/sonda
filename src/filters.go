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

	var diffs []int64
	var min int64

	speedPulsesSinceLastDirectionPulse := 0

	for current := range inputPulses {

		if len(diffs) > 30 &&
			!current.Long &&
			(current.At.Sub(lastSpeedPulse.At).Nanoseconds() < (min * 7)/10) {

			current.Invalid = true
			current.Reason = "t"
		}

		outputPulses <- current

		if current.Long {
			// we will meassure only if there will be pulses close together, we will not meassure in 10° case

			if speedPulsesSinceLastDirectionPulse < 3 {
				if len(diffs) > 30 {
					diffs = append(diffs[1:], current.At.Sub(lastDirectionPulse.At).Nanoseconds())
				} else {
					diffs = append(diffs, current.At.Sub(lastDirectionPulse.At).Nanoseconds())
				}
			}

			min = diffs[0]
			// find smallest
			for _, v := range diffs {
				if v < min {
					min = v
				}
			}
			lastDirectionPulse = current
			speedPulsesSinceLastDirectionPulse = 0;
		} else {
			if(!current.Invalid) {
				lastSpeedPulse = current
				speedPulsesSinceLastDirectionPulse++;
			}
		}
	}
}
