package charter

import (
	gochart "github.com/regorov/go-chart"
	"github.com/regorov/go-chart/drawing"
)

type StyleFunc func(c *Curve) gochart.Style

// StyleLastPrice returns style for last price curve.
func StyleLastPrice(c *Curve) gochart.Style {
	return gochart.Style{
		StrokeColor: drawing.ColorFromAlphaMixedRGBA(240, 0, 255, 255), // gochart.ColorWhite,
		// filling deleted in accordance to Vadim's req
		// FillColor:   gochart.ColorBlue.WithAlpha(50),
		//
		StrokeWidth: 0.75,
		Padding:     gochart.BoxZero,
	}
}

// StylePnL returns style for pnl curve. The color changes between Red and Green
// in accordance to last value (below or above zero).
func StylePnL(c *Curve) gochart.Style {

	res := gochart.Style{
		StrokeColor: gochart.ColorLightGray,
		StrokeWidth: 2,
		Padding:     gochart.BoxZero,
	}
	pnl, _ := c.lastValue()
	if pnl > 0 {
		res.FillColor = drawing.Color{R: 0x00, G: 0xFF, B: 0x00, A: 255} //{R: 0x23, G: 0xd1, B: 0x60, A: 255} // gochart.ColorAlternateGreen
	}
	if pnl < 0 {
		res.FillColor = drawing.Color{R: 255, G: 0, B: 0, A: 255}
	}
	res.StrokeColor = res.FillColor.WithAlpha(255)

	return res
}

func StyleMixedPnL(c *Curve) gochart.Style {

	res := gochart.Style{
		StrokeColor: gochart.ColorTransparent,
		StrokeWidth: 0,
		Padding:     gochart.BoxZero,
	}

	return res
}

func StyleInstancePnL(c *Curve) gochart.Style {

	res := gochart.Style{
		Padding:     gochart.BoxZero,
		StrokeWidth: 0,
		StrokeColor: gochart.ColorTransparent,
	}
	pnl, _ := c.lastValue()
	if pnl > 0 {
		res.FillColor = drawing.ColorFromAlphaMixedRGBA(0, 128, 0, 255) // Wi {R: 116, G: 238, B: 21, A: 40} //ColorGreen.WithAlpha(70)
	}
	if pnl < 0 {
		res.FillColor = drawing.Color{R: 128, G: 0, B: 0, A: 255}
	}
	//res.FillColor = res.StrokeColor.
	return res
}

// PnlColor
func PnlColor(last float64) drawing.Color {
	if last >= 0 {
		return gochart.ColorAlternateGreen.WithAlpha(100)
	}
	return drawing.Color{R: 255, G: 0, B: 0, A: 255}.WithAlpha(100)
}

// StyleThreshold returns style for high and low curves.
func StyleThreshold(color drawing.Color) func(*Curve) gochart.Style {
	return func(*Curve) gochart.Style {
		return gochart.Style{
			StrokeColor: color,
			StrokeWidth: 2,
			Padding:     gochart.BoxZero,
		}
	}
}

// StyleOpeningPrice returns style for opening price curve. It's dashed.
func StyleOpeningPrice(color drawing.Color) func(*Curve) gochart.Style {
	return func(*Curve) gochart.Style {
		return gochart.Style{
			StrokeColor:     color,
			StrokeDashArray: []float64{2, 2},
			StrokeWidth:     2,
			Padding:         gochart.BoxZero,
		}
	}
}
