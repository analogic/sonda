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
	var buffer []Pulse

	for current := range inputPulses {

		if len(buffer) == 2 &&
			!current.Long &&
			current.At.Sub(buffer[1].At).Nanoseconds() < (buffer[1].At.Sub(buffer[0].At).Nanoseconds() / 10) * 7 {

			//fmt.Printf("\ndiff: %v, previous: %v\n", current.at.Sub(buffer[1].at).Nanoseconds() , (buffer[1].at.Sub(buffer[0].at).Nanoseconds() / 10) * 7)

			current.Invalid = true
			current.Reason = "t"
		}

		outputPulses <- current

		if !current.Long && !current.Invalid {
			if len(buffer) > 1 {
				buffer = append(buffer[1:], current)
			} else {
				buffer = append(buffer, current)
			}
		}
	}
}
