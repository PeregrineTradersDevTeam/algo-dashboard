package adash

import (
	//exposes "chart"

	"math/rand"
	"time"
)

func (s *Service) NewPnLChart() {
	dir := 1
	now := time.Now().Truncate(time.Hour)
	for i := 0; i < 3600; i++ {
		s.pnl.XValues = append(s.pnl.XValues, now)
		now = now.Add(1 * time.Minute)

		s.pnl.YValues = append(s.pnl.YValues, float64(rand.Intn(500)*dir))
		if i%100 == 0 {
			z1 := rand.Intn(1)
			if z1 == 0 {
				dir = 1
			} else {
				dir = -1
			}
		}
	}

}
