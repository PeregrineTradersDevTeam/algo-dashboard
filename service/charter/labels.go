package charter

import (
	"fmt"
	"strconv"

	gochart "github.com/regorov/go-chart"
	"github.com/regorov/go-chart/drawing"
)

type LabelGetterFunc func(c *Curve) string

// LabelFunc describes function parameters
type LabelFunc func(c *Curve, r gochart.Renderer, width, height int)

// LabelMaxMin puts max and min values in to the top/left and top/rigt corners.
func LabelMaxMin(c *Curve, rend gochart.Renderer, width, height int) {

	minlabel := "-"
	maxlabel := "-"

	if c.MinIndex >= 0 {
		minlabel = strconv.FormatFloat(c.Value[c.MinIndex], 'f', 0, 64)
	}

	if c.MaxIndex >= 0 {
		maxlabel = strconv.FormatFloat(c.Value[c.MaxIndex], 'f', 0, 64)
	}

	rend.SetFontSize(9)
	//rend.SetFontColor(gochart.ColorRed) // Color{R: 0x00, G: 0x66, B: 0xff, A: 255})
	rend.SetFontColor(drawing.Color{R: 255, G: 0, B: 0, A: 255})
	rend.Text(minlabel, 4, 12)

	rend.SetFontColor(drawing.Color{R: 0x23, G: 0xd1, B: 0x60, A: 255}) // (drawing.Color{R: 0x00, G: 0xFF, B: 0x00, A: 255})
	rend.Text(maxlabel, width-2-rend.MeasureText(maxlabel).Right, height-2)
}

// LabelLastMaxMin puts last, max and min values in to the top/center, top/left and top/rigt corners.
func LabelLastMaxMin(c *Curve, rend gochart.Renderer, width, height int) {

	LabelMaxMin(c, rend, width, height)

	v, ok := c.lastValue()
	if !ok {
		return
	}
	label := fmt.Sprintf("%0.2f", v)

	rend.SetFontColor(drawing.ColorBlack)
	rend.Text(label, width/2-2-rend.MeasureText(label).Right/2, 10)
}

// LabelLast
func LabelLast(c *Curve, rend gochart.Renderer, width, height int) {

	if len(c.Value) == 0 {
		return
	}

	label := fmt.Sprintf("%0.2f", c.Value[len(c.Value)-1])

	rend.SetFontColor(drawing.ColorBlue)
	rend.SetFontSize(8)
	rend.Text(label, width-2-rend.MeasureText(label).Right, 10)
}

// func LabelLast(c *Curve, rend gochart.Renderer, width, height int) {

// 	if len(c.Value) == 0 {
// 		return
// 	}

// 	label := fmt.Sprintf("%0.2f", c.Value[len(c.Value)-1])

// 	rend.SetFontColor(drawing.ColorBlue)
// 	rend.SetFontSize(8)
// 	rend.Text(label, width-2-rend.MeasureText(label).Right, 10)
// }

func Label(x, y int, label LabelGetterFunc, color drawing.Color) LabelFunc {
	return func(c *Curve, rend gochart.Renderer, width, height int) {
		x1, y1 := x, y
		s := label(c)
		mt := rend.MeasureText(s)

		if x < 0 {
			x1 = width + x - mt.Width()
		}
		if y < 0 {
			y1 = height + y
		}
		rend.SetFontColor(color)
		rend.Text(s, x1, y1)
	}
}

func Color(c drawing.Color) func() drawing.Color {
	return func() drawing.Color {
		return c
	}
}

func Last(c *Curve) string {
	if len(c.Value) == 0 {
		return ""
	}
	return fmt.Sprintf("%0.2f", c.Value[len(c.Value)-1])
}

func Min(c *Curve) string {
	if c.MinIndex < 0 {
		return "-"
	}

	return fmt.Sprintf("%0.2f", c.Value[c.MinIndex])
}

func Max(c *Curve) string {
	if c.MaxIndex < 0 {
		return "-"
	}

	return fmt.Sprintf("%0.2f", c.Value[c.MaxIndex])
}

func LabelMinMaxInterval(c *Curve) string {
	min, max := "", ""
	if c.MinIndex >= 0 {
		min = strconv.FormatFloat(c.Value[c.MinIndex], 'f', 0, 64)
	}
	if c.MinIndex >= 0 {
		max = strconv.FormatFloat(c.Value[c.MaxIndex], 'f', 0, 64)
	}

	return min + "..." + max
}

func MaxFloor(c *Curve) string {
	if c.MaxIndex < 0 {
		return "-"
	}

	return strconv.FormatFloat(c.Value[c.MaxIndex], 'f', 0, 64)
}
