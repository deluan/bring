package bring

const (
	MouseLeft = 1 << iota
	MouseMiddle
	MouseRight
	MouseUp
	MouseDown
)

type KeyCode []int

var (
	KeyBackspace = KeyCode{0xFF08}
	KeyEnter     = KeyCode{0xFF0D}
	KeyUp        = KeyCode{0xFF52}
)
