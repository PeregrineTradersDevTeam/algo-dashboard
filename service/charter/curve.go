package charter

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	gochart "github.com/regorov/go-chart"
	"github.com/regorov/go-chart/drawing"
	"github.com/rs/zerolog"
)

var HTMLBackgroundColor = drawing.Color{R: 0x1e, G: 0x1e, B: 0x1e, A: 255}

// Curve holds values and related timestamps and becomes source for chart.
type Curve struct {
	version
	At    []time.Time `json:"at"`
	Value []float64   `json:"v"`

	//
	AtP []time.Time `json:"atp"`
	AtN []time.Time `json:"atn"`

	// ValueP holds only values above or equal to zero
	ValueP []float64 `json:"vp,omitempty"`
	// ValueN holds only values below zero
	ValueN []float64 `json:"vn,omitempty"`

	MaxIndex int
	MinIndex int

	// isFillingExternally holds true if rows into the set added by external party.
	// as instance redis set PNL:BUOYANCY is filled by trading server.
	// but LASTPRICE:XXXX is filled by the adash.
	isFillingExternally bool

	isTimeInUnixNanoFormat bool

	// chartRenderRequired holds true if At or Value has been changed.
	chartRenderRequired bool
	fullCode            string // M:A002.op
	code                string // M:A002
	attr                string // op
	prefix              string
	redisSetName        string
	mux                 sync.RWMutex

	labelFunc       []LabelFunc
	valueReaderFunc ValueReaderFunc
	styleFunc       StyleFunc
	stop            chan struct{}
	pointFunc       PointFunc
}

func (c *Curve) Attr() string {
	return c.attr
}

func NewCurve(code, attr string, vrf ValueReaderFunc, lsf StyleFunc, pf PointFunc, f ...LabelFunc) *Curve {
	c := Curve{
		version:         version{pngbuf: bytes.NewBuffer(nil)},
		At:              make([]time.Time, 0, 1000),
		Value:           make([]float64, 0, 1000),
		code:            code,
		attr:            attr,
		fullCode:        buildFullCode(code, attr),
		MaxIndex:        -1,
		MinIndex:        -1,
		labelFunc:       append([]LabelFunc{}, f...),
		styleFunc:       lsf,
		pointFunc:       pf,
		valueReaderFunc: vrf,
		stop:            make(chan struct{}, 0),
	}

	return &c
}

func (c *Curve) Len() int {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return len(c.Value)
}

func (cu *Curve) LastValue() (float64, bool) {
	cu.mux.RLock()
	res, ok := cu.lastValue()
	cu.mux.RUnlock()
	return res, ok
}

func (cu *Curve) lastValue() (float64, bool) {
	var (
		res float64
		ok  bool
	)
	if k := len(cu.Value); k > 0 {
		res = cu.Value[k-1]
		ok = true
	}
	return res, ok
}

func (c *Curve) ETag() string {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.etag
}

type ValueAt struct {
	At    int64
	Value float64
}

func (c *Curve) attach(at time.Time, val float64) {
	c.Value = append(c.Value, val)
	c.At = append(c.At, at)
}

// Push adds elements in the end of the curve.
func (c *Curve) Push(at int64, val float64) {
	c.mux.Lock()
	c.push(at, val)
	c.mux.Unlock()
}

func (c *Curve) push(at int64, val float64) {

	c.Value = append(c.Value, val)
	c.At = append(c.At, time.Unix(at, 0))
	if c.MaxIndex >= 0 {
		if c.Value[c.MaxIndex] < val {
			c.MaxIndex = len(c.Value) - 1
		}
	} else {
		c.MaxIndex = 0
	}

	if c.MinIndex >= 0 {
		if c.Value[c.MinIndex] > val {
			c.MinIndex = len(c.Value) - 1
		}
	} else {
		c.MinIndex = 0
	}

	c.chartRenderRequired = true
}

func (c *Curve) detach() {
	if len(c.Value) == 0 {
		return
	}
	c.Value = c.Value[0 : len(c.Value)-1]
	c.At = c.At[0 : len(c.At)-1]

	c.chartRenderRequired = true
}

func (c *Curve) Reset() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.reset()
}

func (c *Curve) reset() {
	c.At = nil    //c.At[:0]
	c.Value = nil // c.Value[:0]
	c.ValueN = nil
	c.AtN = nil
	c.ValueP = nil
	c.AtP = nil
	c.MaxIndex = -1
	c.MinIndex = -1
}

func (c *Curve) Max() float64 {
	if c.MaxIndex == -1 {
		return 0
	}
	return c.Value[c.MaxIndex]
}

func (c *Curve) Min() float64 {
	if c.MinIndex == -1 {
		return 0
	}
	return c.Value[c.MinIndex]
}

// AddSlice adds a slice of ValueAt wrapped by a single Mutex.
func (vat *Curve) AddSlice(p []ValueAt) {
	vat.mux.Lock()
	defer vat.mux.Unlock()
	vat.addSlice(p)
}

func (vat *Curve) addSlice(p []ValueAt) {
	for i := range p {
		vat.push(p[i].At, p[i].Value)
	}
	vat.chartRenderRequired = true
}

func convertToValueAtSlice(zl *zerolog.Logger, s []string) []ValueAt {
	var res []ValueAt
	for i := range s {

		// TODO: optimize by not using Split
		el := strings.Split(s[i], ":")
		at, err := strconv.ParseInt(el[0], 10, 64)
		if err != nil {
			zl.Error().Str("errmsg", err.Error()).Str("at", el[0]).Msg("string->int64 parsing failed")
			continue
		}

		val, err := strconv.ParseFloat(el[1], 64)
		if err != nil {
			zl.Error().Str("errmsg", err.Error()).Str("value", el[1]).Msg("string->float64 parsing failed")
			continue
		}

		res = append(res, ValueAt{At: at, Value: val})
	}
	return res
}

func (cu *Curve) render(width, height int, hexcolor string) error {
	cu.mux.Lock()
	defer cu.mux.Unlock()

	if len(cu.At) < 2 {
		cu.pngbuf.Reset()
		if _, err := EmptyImageWithText(width, height, "...", cu.pngbuf); err != nil {
			return err
		}
		return nil
	}

	style := cu.styleFunc(cu)

	graph := gochart.Chart{
		XAxis:  gochart.HideXAxis(),
		YAxis:  gochart.HideYAxis(),
		Width:  width,  // 184,
		Height: height, //67,
		Background: gochart.Style{
			Padding: gochart.BoxZero,
		},
		Canvas: gochart.Style{
			FillColor: HTMLBackgroundColor,
		},
		Series: []gochart.Series{
			gochart.TimeSeries{
				XValues: cu.At,
				YValues: cu.Value,
				Style:   style,
			},
		},
	}
	if len(cu.AtP) > 0 || len(cu.AtN) > 0 {
		style.FillColor = drawing.Color{R: 0x23, G: 0xd1, B: 0x60, A: 255}
		graph.Series[0] = gochart.TimeSeries{
			XValues: cu.AtP,
			YValues: cu.ValueP,
			Style:   style,
		}
		style.FillColor = drawing.Color{R: 0xFF, G: 0x00, B: 0x00, A: 255}
		graph.Series = append(graph.Series,
			gochart.TimeSeries{
				XValues: cu.AtN,
				YValues: cu.ValueN,
				Style:   style,
			})
	}

	// rend, _ := gochart.PNG(width, height)
	// graph.Series[0].Render(rend, gochart.BoxZero, &gochart.ContinuousRange{}, &gochart.ContinuousRange{}, gochart.StyleTextDefaults())

	cu.pngbuf.Reset()
	//	err := graph.Render(gochart.PNG, c.pngbuf)
	rend, err := graph.PreRender(gochart.PNG)
	//err := rend.Save(c.pngbuf)
	if err != nil {
		return err
	}

	for i := range cu.labelFunc {
		cu.labelFunc[i](cu, rend, width, height)
	}

	if err := rend.Save(cu.pngbuf); err != nil {
		return err
	}

	cu.generation++
	cu.etag = strconv.FormatInt(cu.generation, 10)

	return nil
}

func (cu *Curve) renderBackup(width, height int, hexcolor string) error {
	cu.mux.Lock()
	defer cu.mux.Unlock()

	if len(cu.At) < 2 {
		cu.pngbuf.Reset()
		if _, err := EmptyImageWithText(width, height, "...", cu.pngbuf); err != nil {
			return err
		}
		return nil
	}

	style := cu.styleFunc(cu)

	graph := gochart.Chart{
		XAxis:  gochart.HideXAxis(),
		YAxis:  gochart.HideYAxis(),
		Width:  width,  // 184,
		Height: height, //67,
		Background: gochart.Style{
			Padding: gochart.BoxZero,
		},
		Canvas: gochart.Style{
			FillColor: HTMLBackgroundColor,
		},
		Series: []gochart.Series{
			gochart.TimeSeries{
				XValues: cu.At,
				YValues: cu.Value,
				Style:   style,
			},
		},
	}

	// rend, _ := gochart.PNG(width, height)
	// graph.Series[0].Render(rend, gochart.BoxZero, &gochart.ContinuousRange{}, &gochart.ContinuousRange{}, gochart.StyleTextDefaults())

	cu.pngbuf.Reset()
	//	err := graph.Render(gochart.PNG, c.pngbuf)
	rend, err := graph.PreRender(gochart.PNG)
	//err := rend.Save(c.pngbuf)
	if err != nil {
		return err
	}

	for i := range cu.labelFunc {
		cu.labelFunc[i](cu, rend, width, height)
	}

	if err := rend.Save(cu.pngbuf); err != nil {
		return err
	}

	cu.generation++
	cu.etag = strconv.FormatInt(cu.generation, 10)

	return nil
}

func (c *Curve) WriteTo(w io.Writer) (int64, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return io.Copy(w, c.pngbuf)
}

func (c *Curve) WriteJSONTo(w io.Writer, name string, from int) (int, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	buf := bytes.NewBufferString(`{"name":"` + name + `", "data":[`)

	sep := ""

	for i := from; i < len(c.At); i++ {
		buf.WriteString(sep)
		buf.WriteString("[")
		buf.WriteString(strconv.FormatInt(c.At[i].Unix(), 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatFloat(c.Value[i], 'f', 3, 64))
		buf.WriteString("]")
		sep = ","
	}

	buf.WriteString("]}")
	n, err := w.Write(buf.Bytes())
	buf.Reset()
	return n, err
}
