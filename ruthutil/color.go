package ruthutil

type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
	Alpha uint8
}

func NewColor(red uint8, green uint8, blue uint8, alpha uint8) Color {
	return Color{alpha, blue, green, red}
}
