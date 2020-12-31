package charter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"bitbucket.org/peregrinetraders/mc/pkg/store"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/go-chi/chi"
	gochart "github.com/regorov/go-chart"
	"github.com/regorov/go-chart/drawing"
	"github.com/rs/zerolog"
)

type ValueReaderFunc func(c *Curve) error
type PointFunc func(c *Curve) ([]time.Time, []float64)

// RedisPnLReader
func RedisPnLReader(zl *zerolog.Logger, rs store.Storer, redisListName string) ValueReaderFunc {
	return func(cu *Curve) error {

		llen, err := rs.LLen(redisListName)
		if err != nil {
			return err
		}

		culen := cu.Len()
		//fmt.Printf("redis list %s len=%d, cu.len=%d\n", redisListName, llen, culen)
		if llen < culen {
			// platform restart !!!
			cu.Reset()
			culen = 0
		}

		rows, err := rs.ListRange(redisListName, culen, -1)
		if err != nil {
			return err
		}

		if len(rows) == 0 {
			return nil
		}

		vat := convertToValueAtSlice(zl, rows)
		for i := range vat {
			vat[i].At /= 1000
		}

		var (
			lval float64 // last value
			lat  int64   // last at
			lok  bool    // is last value exist (slice is not empty)
		)
		cu.mux.Lock()
		defer cu.mux.Unlock()
		if len(cu.At) > 0 {
			lok = true
			lat = cu.At[len(cu.At)-1].Unix()
			lval = cu.Value[len(cu.Value)-1]
		} else {

		}
		cu.addSlice(vat)

		// Stores separately positive and negative values for having
		// combined red/green chart.

		for i := range vat {
			if lok { // && lval != 0 {
				mid := (lat + vat[i].At) / 2

				if (lval < 0 && vat[i].Value > 0) || (lval > 0 && vat[i].Value <= 0) {

					cu.AtN = append(cu.AtN, time.Unix(mid, 0))
					cu.ValueN = append(cu.ValueN, 0)

					cu.AtP = append(cu.AtP, time.Unix(mid, 0))
					cu.ValueP = append(cu.ValueP, 0)
				}
			}

			lval = vat[i].Value
			lat = vat[i].At
			lok = true

			if vat[i].Value >= 0 {
				cu.AtP = append(cu.AtP, time.Unix(vat[i].At, 0))
				cu.ValueP = append(cu.ValueP, vat[i].Value)
				continue
			}
			cu.AtN = append(cu.AtN, time.Unix(vat[i].At, 0))
			cu.ValueN = append(cu.ValueN, vat[i].Value)
		}
		return nil
	}
}

func RedisPlatformPnLFiller(zl *zerolog.Logger, rs store.Storer, redisListName string, c *Charter) ValueReaderFunc {
	isFirst := true

	return func(cu *Curve) error {

		var vat []ValueAt
		if isFirst {
			isFirst = false
			rows, err := rs.ListRange(redisListName, 0, -1)
			if err != nil {
				return err
			}

			if len(rows) == 0 {
				return nil
			}

			vat = convertToValueAtSlice(zl, rows)

		} else {

			llen, err := rs.LLen(redisListName)
			if err != nil {
				return err
			}

			culen := cu.Len()
			//fmt.Printf("redis list %s len=%d, cu.len=%d\n", redisListName, llen, culen)
			if llen < culen {
				// platform restart !!!
				//println("reset")
				cu.Reset()
				culen = 0
			}

			sum := 0.0
			found := false
			c.Traverse(func(a *Curve) {
				if !strings.HasPrefix(a.code, "PNL:BU:") {
					return
				}
				v, ok := a.lastValue()
				//println(a.code, v, ok)
				if ok {
					found = ok
					sum += v
				}
			})

			if !found {
				return nil
			}

			// if PNL:BU:* started filling already but with zero values
			// becase trading was not started yet. No reason to fill with zeroes
			// the list PNL:BUOYANCE.
			if (int64(sum*10_000_000)/10_000_000) == 0 && culen == 0 {
				return nil
			}

			sum = float64(int64(sum*10_000_000)) / 10_000_000.00

			if lv, ok := cu.LastValue(); ok && lv == sum { //int64(lv*1_000_000) == int64(sum*1_000_000) {
				return nil
			} else {
				//fmt.Println(lv, sum)
			}
			vat = []ValueAt{{
				Value: sum,
				At:    time.Now().Unix(),
			}}

			if err := rs.Rpush(redisListName, fmt.Sprintf("%d:%0.2f", vat[0].At, vat[0].Value)); err != nil {
				zl.Error().Str("errmsg", err.Error()).Msg("rpush to PNL:BUOYANCY failed")
			}
		}

		var (
			lval float64 // last value
			lat  int64   // last at
			lok  bool    // is last value exist (slice is not empty)
		)
		cu.mux.Lock()
		if len(cu.At) > 0 {
			lok = true
			lat = cu.At[len(cu.At)-1].Unix()
			lval = cu.Value[len(cu.Value)-1]
		} else {

		}

		cu.addSlice(vat)

		defer cu.mux.Unlock()
		for i := range vat {
			if lok { // && lval != 0 {
				mid := (lat + vat[i].At) / 2

				if (lval < 0 && vat[i].Value > 0) || (lval > 0 && vat[i].Value <= 0) {

					cu.AtN = append(cu.AtN, time.Unix(mid, 0))
					cu.ValueN = append(cu.ValueN, 0)

					cu.AtP = append(cu.AtP, time.Unix(mid, 0))
					cu.ValueP = append(cu.ValueP, 0)
				}
			}

			lval = vat[i].Value
			lat = vat[i].At
			lok = true

			if vat[i].Value >= 0 {
				cu.AtP = append(cu.AtP, time.Unix(vat[i].At, 0))
				cu.ValueP = append(cu.ValueP, vat[i].Value)
				continue
			}
			cu.AtN = append(cu.AtN, time.Unix(vat[i].At, 0))
			cu.ValueN = append(cu.ValueN, vat[i].Value)
		}
		return nil
	}
}

// RedisHashAttrReader
func RedisHashAttrReader(zl *zerolog.Logger, rs store.Storer, hashName, attr string, listprefix string) ValueReaderFunc {

	isFirst := true

	return func(cu *Curve) error {

		if isFirst {
			isFirst = false
			rows, err := rs.ListRange(listprefix+":"+hashName, 0, -1)
			if err != nil {
				return err
			}

			if len(rows) == 0 {
				return nil
			}

			vat := convertToValueAtSlice(zl, rows)
			cu.AddSlice(vat)
			for i := range vat {
				vat[i].At *= 1000
			}
			return nil
		}

		sval, err := rs.HGet(hashName, attr)
		if err != nil {
			return err
		}
		if len(sval) == 0 {
			return nil
		}

		last, err := strconv.ParseFloat(sval, 64)
		if err != nil {
			return err
		}
		if math.IsNaN(last) || math.IsInf(last, 0) {
			return nil
		}

		lasti := int64(last * 1000_000)

		// do not store ZERO values if the value if first in the curve
		if lasti == 0 && cu.Len() == 0 {
			return nil
		}

		//if len(cu.Value) > 0 && cu.Value[len(cu.Value)-1] == last {
		if lval, ok := cu.LastValue(); ok && int64(lval*1000_000) == lasti {
			return nil
		}

		at := time.Now().Unix()
		cu.Push(at, last)

		save := []string{
			listprefix + ":" + hashName,
			strconv.FormatInt(at, 10) + ":" + strconv.FormatFloat(last, 'f', 6, 64),
		}

		if err := rs.Do("RPUSH", save...); err != nil {
			zl.Error().Str("list", save[0]).Str("errmsg", err.Error()).Msg("saving values to redis list failed")
		}

		return nil
	}
}

type version struct {
	generation int64
	etag       string
	pngbuf     *bytes.Buffer
}

func (c *Charter) readValueAtSet(redisSetName string, from int) ([]ValueAt, error) {
	rows, err := c.store.ListRange(redisSetName, from, -1)
	if err != nil {
		return nil, err
	}

	return convertToValueAtSlice(&c.log, rows), nil
}
func (c *Charter) Curve(id string) *Curve {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.curve[id]
}

type Charter struct {
	store store.Storer
	log   zerolog.Logger
	mux   sync.RWMutex
	//list  []*Curve
	//idx   map[string]int

	curve map[string]*Curve

	mulmux sync.RWMutex
	// multi holds complex PnL charts
	multi map[string]version
}

func NewCharter(s store.Storer, logger *zerolog.Logger) *Charter {
	return &Charter{
		store: s,
		log:   logger.With().Str("layer", "charter").Logger(),
		// list:  make([]*Curve, 0, 500),
		// idx:   make(map[string]int, 0),
		multi: make(map[string]version),
		curve: make(map[string]*Curve),
	}
}

func (c *Charter) Reset(key, value string) {
	c.log.Info().Msg("cache reset requested")
	c.mux.RLock()
	defer c.mux.RUnlock()
	for _, cu := range c.curve {
		cu.Reset()
	}
}

func buildFullCode(code, attr string) string {

	if attr == "" {
		return code
	}

	return code + "." + attr
}

func (c *Charter) RegisterCurve(cu *Curve) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if _, ok := c.curve[cu.fullCode]; ok {
		return
	}
	c.addCurve(cu)
	go func() {
		t := time.NewTicker(time.Second)
		for {
			select {
			case <-t.C:
				if err := cu.valueReaderFunc(cu); err != nil {
					c.log.Error().Str("errmsg", err.Error()).Msg("redis reading failed")
				}
			case <-cu.stop:
				c.log.Info().Str("hash", cu.fullCode).Msg("stop sync request accepted")
				return
			}
		}
	}()
}
func (c *Charter) DeleteCurve(code string) {
	c.mux.Lock()
	cu, ok := c.curve[code]
	if ok {
		cu.stop <- struct{}{}
		delete(c.curve, code)
	}
	c.mux.Unlock()
}

// RedisSetNames returns map of redis sets starting with prefix. The key of map is set name without prefix.
// Value is redis set name.
func (c *Charter) RedisSetNames(ctx context.Context, prefix string) (map[string]string, error) {

	res := make(map[string]string, 0)

	if !strings.HasSuffix(prefix, ":") {
		prefix += ":"
	}

	keys := []string{}
	err := c.store.Keys(prefix+"*", &keys)
	if err != nil {
		return nil, err
	}

	for j := range keys {
		code := keys[j][len(prefix):]
		res[code] = keys[j]
	}
	return res, nil
}

// AddCurve adds value curve. It's ignored if exist already.
func (c *Charter) AddCurve(vac *Curve) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.addCurve(vac)
}

func (c *Charter) addCurve(vac *Curve) {
	if _, ok := c.curve[vac.fullCode]; ok {
		return
	}
	// c.list = append(c.list, vac)
	// c.idx[vac.code] = len(c.list) - 1
	c.curve[vac.fullCode] = vac
}

// LastValue returns last value related to curve code. Return false if
// curve not found or there is no values.
func (c *Charter) LastValue(code string) (float64, bool) {
	var res float64

	c.mux.RLock()
	cu, ok := c.curve[code]
	if ok {
		cu.mux.RLock()
		res, ok = cu.lastValue()
		cu.mux.RUnlock()
	}
	c.mux.RUnlock()
	return res, ok
}

func (c *Charter) Traverse(f func(c *Curve)) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	for _, c := range c.curve {
		f(c)
	}
}

func (c *Charter) AddValue(code string, vat ValueAt) bool {

	c.mux.RLock()
	cu, ok := c.curve[code]
	if ok {
		cu.Push(vat.At, vat.Value)
	}
	c.mux.RUnlock()
	return ok
}

// Len returns amount of Curve elements identified by code.
func (c *Charter) Len(code string) int {
	c.mux.RLock()
	defer c.mux.RUnlock()
	if cu, ok := c.curve[code]; ok {
		return cu.Len()
	}
	return 0
}

// MultiChartImage returns a chart image representing a single curve.
func (c *Charter) MultiChartImage(w http.ResponseWriter, r *http.Request) {

	a := r.FormValue("a")
	if a == "" {
		c.SingleChartImage(w, r)
		return
	}

	codes := strings.Split(a, ",")

	colors := []string{}

	if cc := r.FormValue("colors"); len(cc) > 0 {
		colors = strings.Split(cc, ",")
	}
	//fmt.Printf("len(cc)=%d %v; %v\n", len(cc), cc, codes)
	if len(colors) > 0 && len(codes) > len(colors) {
		w.WriteHeader(400)
		w.Write([]byte("amount of colors shall be equal to amount of attributes"))
		return
	}
	if len(colors) > 0 {
		for i := range colors {
			if len(colors[i]) != 3 && len(colors[i]) != 6 {
				w.WriteHeader(400)
				w.Write([]byte("a single color only 3 or 6 hexadecimal symbols"))
				return
			}
		}
	}

	widths := []string{}

	if cc := r.FormValue("widths"); len(cc) > 0 {
		widths = strings.Split(cc, ",")
	}
	if len(widths) > 0 && len(codes) > len(widths) {
		w.WriteHeader(400)
		w.Write([]byte("amount of widths shall be equal to amount of attributes"))
		return
	}

	width, err := strconv.Atoi(r.FormValue("w"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("mandatory url parameter 'w' is invalid"))
		return
	}

	height, err := strconv.Atoi(r.FormValue("h"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("mandatory url parameter 'h' is invalid)"))
		return
	}

	id := chi.URLParam(r, "id")

	//fmt.Printf("id=%s a=%s\n", id, a)

	var cus []*Curve

	isLong := false

	c.mux.RLock()

	var (
		xcolors []string
		xwidths []float64
	)
	for i, code := range codes {
		key := id + "." + strings.TrimSpace(code)
		cu, ok := c.curve[key]
		if !ok {
			//c.log.Error().Str("key", key).Msg("not found multichart attribute")
			continue
		}
		//fmt.Printf("key %s=%v", key, cu)
		cu.mux.Lock()
		n := len(cu.Value)
		if n == 0 {
			cu.mux.Unlock()
			continue
		}

		if n > 1 {
			isLong = true
		}
		defer cu.mux.Unlock()
		cus = append(cus, cu)
		if len(colors) > 0 {
			xcolors = append(xcolors, colors[i])
		}
		if len(widths) > 0 {
			v, _ := strconv.ParseFloat(widths[i], 64)
			xwidths = append(xwidths, v)
		}
	}
	c.mux.RUnlock()

	if len(cus) == 0 {
		w.WriteHeader(404)
		if _, err := EmptyImageWithText(width, height, "---", w); err != nil {
			c.log.Error().Str("errmsg", err.Error()).Msg("building empty png failed")
		}
		return
	}

	if !isLong {
		w.WriteHeader(404)
		if _, err := EmptyImageWithText(width, height, "...", w); err != nil {
			c.log.Error().Str("errmsg", err.Error()).Msg("building empty png failed")
		}
		return
	}

	// cetag := r.Header.Get("If-None-Match")
	// setag := cu.ETag()
	// if cetag != "" && cetag == setag {
	// 	w.WriteHeader(304)
	// 	return
	// }

	w.Header().Set("Cache-Control", "no-store")

	w.Header().Set("Content-Type", "image/png")

	pngbuf := bytes.NewBuffer(nil)
	if err := RenderMulti(cus, width, height, xcolors, xwidths, pngbuf); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		c.log.Error().Str("errmsg", err.Error()).Msg("chart render failed")
		return
	}

	//w.Header().Set("ETag", setag)
	_, err = pngbuf.WriteTo(w)
	if err != nil {
		c.log.Error().Str("errmsg", err.Error()).Msg("http response writing failed")
	}
}

func (c *Charter) SingleChartImage(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	width, err := strconv.Atoi(r.FormValue("w"))
	if err != nil {
		w.WriteHeader(400)
		return
	}

	height, err := strconv.Atoi(r.FormValue("h"))
	if err != nil {
		w.WriteHeader(400)
		return
	}

	c.mux.RLock()
	//fmt.Printf("curves=%v\n", c.curve)
	cu, ok := c.curve[id]
	if !ok {
		c.mux.RUnlock()
		w.WriteHeader(404)
		if _, err := EmptyImageWithText(width, height, "---", w); err != nil {
			c.log.Error().Str("errmsg", err.Error()).Msg("building empty png failed")
		}
		c.log.Error().Str("code", id).Msg("curve not found")
		return
	}
	c.mux.RUnlock()

	cetag := r.Header.Get("If-None-Match")
	setag := cu.ETag()
	if cetag != "" && cetag == setag {
		w.WriteHeader(304)
		return
	}

	w.Header().Set("Cache-Control", "no-store")

	w.Header().Set("Content-Type", "image/png")

	color := "008000"
	if lval, ok := cu.LastValue(); ok && lval < 0 {
		color = "FF0000"
	}

	if err := cu.render(width, height, color); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		c.log.Error().Str("errmsg", err.Error()).Msg("chart render failed")
		return
	}

	w.Header().Set("ETag", setag)
	_, err = cu.WriteTo(w)
	if err != nil {
		c.log.Error().Str("errmsg", err.Error()).Msg("http response writing failed")
	}
}

// SingleChartJSON returns a chart representing a single curve.
func (c *Charter) SingleChartJSON(w http.ResponseWriter, r *http.Request) {

	var err error
	from := 0

	if fstr := r.FormValue("from"); fstr != "" {
		from, err = strconv.Atoi(fstr)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("parameter from has invalid value"))
			return
		}
	}

	id := chi.URLParam(r, "id")

	c.mux.RLock()
	cu, ok := c.curve[id]
	if !ok {
		c.mux.RUnlock()
		w.WriteHeader(404)
		return
	}
	c.mux.RUnlock()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	_, err = cu.WriteJSONTo(w, id, from)
	if err != nil {
		c.log.Error().Str("errmsg", err.Error()).Msg("http response writing failed")
	}
}

type GroupPnL struct {
	Name  string  `json:"name"`
	Value float64 `json:"val"`
}

func (a *Charter) GroupPnLHandler(w http.ResponseWriter, r *http.Request) {

	t := r.FormValue("type")
	except := r.FormValue("except")

	var res []GroupPnL
	a.Traverse(func(c *Curve) {
		if c.code == except {
			return
		}

		//println(c.code, "-", c.fullCode, ";", c.redisSetName)
		if !strings.HasPrefix(c.code, "PNL:") {
			return
		}
		v := 0.0
		switch t {
		case "min":
			v = c.Min()
		case "max":
			v = c.Max()
		default:
			v, _ = c.LastValue()
		}

		res = append(res, GroupPnL{
			Name:  c.code[7:],
			Value: v,
		})
	})
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	_ = json.NewEncoder(w).Encode(res)
}
func (c *Charter) DebugHandler(w http.ResponseWriter, r *http.Request) {

	c.mux.RLock()
	cu, ok := c.curve["PNL:BUOYANCY"]
	c.mux.RUnlock()
	if !ok {
		return
	}
	cu.mux.RLock()
	err := json.NewEncoder(w).Encode(cu)
	if err != nil {
		c.log.Error().Str("errmsg", err.Error()).Msg("json marshaling failed")
	}
	cu.mux.RUnlock()
}

func (c *Charter) MultiChartJSON(w http.ResponseWriter, r *http.Request) {

	var err error

	a := r.FormValue("a")
	if a == "" {
		w.WriteHeader(400)
		w.Write([]byte("parameter a is expected"))
		return
	}

	id := chi.URLParam(r, "id")

	var cus []*Curve

	isLong := false

	c.mux.RLock()
	for _, code := range strings.Split(a, ",") {
		key := id + "." + strings.TrimSpace(code)
		cu, ok := c.curve[key]
		if !ok {
			//c.log.Error().Str("key", key).Msg("not found multichart attribute")
			continue
		}
		//fmt.Printf("key %s=%v", key, cu)
		cu.mux.Lock()
		n := len(cu.Value)
		if n == 0 {
			cu.mux.Unlock()
			continue
		}

		if n > 1 {
			isLong = true
		}
		defer cu.mux.Unlock()
		cus = append(cus, cu)
	}
	c.mux.RUnlock()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if len(cus) == 0 {
		w.WriteHeader(404)
		w.Write([]byte(id + "[" + a + "] not found!"))
		return
	}

	if !isLong {
		w.WriteHeader(404)
		w.Write([]byte("No values"))
		return
	}

	// cetag := r.Header.Get("If-None-Match")
	// setag := cu.ETag()
	// if cetag != "" && cetag == setag {
	// 	w.WriteHeader(304)
	// 	return
	// }

	//w.Header().Set("Cache-Control", "no-store")

	//w.Header().Set("Content-Type", "image/png")

	_, err = WriteMultiJSONTo(cus, w)
	if err != nil {
		c.log.Error().Str("errmsg", err.Error()).Msg("http response writing failed")
	}

}

// DownloadFile
func (c *Charter) DownloadFile(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	c.mux.RLock()
	cu, ok := c.curve[id]
	if !ok {
		c.mux.RUnlock()
		return
	}
	c.mux.RUnlock()

	f := excelize.NewFile()
	sch := f.GetSheetName(f.GetActiveSheetIndex())

	cell, _ := excelize.CoordinatesToCellName(1, 1)
	f.SetCellStr(sch, cell, "At")
	cell, _ = excelize.CoordinatesToCellName(2, 1)
	f.SetCellStr(sch, cell, "Value")

	cu.mux.RLock()
	for i := range cu.At {

		cell, err := excelize.CoordinatesToCellName(1, i+2)
		if err != nil {
			c.log.Error().Str("errmsg", err.Error()).Msg("cell coordinates failed")
			continue
		}
		if err := f.SetCellStr(sch, cell, cu.At[i].Format("02.01.2006 15:04:05")); err != nil {
			c.log.Error().Str("errmsg", err.Error()).Msg("set cell failed")
			continue
		}

		cell, err = excelize.CoordinatesToCellName(2, i+2)
		if err != nil {
			c.log.Error().Str("errmsg", err.Error()).Msg("cell coordinates failed")
			continue
		}
		if err := f.SetCellFloat(sch, cell, cu.Value[i], 3, 64); err != nil {
			c.log.Error().Str("errmsg", err.Error()).Msg("set cell failed")
			continue
		}
	}
	cu.mux.RUnlock()

	f.SetSheetName(sch, id)

	fname := id + ".xlsx"

	//fname = "/home/gera/temp/kot.txt"

	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+fname)
	w.Header().Set("Content-Type", "application/octet-stream")

	//w.Header().Set("Transfer-Encoding", "chunked")
	//w.Header().Set("Content-Type", "text/plain")
	//	w.Header().Set("Content-Encoding", "gzip")

	t := time.Now()

	if _, err := f.WriteTo(w); err != nil {
		if err != io.EOF {
			panic(err)
			//render.Render(w, r, ErrRender(err))
			return
		}
	}
	c.log.Info().Str("filename", fname).Str("dur", time.Since(t).String()).Msg("csv curve downloaded")
}

func envVar(name string, def int) int {
	var (
		res = def
		err error
	)
	e, ok := os.LookupEnv(name)
	if ok {
		res, err = strconv.Atoi(e)
		if err != nil {
			res = def
		}
	}
	return res
}

// RenderMulti renders a few Curves at once.
func RenderMulti(cus []*Curve, width, height int, colors []string, widths []float64, pngbuf *bytes.Buffer) error {

	// maxat holds last At existing among all in cus.
	var maxat time.Time

	// maxin holds index of cus, where maximal time is found.
	maxin := -1

	// maxlen holds maximal Curve length among all in cus.
	// maxlen helps to draw "No value response" if all Curves are empty.
	maxlen := 0
	maxval := -math.MaxFloat64
	maxset := false
	minval := math.MaxFloat64
	minset := false
	// for i := range cus {
	// 	fmt.Printf("%#s.%s ", cus[i].code, cus[i].attr)
	// }
	// fmt.Println("")

	for i := range cus {
		//fmt.Printf("%#s.%s\n", cus[i].code, cus[i].attr)
		n := len(cus[i].Value)

		if n == 0 {
			continue
		}

		if n > maxlen {
			maxlen = n
		}

		if mx := cus[i].At[n-1]; mx.After(maxat) {
			maxin = i
			maxat = mx
		}

		if cus[i].attr == "pnl" {
			continue
		}

		//fmt.Printf("%s min=%f, max=%f\n", cus[i].code, cus[i].Min(), cus[i].Max())
		if mv := cus[i].Max(); mv > maxval {
			maxval = mv
			maxset = true
		}

		if mv := cus[i].Min(); mv < minval {
			minval = mv
			minset = true
		}
	}

	if maxlen < 2 || maxin == -1 {
		if _, err := EmptyImageWithText(width, height, "...", pngbuf); err != nil {
			return err
		}
		return nil
	}

	//fmt.Printf("maxin=%d, maxat=%d\n", maxin, maxat.Unix())
	series := []gochart.Series{}

	for i := range cus {
		if cus[i] == nil {
			continue
		}
		if len(cus[i].At) == 0 {
			continue
		}

		ts := gochart.TimeSeries{
			Style: cus[i].styleFunc(cus[i]),
		}
		if len(colors) > 0 {
			if i == 0 {
				ts.Style.FillColor = drawing.ColorFromHex(colors[i])
				ts.Style.StrokeColor = drawing.ColorTransparent
				ts.Style.StrokeWidth = 0

			} else {
				ts.Style.StrokeColor = drawing.ColorFromHex(colors[i])
			}
		}

		if len(widths) > 0 && i > 0 {
			ts.Style.StrokeWidth = widths[i]
		}

		if maxin != i {
			cus[i].attach(maxat, cus[i].Value[len(cus[i].Value)-1])
		}
		//ts.XValues, ts.YValues = cus[i].At, cus[i].Value

		if cus[i].attr != "pnl" {
			ts.XValues, ts.YValues = cus[i].pointFunc(cus[i])
			series = append(series, ts)
			continue
		}
		//	fmt.Printf("%s min=%f (%b) max=%f(%b)\n", cus[i].code, minval, minset, maxval, maxset)
		if maxset && minset {
			ts.XValues, ts.YValues = adoptedPoints(cus[i], minval, maxval)
			//	fmt.Printf("corrected pnl=%v\n", ts.YValues)
			series = append(series, ts)
		}
	}
	if len(series) == 0 {
		if _, err := EmptyImageWithText(width, height, "///", pngbuf); err != nil {
			return err
		}
		return nil
	}

	graph := gochart.Chart{
		XAxis:  gochart.HideXAxis(),
		YAxis:  gochart.HideYAxis(),
		Width:  width,  // 184,
		Height: height, // 67,
		Background: gochart.Style{
			Padding: gochart.BoxZero,
		},
		Canvas: gochart.Style{
			FillColor: HTMLBackgroundColor,
		},

		Series: series,
	}

	//	err := graph.Render(gochart.PNG, c.pngbuf)
	rend, err := graph.PreRender(gochart.PNG)
	if err != nil {
		return err
	}

	rend.SetFontSize(8)
	for i := range cus {
		if cus[i] == nil {
			continue
		}

		if len(cus[i].Value) == 0 {
			continue
		}
		if i != maxin {
			cus[i].detach()
		}
		for j := range cus[i].labelFunc {
			cus[i].labelFunc[j](cus[i], rend, width, height)
		}
	}

	if err := rend.Save(pngbuf); err != nil {
		return err
	}

	return nil
}

func WriteMultiJSONTo(cus []*Curve, w io.Writer) (int, error) {

	// maxat holds last At existing among all in cus.
	var maxat time.Time

	// maxin holds index of cus, where maximal time is found.
	maxin := -1

	// maxlen holds maximal Curve length among all in cus.
	// maxlen helps to draw "No value response" if all Curves are empty.
	maxlen := 0

	for i := range cus {
		//fmt.Printf("%#v\n", cus)
		n := len(cus[i].Value)

		if n == 0 {
			continue
		}

		if n > maxlen {
			maxlen = n
		}

		if mx := cus[i].At[n-1]; mx.After(maxat) {
			maxin = i
			maxat = mx
		}
	}

	if maxlen < 2 || maxin == -1 {
		return 0, errors.New("No values")
	}

	buf := bytes.NewBufferString(`{`)
	csep := ""

	for i := range cus {
		c := cus[i]
		if c == nil {
			continue
		}
		if len(c.At) == 0 {
			continue
		}

		if maxin != i {
			c.attach(maxat, c.Value[len(c.Value)-1])
		}

		// Do JSON
		buf.Write([]byte(fmt.Sprintf(`%s "%s" : {"color" : "%s", "data" : [`, csep, c.attr, c.styleFunc(c).StrokeColor.String())))

		sep := ""

		for j := 0; j < len(c.At); j++ {
			buf.WriteString(sep)
			buf.WriteString("[")
			buf.WriteString(strconv.FormatInt(c.At[j].Unix(), 10))
			buf.WriteString(",")
			buf.WriteString(strconv.FormatFloat(c.Value[j], 'f', 3, 64))
			buf.WriteString("]")
			sep = ","
		}
		//buf.WriteString(sep)
		buf.WriteString("]}")

		csep = ","
		if i != maxin {
			c.detach()
		}
	}

	buf.WriteString("}")
	n, err := w.Write(buf.Bytes())
	buf.Reset()
	return n, err
}

// EmptyImageWithText writes PNG picture with text in the center.
func EmptyImageWithText(width, height int, text string, w io.Writer) (int, error) {

	png, err := gochart.PNG(width, height)
	if err != nil {
		return 0, err
	}

	gochart.Draw.Box(png, gochart.NewBox(0, 0, width, height), gochart.Style{
		FillColor: HTMLBackgroundColor,
	})

	fs := gochart.StyleTextDefaults()
	fs.FontSize = 8
	fs.FontColor = drawing.Color{R: 192, G: 192, B: 192, A: 255}

	tm := gochart.Draw.MeasureText(png, text, fs)
	gochart.Draw.Text(png, text, (width/2 - tm.Width()/2), height/2+tm.Height()/2, fs)

	return 0, png.Save(w)
}

// PointsAsIs returns curve points as is without changes.
func PointsAsIs(c *Curve) ([]time.Time, []float64) {
	return c.At, c.Value
}

// PointsSquared returns points always horizontal or vertical lines.
func PointsSquared(c *Curve) ([]time.Time, []float64) {
	var (
		t []time.Time
		v []float64
	)
	for i := range c.At {
		if i > 0 {
			t = append(t, c.At[i])
			v = append(v, c.Value[i-1])
		}
		t = append(t, c.At[i])
		v = append(v, c.Value[i])
	}
	return t, v
}

func adoptedPoints(c *Curve, low, hi float64) ([]time.Time, []float64) {
	var (
		t []time.Time
		v []float64
	)
	// 	Ginatullin Ildar, [16.10.20 11:22]
	// min max находишь (0, 170)

	// Ginatullin Ildar, [16.10.20 11:22]
	// значит это будет соответствовать (2, 10)

	// Ginatullin Ildar, [16.10.20 11:22]
	// значит шкалы соответствуют 170 к 8

	// Ginatullin Ildar, [16.10.20 11:24]
	// то есть 2 + [ arr[i] - min(arr) ] * (8/170)

	min, max := c.Min(), c.Max()
	max += 0.000001
	//fmt.Printf("pnl %s min=%f, max=%f low=%f, hi=%f\n", c.code, min, max, low, hi)

	for i := range c.At {

		nv := low + (c.Value[i]-min)*((hi-low)/max)
		v = append(v, nv)
		t = append(t, c.At[i])
	}
	return t, v
}

// func TimeArrayToFloat(in []time.Time) []float64 {
// 	res := make([]float64, len(in))
// 	for i := range in {
// 		res[i] = gochart.TimeToFloat64(in[i])
// 	}
// 	return res
// }

func WritePositiveNegativeJSONTo(cu *Curve, from int, w io.Writer) (int, error) {

	cu.mux.RLock()
	defer cu.mux.RUnlock()

	if len(cu.At) == 0 {
		return 0, nil
	}

	if from >= len(cu.AtN) && from >= len(cu.AtP) {
		return 0, nil
	}

	// Do JSON
	buf := bytes.NewBufferString(`{ "positive" : {"color" : "green", "data" : [`)

	sep := ""

	for j := from; j < len(cu.AtP); j++ {
		buf.WriteString(sep)
		buf.WriteString("[")
		buf.WriteString(strconv.FormatInt(cu.AtP[j].Unix(), 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatFloat(cu.ValueP[j], 'f', 3, 64))
		buf.WriteString("]")
		sep = ","
	}
	buf.WriteString(`]},"negative" : { "color" : "red", "data" : [`)
	sep = ""
	for j := from; j < len(cu.AtN); j++ {
		buf.WriteString(sep)
		buf.WriteString("[")
		buf.WriteString(strconv.FormatInt(cu.AtN[j].Unix(), 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatFloat(cu.ValueN[j], 'f', 3, 64))
		buf.WriteString("]")
		sep = ","
	}

	buf.WriteString("]}}")

	n, err := w.Write(buf.Bytes())
	buf.Reset()
	return n, err
}
