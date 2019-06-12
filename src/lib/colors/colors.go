package colors

import "fmt"

type Color int

const (
	Reset        Color = 0
	Black        Color = 30
	Red          Color = 31
	Green        Color = 32
	Yellow       Color = 33
	Blue         Color = 34
	Magenta      Color = 35
	Cyan         Color = 36
	Defaut       Color = 39
	LightGray    Color = 37
	DarkGray     Color = 90
	LightRed     Color = 91
	LightGreen   Color = 92
	LightYellow  Color = 93
	LightBlue    Color = 94
	LightMagenta Color = 95
	LightCyan    Color = 96
	White        Color = 97
)

func (c Color) ToString() string {
	return fmt.Sprintf("\033[%dm", c)
}
