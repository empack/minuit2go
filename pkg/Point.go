package minuit

type Point struct {
	first  float64
	second float64
}

func NewPoint(first, second float64) *Point {
	return &Point{
		first:  first,
		second: second,
	}
}
