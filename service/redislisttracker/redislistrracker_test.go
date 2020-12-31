package redislisttracker

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"bitbucket.org/peregrinetraders/mc/pkg/store"
	"github.com/rs/zerolog"
)

var (
	zl        zerolog.Logger
	lps       *RedisListTracker
	instances = []string{"A001", "A002", "A003", "A004", "A005", "A006"}
)

func TestMain(t *testing.M) {

	zl = zerolog.New(os.Stdout)
	con := store.NewRedisStore("localhost:6379", 0, &zl)
	if err := con.Connect(); err != nil {
		fmt.Printf("redis connection failed: %s", err.Error())
		os.Exit(2)
	}

	lps = New(con, &zl)
	//lps.Start(context.Background())

	os.Exit(t.Run())
}

func TestLastPricer_saveLastPrice(t *testing.T) {

	rand.Seed(time.Now().Unix())
	save := []string{redisSetName}
	at := time.Now().Unix()
	for i := 0; i < 100; i++ {
		price := float64(rand.Intn(10)) + float64(rand.Intn(100))/100.0
		save = append(save, lpFormat(instances[rand.Intn(6)], at+int64(i), fmt.Sprintf("%0.f", price)))
	}

	if err := lps.store.Do("RPUSH", save...); err != nil {
		t.Error()
	}
}
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
