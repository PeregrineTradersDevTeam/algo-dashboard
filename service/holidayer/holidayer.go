package holidayer

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
)

// Response defines REST API response structure.
type Response struct {
	Today            string    `json:"today"`
	Dates            []Holiday `json:"dates"`
	IsHolidayToday   bool      `json:"isHolidayToday"`
	IsUpdateRequired bool      `json:"isUpdateRequired"`
	ServerTime       string    `json:"serverTime"`
}

// Holiday describes a single holiday attributes.
type Holiday struct {
	Dt        string `json:"dt"`
	Name      string `json:"name"`
	Market    string `json:"market"`
	OpenTime  string `json:"openTime"`
	CloseTime string `json:"closeTime"`
	Mic       string `json:"mic"`
}

type Holidayer struct {
	log   zerolog.Logger
	fname string

	mux   sync.RWMutex
	cache Response
}

// New returns Holidayer.
func New(logger zerolog.Logger, fname string) *Holidayer {
	return &Holidayer{
		log:   logger.With().Str("layer", "holidayer").Logger(),
		fname: fname,
		cache: Response{IsUpdateRequired: true},
	}
}

func (h *Holidayer) Start(ctx context.Context) error {
	if h.fname == "" {
		return nil
	}

	h.refreshCache()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err = watcher.Add(h.fname); err != nil {
		return err
	}

	go h.runner(ctx, watcher)
	return nil
}

func (h *Holidayer) runner(ctx context.Context, watcher *fsnotify.Watcher) {

	i := 0
	defer watcher.Close()
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-watcher.Events:
			h.log.Debug().Int("i", i).Str("event", event.Op.String()).Str("file", h.fname).Msg("file system notifier")
			if !ok {
				break
			}
			i++
			if err := h.refreshCache(); err != nil {
				h.log.Error().Str("errmsg", err.Error()).Str("file-name", h.fname).Msg("file reading failed")
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				break
			}
			h.log.Error().Str("errmsg", err.Error()).Msg("file system notify failed")
		}
	}
}

func (h *Holidayer) refreshCache() error {

	buf, err := ioutil.ReadFile(h.fname)
	if err != nil {
		return err
	}
	if len(buf) == 0 {
		return nil
	}

	rows := strings.Split(string(buf), "\n")

	now := time.Now().UTC()

	h.mux.Lock()
	defer h.mux.Unlock()

	h.cache.Today = now.Format("2006-01-02")
	h.cache.Dates = h.cache.Dates[:0]
	h.cache.IsHolidayToday = false
	h.cache.IsUpdateRequired = true
	titlePassed := false
	for i := range rows {
		row := strings.TrimSpace(rows[i])
		if len(row) == 0 {
			continue
		}

		limit := 5
		if len(row) < 5 {
			limit = len(row)
		}

		if strings.Contains(row[0:limit], "#") {
			continue
		}
		if !titlePassed {
			titlePassed = true
			continue
		}
		fields := strings.Split(row, ",")

		if fields[0] < h.cache.Today {
			continue
		}

		h.cache.Dates = append(h.cache.Dates, Holiday{
			Dt:        fields[0],
			Name:      fields[1],
			Market:    fields[2],
			OpenTime:  fields[3],
			CloseTime: fields[4],
			Mic:       fields[5],
		})

		if fields[0] == h.cache.Today {
			h.cache.IsHolidayToday = true
		}
	}
	if len(h.cache.Dates) == 0 {
		return nil
	}

	t, err := time.Parse("2006-01-02", h.cache.Dates[len(h.cache.Dates)-1].Dt)
	if err != nil {
		return err
	}

	h.cache.IsUpdateRequired = t.Sub(now) < time.Hour*24*90
	return nil
}

// JSON returns holidays in JSON structure.
func (h *Holidayer) JSON(w io.Writer) error {
	h.mux.RLock()
	defer h.mux.RUnlock()
	h.cache.ServerTime = time.Now().Format("2006-01-02 15:04:05 -0700 MST")
	return json.NewEncoder(w).Encode(h.cache)
}
