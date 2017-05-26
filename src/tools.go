package sonda

import "math"

func AverageSpeed(s *[]float32) float32 {
	var speedSum float32
	speedSum = 0
	for _, num := range *s {
		speedSum += num
	}

	return speedSum / float32(len(*s))
}

func MaxSpeed(s *[]float32) float32 {
	// TODO proper gust calculation
	max := *s[0]
	for _, v := range *s {
		if v > max {
			max = v
		}
	}
	return max
}

func AverageDirection(d *[]int) int {
	var sins []float64
	var coss []float64

	for _, direction := range *d {
		sins = append(sins, math.Sin(Rad(float64(direction))))
		coss = append(coss, math.Cos(Rad(float64(direction))))
	}

	return int(Deg(math.Atan2(SumFloat64(&sins) / len(sins), SumFloat64(&coss) / len(coss))))
}

func SumFloat64(a *[]float64) float64 {
	var sum float64
	sum = 0
	for _, num := range *a {
		sum += num
	}
	return sum
}

func Rad(d float64) float64 {
	return math.Pi/180 * d
};

func Deg(r float64) float64{
	return r / (math.Pi/180)
}
