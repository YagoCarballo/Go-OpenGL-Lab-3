package objects

type DrawMode int32
type ColorMode int32

const (
	_ = iota // ignore first value by assigning to blank identifier
	DRAW_POINTS DrawMode = 0 + iota
	DRAW_LINES
	DRAW_POLYGONS
)

const (
	_ = iota // ignore first value by assigning to blank identifier
	COLOR_PER_SIDE ColorMode = 0 + iota
	COLOR_SOLID
)


var drawModeNames = [...]string{
	"_",
	"Draw Points",
	"Draw Lines",
	"Draw Polygons",
}

var colorModeNames = [...]string{
	"Color per side",
	"Solid Color",
}

func (drawMode DrawMode) String() string {
	return drawModeNames[drawMode]
}

func (colorMode ColorMode) String() string {
	return colorModeNames[colorMode]
}