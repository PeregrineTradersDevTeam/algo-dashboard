package redishashtracker

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"bitbucket.org/peregrinetraders/mc/apps/adash/service/charter"
	"bitbucket.org/peregrinetraders/mc/pkg/store"
	"github.com/rs/zerolog"
)

type Price struct {
	ID        string
	sprice    float64
	isChanged bool
}

type key struct {
	prefix string
	attr   string
}

type trackingKey struct {
	hash string
	attr string
}

type trackingValue struct {
	lastValue float64
	attr      string
}

type attr struct {
	name       string
	labelfunc  []charter.LabelFunc
	stylefunc  charter.StyleFunc
	listprefix string
	pointFunc  charter.PointFunc
}

type TrackingEngine struct {
	store store.Storer
	log   zerolog.Logger
	// mux   sync.RWMutex
	// list  []*PriceCurve
	// idx   map[string]int
	pch *charter.Charter

	mux sync.RWMutex

	// hashprefix holds hash prefixes registered on init phase
	//
	hashprefix map[string][]attr

	// ha holds hash tracking attributes. key is full hashname, values are attributes
	ha map[string][]attr

	// pa value is a prefix of the redis list storing
	// historical values.
	//pa        map[key]prop
	lastvalue map[trackingKey]float64
}

func New(s store.Storer, logger *zerolog.Logger, pch *charter.Charter) *TrackingEngine {
	return &TrackingEngine{
		store:      s,
		log:        logger.With().Str("layer", "tracker").Logger(),
		hashprefix: make(map[string][]attr),
		ha:         make(map[string][]attr),

		pch: pch}
}

func (t *TrackingEngine) TrackHashAttribute(hashprefix string, attrname string, listprefix string, cf charter.StyleFunc, pf charter.PointFunc, labelFunc ...charter.LabelFunc) {
	t.mux.Lock()
	defer t.mux.Unlock()

	a := t.hashprefix[hashprefix]
	a = append(a, attr{
		name:       attrname,
		labelfunc:  append([]charter.LabelFunc{}, labelFunc...),
		stylefunc:  cf,
		listprefix: listprefix,
		pointFunc:  pf,
	})
	t.hashprefix[hashprefix] = a
}

func (t *TrackingEngine) Start(ctx context.Context) error {
	go t.runner(ctx)
	return nil
}

func (lp *TrackingEngine) syncTrackers() error {

	found := ""
	deleted := ""
	sepf := ""
	sepd := ""
	for hp, as := range lp.hashprefix {
		var keys []string
		if err := lp.store.Scan(hp+":*", &keys); err != nil {
			return err
		}

		// add if new hash appeared
		for i := range keys {
			if _, ok := lp.ha[keys[i]]; ok {
				continue
			}

			found += sepf + keys[i]
			sepf = ","
			lp.ha[keys[i]] = as
			for _, attr := range as {
				lp.pch.RegisterCurve(charter.NewCurve(
					keys[i],
					attr.name,
					charter.RedisHashAttrReader(&lp.log, lp.store, keys[i], attr.name, attr.listprefix),
					attr.stylefunc,
					attr.pointFunc,
					attr.labelfunc...,
				),
				)
			}
		}

		// del if hash dissapeared. can be optimized
		for h, a := range lp.ha {
			found := false
			for i := range keys {
				if keys[i] == h {
					found = true
					break
				}
			}
			if !found {
				for _, attr := range a {
					n := h + "." + attr.name
					lp.pch.DeleteCurve(n)
					deleted += sepd + n
					sepd = ","
				}
				delete(lp.ha, h)
			}
		}
	}
	if found != "" {
		lp.log.Info().Msgf("hashes found: %s", found)
	}

	if deleted != "" {
		lp.log.Info().Msgf("hashes deleted: %s", deleted)
	}

	return nil
}

func (t *TrackingEngine) runner(ctx context.Context) {

	sfunc := func() {
		if err := t.syncTrackers(); err != nil {
			t.log.Error().Str("errmsg", err.Error()).Msg("reading from redis failed")
		}
	}
	sfunc()
	tick := time.NewTicker(time.Second * time.Duration(envVar("ADASH_INSTANCE_SYNC_INTERVAL", 60)))

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			sfunc()
		}
	}
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
