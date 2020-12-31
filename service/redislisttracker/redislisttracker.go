package redislisttracker

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

type RedisListTracker struct {
	store store.Storer
	log   zerolog.Logger

	pch *charter.Charter

	mux sync.RWMutex

	src  map[string]attr
	list map[string]attr
}

// New returns RedisListTracker instance.
func New(s store.Storer, logger *zerolog.Logger, pch *charter.Charter) *RedisListTracker {
	return &RedisListTracker{
		store: s,
		log:   logger.With().Str("layer", "list-tracker").Logger(),
		src:   make(map[string]attr),
		list:  make(map[string]attr),
		pch:   pch}
}

// TrackList.
func (t *RedisListTracker) TrackList(
	listprefix string,
	cf charter.StyleFunc,
	pf charter.PointFunc,
	labelFunc ...charter.LabelFunc) {

	t.mux.Lock()
	defer t.mux.Unlock()

	t.src[listprefix] = attr{
		name:       listprefix,
		labelfunc:  append([]charter.LabelFunc{}, labelFunc...),
		stylefunc:  cf,
		listprefix: listprefix,
		pointFunc:  pf,
	}
}

func (t *RedisListTracker) Start(ctx context.Context) error {
	go t.runner(ctx)
	return nil
}

func (lp *RedisListTracker) syncTrackers() error {

	found := ""
	deleted := ""
	sepf := ""
	sepd := ""
	for p, a := range lp.src {
		var keys []string
		if err := lp.store.Scan(p+":*", &keys); err != nil {
			return err
		}

		// add if new hash appeared
		for i := range keys {
			if _, ok := lp.list[keys[i]]; ok {
				continue
			}

			found += sepf + keys[i]
			sepf = ","
			lp.list[keys[i]] = a
			lp.pch.RegisterCurve(charter.NewCurve(
				keys[i],
				"",
				charter.RedisPnLReader(&lp.log, lp.store, keys[i]),
				a.stylefunc,
				a.pointFunc,
				a.labelfunc...,
			),
			)

		}

		// del if hash dissapeared. can be optimized
		for h := range lp.list {
			found := false
			for i := range keys {
				if keys[i] == h {
					found = true
					break
				}
			}
			if !found {
				lp.pch.DeleteCurve(h)
				delete(lp.list, h)
				deleted += sepd + h
				sepd = ","
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

func (t *RedisListTracker) runner(ctx context.Context) {

	sfunc := func() {
		if err := t.syncTrackers(); err != nil {
			t.log.Error().Str("errmsg", err.Error()).Msg("reading from redis failed")
		}
	}
	sfunc()
	tick := time.NewTicker(time.Second * time.Duration(envVar("ADASH_PNL_SYNC_INTERVAL", 5)))

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
