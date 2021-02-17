package owntracks

import (
	"sync"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
)

type Bind struct {
	di.ConfigBind
	di.MetricsBind
	di.MQTTBind
	di.WorkersBind

	lat     *atomic.Float32Null
	lon     *atomic.Float32Null
	geoHash *atomic.String
	conn    *atomic.String
	acc     *atomic.Int32Null
	alt     *atomic.Int32Null
	batt    *atomic.Float32Null
	vel     *atomic.Int32Null

	mutex    sync.RWMutex
	regions  map[string]Point
	checkers map[string]*atomic.BoolNull
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	b.lat.Nil()
	b.lon.Nil()
	b.geoHash.Set("")
	b.conn.Set("")
	b.acc.Nil()
	b.alt.Nil()
	b.batt.Nil()
	b.vel.Nil()

	b.regions = make(map[string]Point)
	b.checkers = make(map[string]*atomic.BoolNull)

	cfg := b.config()

	for name, region := range cfg.Regions {
		if region.Radius <= 0 {
			region.Radius = cfg.PointRadius
		}

		cfg.Regions[name] = region
		b.registerRegion(name, region)
	}

	return nil
}

func (b *Bind) validAccuracy(acc *int64, maxAccuracy int64) bool {
	if acc == nil {
		return false
	}

	value := *acc
	if value == 0 {
		return false
	}

	if maxAccuracy > 0 && value > maxAccuracy {
		return false
	}

	return true
}

func (b *Bind) getAllRegions() map[string]Point {
	b.mutex.RLock()
	all := make(map[string]Point, len(b.regions))

	for k, v := range b.regions {
		all[k] = v
	}

	b.mutex.RUnlock()

	return all
}

func (b *Bind) getAllRegionCheckers() map[string]*atomic.BoolNull {
	b.mutex.RLock()
	all := make(map[string]*atomic.BoolNull, len(b.checkers))

	for k, v := range b.checkers {
		all[k] = v
	}

	b.mutex.RUnlock()

	return all
}

func (b *Bind) getRegionChecker(name string) (*atomic.BoolNull, bool) {
	b.mutex.RLock()
	check, ok := b.checkers[name]
	b.mutex.RUnlock()

	return check, ok
}

func (b *Bind) getOrSetRegionChecker(name string) *atomic.BoolNull {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if check, ok := b.checkers[name]; ok {
		return check
	}

	b.checkers[name] = atomic.NewBoolNull()

	return b.checkers[name]
}

func (b *Bind) registerRegion(name string, region Point) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.regions[name] = region

	if _, ok := b.checkers[name]; !ok {
		b.checkers[name] = atomic.NewBoolNull()
	}
}
