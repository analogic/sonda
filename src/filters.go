package sonda

//import "fmt"

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

//        var cnt = 0;

	for current := range inputPulses {

		if len(diffs) > 30 &&
			!current.Long &&
			(current.At.Sub(lastSpeedPulse.At).Nanoseconds() < (min * 7)/10) {

			current.Invalid = true
			current.Reason = "t"
		}

//                cnt++
//                if(cnt == 36) {
//                fmt.Printf("min: %v cur: %v\n", min, current.At.Sub(lastSpeedPulse.At).Nanoseconds());
//                cnt = 0
//                }
		outputPulses <- current

		if current.Long {
			if len(diffs) > 30 {
				n := current.At.Sub(lastDirectionPulse.At).Nanoseconds()
				//if n < 3*min {
					diffs = append(diffs[1:], n)
				//}
			} else {
				diffs = append(diffs, current.At.Sub(lastDirectionPulse.At).Nanoseconds())
			}

			min = diffs[0]
			// find smallest
			for _, v := range diffs {
				if v < min {
					min = v
				}
			}
			lastDirectionPulse = current
		} else {
			if(!current.Invalid) {
				lastSpeedPulse = current
			}
		}
	}
}
