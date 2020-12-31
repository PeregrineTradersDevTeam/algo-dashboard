package adash

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	//"bitbucket.org/elephantsoft/rdk/tracker"
	"bitbucket.org/peregrinetraders/mc/pkg/websocket"
	"github.com/regorov/go-chart"
	"github.com/rs/zerolog"
)

var (
	// ErrImcomingMessageEmpty defines error.
	ErrImcomingMessageEmpty = errors.New("imcoming message is empty")
)

type InstanceEx struct {
	Instance
	Color string `json:"color"`
	// IsBeforeZero bool   `json:"ibz"` deleted in accordance to Floris request
}

type InstanceKeySet struct {
	ID           string
	I            string
	M            string
	S            string
	L            string
	AlgoType     string
	CurrencyCode string
	Group        string
}

// Service  .
type Service struct {
	store Storer

	mux  sync.RWMutex
	keys []InstanceKeySet

	muxg sync.RWMutex

	log zerolog.Logger

	hub *websocket.Hub

	//wsproxy chan tracker.Change

	cfgmux sync.RWMutex
	cfg    map[string]string

	scanInterval chan time.Duration

	pnl      chart.TimeSeries
	lcfolder string
}

func NewService(store Storer, logger *zerolog.Logger, lcfolder string) *Service {
	s := Service{store: store,
		log:          logger.With().Str("layer", "service").Logger(),
		scanInterval: make(chan time.Duration),
		cfg:          make(map[string]string),
		lcfolder:     lcfolder,
	}
	s.pnl.Name = "PnL"
	s.pnl.Style = chart.StyleTextDefaults()
	s.pnl.YAxis = chart.YAxisSecondary
	s.NewPnLChart()
	return &s
}

func (s *Service) Start(ctx context.Context) error {

	_ = s.RefreshConfig()

	go s.refreshConfigRunner(ctx)
	go s.instanceScanner(ctx)
	return nil
}

func (s *Service) SetSession(k, v string) {
	s.cfgmux.Lock()
	s.cfg["session"] = v
	s.cfgmux.Unlock()
}

func (s *Service) RefreshConfig() error {
	var (
		zone, platform, site, env, ppnl, build, pdlc string
	)

	if err := s.store.Get("CFG:ZONE", &zone); err != nil {
		return err
	}

	if err := s.store.Get("CFG:ENV", &env); err != nil {
		return err
	}

	if err := s.store.Get("CFG:PLATFORM", &platform); err != nil {
		return err
	}

	if err := s.store.Get("CFG:SITE", &site); err != nil {
		return err
	}

	if err := s.store.Get("CFG:PPNL", &ppnl); err != nil {
		return err
	}

	if err := s.store.Get("CFG:BUILD", &build); err != nil {
		return err
	}

	if err := s.store.Get("CFG:PDLC", &pdlc); err != nil {
		return err
	}

	s.cfgmux.Lock()
	s.cfg["zone"] = zone
	s.cfg["env"] = env
	s.cfg["platform"] = platform
	s.cfg["site"] = site
	s.cfg["ppnl"] = ppnl
	s.cfg["build"] = build
	s.cfg["pdlc"] = pdlc

	s.cfgmux.Unlock()
	return nil

}

func (s *Service) refreshConfigRunner(ctx context.Context) {

	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.RefreshConfig(); err != nil {
				s.log.Error().Str("errmsg", err.Error()).Msg("reading CFG:* key failed")
			}
		}
	}
}

func (s *Service) instanceScanner(ctx context.Context) {

	var (
		keys []string
	)

	f := func() {
		oldlen := len(keys)
		if err := s.store.Scan(ScanInputKeys, &keys); err != nil {
			s.log.Error().Str("errmsg", err.Error()).Msg("redis SCAN instance names failed")
			return
		}
		// keys has values like I:TAWA, I:TIBO... I:*

		var err error
		altp := make(map[string]string, len(keys))
		curcode := make(map[string]string, len(keys))
		group := make(map[string]string, len(keys))

		for i := range keys {
			altp[keys[i]], err = s.store.HGet(keys[i], "algoType")
			if err != nil {
				s.log.Error().Str("errmsg", err.Error()).Str("key", keys[i]).Str("attr", "algoType").Msg("redis HGET key attr failed")
				continue
			}

			group[keys[i]], err = s.store.HGet(keys[i], "g")
			if err != nil {
				s.log.Error().Str("errmsg", err.Error()).Str("key", keys[i]).Str("attr", "g").Msg("redis HGET key attr failed")
				continue
			}
			curcode[keys[i]], err = s.store.HGet(keys[i], "currencyCode")
			if err != nil {
				s.log.Error().Str("errmsg", err.Error()).Str("key", keys[i]).Str("attr", "currencyCode").Msg("redis HGET key attr failed")
				continue
			}
		}

		s.mux.Lock()
		s.keys = s.keys[0:0]
		for i := range keys {
			iks := InstanceKeySet{
				ID:           keys[i][2:],
				M:            KeyPrefixModel + keys[i][2:],
				S:            KeyPrefixStatus + keys[i][2:],
				I:            KeyPrefixInput + keys[i][2:],
				L:            KeyPrefixLog + keys[i][2:],
				AlgoType:     altp[keys[i]],
				CurrencyCode: curcode[keys[i]],
				Group:        group[keys[i]],
			}
			s.keys = append(s.keys, iks)
		}
		s.mux.Unlock()

		if oldlen != len(keys) {
			s.log.Debug().Int("cnt", len(keys)).Strs("keys", keys).Msg("instance names updated")
		}
	}

	f()
	ticker := time.NewTicker(15 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case interval := <-s.scanInterval:
			ticker.Stop()
			ticker = time.NewTicker(interval)
		case <-ticker.C:
			f()

		}
	}
}

// Instances retrives instances' figures.
func (s *Service) Instances(result *[]InstanceEx) error {

	var keys []InstanceKeySet
	s.mux.RLock()
	keys = make([]InstanceKeySet, len(s.keys))
	copy(keys, s.keys)
	s.mux.RUnlock()

	*result = (*result)[0:0]

	res := make([]InstanceEx, len(keys))

	var wg sync.WaitGroup
	wg.Add(3)
	//dummy := []string{"BUOYANCY", "TREND_TRADER"}
	go func() {
		for i := range keys {
			res[i].ID = keys[i].ID
			res[i].AlgoType = keys[i].AlgoType
			res[i].CurrencyCode = keys[i].CurrencyCode
			res[i].Group = keys[i].Group
			if err := s.store.GetObject(keys[i].M, &res[i].Value); err != nil {
				s.log.Error().Str("errmsg", err.Error()).Msg("key not found")
				continue
			}

			res[i].Value.ReduceLength()
		}
		wg.Done()
	}()

	go func() {
		for i := range keys {
			if err := s.store.GetObject(keys[i].S, &res[i].Status); err != nil {
				s.log.Error().Str("errmsg", err.Error()).Msg("key not found")
				continue
			}

			if (res[i].Status.Status == READY || res[i].Status.Status == RUNNING) && time.Since(time.Unix(res[i].Status.At/1000, 0)) > MarkAsLostAfter {
				res[i].Status.Status = LOST
			}

		}
		wg.Done()
	}()

	go func() {
		var err error
		for i := range keys {
			res[i].LogItemsCount, err = s.store.ListLen(keys[i].L)
			if err != nil {
				s.log.Error().Str("errmsg", err.Error()).Msg("list  not found")
				continue
			}
		}
		wg.Done()
	}()

	//now := time.Now()
	//s.log.Debug().Msg("waiting WG")
	wg.Wait()
	//s.log.Debug().Str("dur", time.Since(now).String()).Msg("waiting done")

	*result = append(*result, res...)

	return nil
}

func (s *Service) PublishLaunchConfig(action, fname string, content string) (bool, error) {

	prefix := ""

	switch action {
	case "launch":
		prefix = "START:"
	case "update":
		prefix = "UPDATE:"
	default:
		return false, errors.New("unknown action code. Supported launch or update")
	}

	key := "FILE:" + time.Now().Format("20060102T150405") + "@" + fname
	err := s.store.Set(key, content)

	if err != nil {
		return false, err
	}
	return s.store.Publish(prefix + key)

}

func (s *Service) Publish(code string) (bool, error) {
	return s.store.Publish(code)
}

func (s *Service) ListLen(key string) (int, error) {
	return s.store.ListLen(key)
}

func (s *Service) ListRange(key string, from, to int) ([]string, error) {
	return s.store.ListRange(key, from, to)
}

func (s *Service) GetObject(key string, result interface{}) error {
	return s.store.GetObject(key, result)
}

func (s *Service) Config() map[string]string {
	res := map[string]string{}
	s.cfgmux.RLock()
	for k, v := range s.cfg {
		res[k] = v
	}
	s.cfgmux.RUnlock()
	return res
}

func (s *Service) PlatformStatus() (*PlatformStatus, error) {

	var (
		err  error
		res  PlatformStatus
		keys []string
		vals []string

		sval Status
	)

	if err := s.store.GetObject("S:PTF", &sval); err != nil {
		s.log.Error().Str("errmsg", err.Error()).Msg("reading hashset S:PTF failed")
	} else {
		res.Status = sval.Status
		res.StatusAt = sval.At
	}

	if res.LogItemsCount, err = s.store.ListLen("L:PTF"); err != nil {
		s.log.Error().Str("errmsg", err.Error()).Msg("reading length of L:PTF is failed")
		return nil, err
	}

	if err := s.store.Scan(ConnectorPrefix+"*", &keys); err != nil {
		s.log.Error().Str("errmsg", err.Error()).Msg("redis SCAN connector names failed")
		return nil, err
	}

	res.ServerTime = time.Now().Format("15:04:05 MST")

	// keys has values like CS:MARKET_DATA, CS:ORDERS
	if len(keys) == 0 {
		// if there is no keys found.
		return &res, nil
	}

	if err := s.store.MGet(keys, &vals); err != nil {
		s.log.Error().Str("errmsg", err.Error()).Msg("redis MGET connector values failed")
		return nil, err
	}

	// amount of keys = amount of vals
	suf := "Connector"
	for i := range keys {
		c := Connector{Name: keys[i][3:]} // CS:

		if strings.HasSuffix(c.Name, suf) {
			c.Name = string([]rune(c.Name)[0 : len(c.Name)-len(suf)])
		}

		val := vals[i]

		if len(val) > 1 {
			v := strings.Split(val, ":")
			c.Status, err = strconv.Atoi(v[0])
			if err != nil {
				return nil, err
			}

			c.StatusAt, err = strconv.ParseUint(v[1], 10, 64)
			if err != nil {
				return nil, err
			}
		}
		res.Connectors = append(res.Connectors, c)
	}

	sort.Slice(res.Connectors, func(i, j int) bool {
		return res.Connectors[i].Name < res.Connectors[j].Name
	})

	return &res, nil
}

type fileStrem struct {
	*zip.Reader
	f *os.File
}

func (fs *fileStrem) Close() error {
	return fs.f.Close()
}

func (s *Service) FileStream(fname string) (io.ReadCloser, error) {

	//Check if file exists and open
	f, err := os.Open(fname)
	if err != nil {
		//File not found, send 404
		return nil, err
	}

	//File is found, create and send the correct headers

	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	//FileHeader := make([]byte, 512)
	//Copy the headers into the FileHeader buffer
	//f.Read(FileHeader)
	//Get content type of file
	//FileContentType :=

	//Get the file size
	//FileStat, _ := Openfile.Stat()                     //Get info from file
	//FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	return f, nil
}

func (s *Service) SaveFileToDisk(fname string, buf bytes.Buffer) error {
	full := filepath.Join(s.lcfolder, fname)
	s.log.Debug().Str("fname", full).Msg("store launcher config file on disk")
	return ioutil.WriteFile(full, buf.Bytes(), 0666)
}

func (s *Service) MatrixI() (map[string]map[string]string, error) {

	var (
		keys []string
	)

	if err := s.store.Scan("I:*", &keys); err != nil {
		s.log.Error().Str("errmsg", err.Error()).Msg("redis SCAN I:*s failed")
		return nil, err
	}

	res := map[string]map[string]string{}
	for i := range keys {
		obj := map[string]string{}
		if err := s.store.GetObject(keys[i], &obj); err != nil {
			s.log.Error().Str("errmsg", err.Error()).Msg("key not found")
			return nil, err
		}
		res[keys[i][2:]] = obj
	}
	return res, nil
}
