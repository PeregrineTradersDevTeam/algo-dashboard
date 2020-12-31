package charter_test

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"bitbucket.org/peregrinetraders/mc/apps/adash/service/charter"
	"bitbucket.org/peregrinetraders/mc/pkg/store"
	"github.com/regorov/go-chart"
	gochart "github.com/regorov/go-chart"
	"github.com/rs/zerolog"
)

var (
	zl        zerolog.Logger
	chr       *charter.Charter
	instances = []string{"A001", "A002", "A003", "A004", "A005", "A006"}
	prefixes  = []string{"LASTPRICE", "LASTPNL"}
	interval  = [][2]float64{{0.05, 10}, {-2.0, 200}}
	con       store.Storer
)

func TestMain(m *testing.M) {

	zl = zerolog.New(os.Stdout)
	con = store.NewRedisStore("localhost:6379", 0, &zl)
	if err := con.Connect(); err != nil {
		fmt.Printf("redis connection failed: %s", err.Error())
		os.Exit(2)
	}

	chr = charter.NewCharter(con, &zl)

	os.Exit(m.Run())
}

/*func TestMultiChartRender(t *testing.T) {
	tr := redishashtracker.New(con, &zl, chr)

	tr.TrackHashAttribute("M", "l", "LASTPRICE", charter.LabelLast, charter.StyleLastPrice, charter.DataAsIs)
	tr.TrackHashAttribute("M", "hi", "HIPRICE", charter.LabelLastMaxMin, charter.StyleThreshold, charter.PointsSquared)
	//lps.TrackHashAttribute("M", "op", "OPENINGPRICE", charter.CurMaxMin, charter.LastPnLColor, charter.AsIs)

	//chr.Start(ctx)

	tr.Start(context.Background())
	time.Sleep(time.Second)

	var cus []*charter.Curve
	cus = append(cus, chr.Curve("M:A002.l"), chr.Curve("M:A002.hi"))

	var maxat time.Time
	maxin := -1
	maxlen := 0
	width := 184
	height := 67

	series := []gochart.Series{}

	for i := range cus {
		if l := len(cus[i].Value); l > maxlen {
			maxlen = l
		}

		if len(cus[i].At) == 0 {
			continue
		}

		if cus[i].At[len(cus[i].At)-1].After(maxat) {
			maxin = i
		}

	}

	if maxlen < 2 || maxin == -1 {
		t.Error("no values")
	}

	for i := range cus {
		ts := gochart.TimeSeries{
			Style: cus[i].styleFunc(cus[i]),
		}
		ts.XValues, ts.YValues = cus[i].pointFunc(cus[i])
		// XValues: cus[i].At,
		// 	YValues: cus[i].Value,

		series = append(series, ts)
	}

	graph := gochart.Chart{
		XAxis:  gochart.HideXAxis(),
		YAxis:  gochart.HideYAxis(),
		Width:  width,  // 184,
		Height: height, //67,
		Background: gochart.Style{
			Padding: gochart.BoxZero,
		},

		Series: series, ///, []gochart.Series{
		// 	gochart.TimeSeries{
		// 		XValues: c.At,
		// 		YValues: c.Value,
		// 		Style: gochart.Style{
		// 			//StrokeColor: gochart.ColorBlue,
		// 			StrokeColor: color,                //gochart.ColorBlue,                ,
		// 			FillColor:   color.WithAlpha(100), //gochart.ColorBlue.WithAlpha(100), //
		// 			StrokeWidth: 0.75,
		// 			Padding:     gochart.BoxZero,
		// 		},
		// 	},
		// 	// gochart.AnnotationSeries{
		// 	// 	Annotations: []gochart.Value2{
		// 	// 		{XValue: gochart.TimeToFloat64(c.At[3]), YValue: c.Value[3], Label: "One"},
		// 	// 	},
		// 	// },
		// },
	}

	// rend, _ := gochart.PNG(width, height)
	// graph.Series[0].Render(rend, gochart.BoxZero, &gochart.ContinuousRange{}, &gochart.ContinuousRange{}, gochart.StyleTextDefaults())

	pngbuf := bytes.NewBuffer(nil)
	//	err := graph.Render(gochart.PNG, c.pngbuf)
	rend, err := graph.RenderX(gochart.PNG)
	err := rend.Save(c.pngbuf)
	if err != nil {
		t.Errorf("render failed. %v", err)
	}

	// for i := range cus {
	// 	cus[i].labelFunc(cus[i], rend, width, height)
	// }

	fo, err := os.Create("twolines.png")
	if err != nil {
		t.Errorf("create file failed. %v", err)
	}

	if _, err := fo.Write(png.Bytes()); err != nil {
		t.Errorf("write file failed. %v", err)
	}

}

// func TestCharter_Init(t *testing.T) {
// 	if err := lps.Cache(context.Background(), "PNL", true); err != nil {
// 		t.Error(err)
// 	}
// 	if err := lps.Cache(context.Background(), "LASTPRICE", false); err != nil {
// 		t.Error(err)
// 	}

// 	lps.Must("BUOYANCY", "PNL:BUOYANCY", true)

// }

*/
func TestLastPricer_GenLastPrice(t *testing.T) {

	rand.Seed(time.Now().Unix())
	at := time.Now().Add(-2 * time.Minute).Unix()
	save := make([]string, 2, 2)

	for i := range instances {

		for p := range prefixes {

			//price := float64(rand.Intn(10)) + float64(rand.Intn(limits[p]))/100.0

			save[0] = prefixes[p] + ":M:" + instances[i]
			if err := con.Do("DEL", save[0]); err != nil {
				t.Error()
			}
			for s := 0; s < 100; s++ {
				ats := at + int64(s)
				price := interval[p][0] + float64(rand.Intn(int(interval[p][1])+1))
				save[1] = fmt.Sprintf("%d:%0.06f", ats, price)

				if err := con.Do("RPUSH", save...); err != nil {
					t.Error()
				}
			}
		}
	}

	if err := con.Do("DEL", "HIPRICE:M:A002", "LOPRICE:M:A002", "OPPRICE:M:A002"); err != nil {
		t.Error()
	}
	if err := con.Do("RPUSH", "HIPRICE:M:A002", fmt.Sprintf("%d:%0.06f", at, 10.02)); err != nil {
		t.Error()
	}
	if err := con.Do("RPUSH", "LOPRICE:M:A002", fmt.Sprintf("%d:%0.06f", at, 0.3)); err != nil {
		t.Error()
	}
	if err := con.Do("RPUSH", "OPPRICE:M:A002", fmt.Sprintf("%d:%0.06f", at, 1.0)); err != nil {
		t.Error()
	}

}

/*
func TestLastPricer_cache(t *testing.T) {

	lps.list = nil
	lps.idx = map[string]int{}
	if err := lps.Init(context.Background()); err != nil {
		t.Error(err)
	}

	lps.Traverse("A001", func(at time.Time, price float64) {
		//t.Logf("%s - %0.f", at, price)
	})

}
func TestLastPricer_Chart(t *testing.T) {

	buf := bytes.NewBuffer(nil)

	if err := lps.Init(context.Background()); err != nil {
		t.Error(err)
	}

	if err := lps.Chart("A001", buf); err != nil {
		t.Error(err)
	}
	if err := ioutil.WriteFile("./chart.png", buf.Bytes(), 0600); err != nil {
		t.Error(err)
	}
	buf.Reset()
}
*/

func TestTwoLines(t *testing.T) {
	var b float64
	b = 1000

	ts1 := gochart.ContinuousSeries{ //TimeSeries{
		Name:    "Time Series",
		XValues: []float64{10 * b, 20 * b, 30 * b, 40 * b, 50 * b, 60 * b, 70 * b, 80 * b},
		YValues: []float64{1.0, 2.0, 30.0, 4.0, 50.0, 6.0, 7.0, 88.0},
	}

	ts2 := chart.ContinuousSeries{ //TimeSeries{
		Style: chart.Style{
			StrokeColor: chart.GetDefaultColor(1),
		},

		XValues: []float64{10 * b, 20 * b, 30 * b, 40 * b, 50 * b, 60 * b, 70 * b, 80 * b},
		YValues: []float64{15.0, 52.0, 30.0, 42.0, 50.0, 26.0, 77.0, 38.0},
	}

	graph := chart.Chart{

		XAxis: chart.XAxis{
			Name:           "The XAxis",
			ValueFormatter: chart.TimeMinuteValueFormatter, //TimeHourValueFormatter,
		},

		YAxis: chart.YAxis{
			Name: "The YAxis",
		},

		Series: []chart.Series{
			ts1,
			ts2,
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		log.Fatal(err)
	}

	fo, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}

	if _, err := fo.Write(buffer.Bytes()); err != nil {
		panic(err)
	}
}
