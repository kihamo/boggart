package di

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

const (
	metricsSizeOfChannel = 100
)

type MetricsContainerSupport interface {
	SetMetrics(*MetricsContainer)
	Metrics() *MetricsContainer
}

func MetricsContainerBind(bind boggart.Bind) (*MetricsContainer, bool) {
	if support, ok := bind.(MetricsContainerSupport); ok {
		container := support.Metrics()
		return container, container != nil
	}

	return nil, false
}

type MetricsBind struct {
	mutex     sync.RWMutex
	container *MetricsContainer
}

func (b *MetricsBind) SetMetrics(container *MetricsContainer) {
	b.mutex.Lock()
	b.container = container
	b.mutex.Unlock()
}

func (b *MetricsBind) Metrics() *MetricsContainer {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.container
}

type MetricsContainer struct {
	snitch.Collector

	bind              boggart.BindItem
	collector         snitch.Collector
	prefix            string
	descriptionsCount uint64
	collectCount      uint64
	emptyCount        uint64
}

func NewMetricsContainer(bind boggart.BindItem) *MetricsContainer {
	c := &MetricsContainer{
		bind:   bind,
		prefix: boggart.ComponentName + "_bind_" + strings.ReplaceAll(bind.Type(), ":", "_") + "_",
	}

	if collector, ok := c.bind.Bind().(snitch.Collector); ok {
		c.collector = collector
	}

	return c
}

func (c *MetricsContainer) description(desc *snitch.Description) *snitch.Description {
	labels := make([]string, 0, len(desc.Labels()))
	for _, l := range desc.Labels() {
		labels = append(labels, l.Key, l.Value)
	}

	return snitch.NewDescription(c.prefix+desc.Name(), desc.Help(), desc.Type(), labels...)
}

func (c *MetricsContainer) Describe(ch chan<- *snitch.Description) {
	status := c.bind.Status()
	if !status.IsStatusOnline() && !status.IsStatusOffline() || c.collector == nil {
		return
	}

	tmpCh := make(chan *snitch.Description, metricsSizeOfChannel)
	var count uint64

	go func() {
		c.collector.Describe(tmpCh)
		close(tmpCh)
	}()

	for v := range tmpCh {
		count++
		ch <- c.description(v)
	}

	atomic.StoreUint64(&c.descriptionsCount, count)
}

func (c *MetricsContainer) Collect(ch chan<- snitch.Metric) {
	status := c.bind.Status()
	if !status.IsStatusOnline() && !status.IsStatusOffline() || c.collector == nil {
		return
	}

	tmpCh := make(chan snitch.Metric, metricsSizeOfChannel)
	var (
		count      uint64
		emptyCount uint64
	)

	go func() {
		c.collector.Collect(tmpCh)
		close(tmpCh)
	}()

	for v := range tmpCh {
		count++

		if m, err := v.Measure(); err == nil && m.SampleCount != nil && *m.SampleCount == 0 {
			emptyCount++
		}

		ch <- &metricWrapper{
			metric:    v,
			container: c,
		}
	}

	atomic.StoreUint64(&c.collectCount, count)
	atomic.StoreUint64(&c.emptyCount, emptyCount)
}

func (c *MetricsContainer) DescriptionsCount() (count uint64) {
	if c.collector == nil {
		return 0
	}

	count = atomic.LoadUint64(&c.descriptionsCount)
	if count == 0 {
		tmpCh := make(chan *snitch.Description, metricsSizeOfChannel)

		go func() {
			c.collector.Describe(tmpCh)
			close(tmpCh)
		}()

		for range tmpCh {
			count++
		}

		atomic.StoreUint64(&c.descriptionsCount, count)
	}

	return atomic.LoadUint64(&c.descriptionsCount)
}

func (c *MetricsContainer) CollectCount() uint64 {
	return atomic.LoadUint64(&c.collectCount)
}

func (c *MetricsContainer) EmptyCount() uint64 {
	return atomic.LoadUint64(&c.emptyCount)
}

func (c *MetricsContainer) Gather() (snitch.Measures, error) {
	tmpCh := make(chan snitch.Metric, metricsSizeOfChannel)
	measures := make(snitch.Measures, 0)
	var (
		count      uint64
		emptyCount uint64
	)

	go func() {
		c.Collect(tmpCh)
		close(tmpCh)
	}()

	for metric := range tmpCh {
		count++

		value, err := metric.Measure()
		if err != nil {
			return nil, err
		}

		if value.SampleCount != nil && *value.SampleCount == 0 {
			emptyCount++
		}

		measures = append(measures, &snitch.Measure{
			Description: metric.Description(),
			CreatedAt:   time.Now(),
			Value:       value,
		})
	}

	atomic.StoreUint64(&c.descriptionsCount, count)
	atomic.StoreUint64(&c.emptyCount, emptyCount)

	return measures, nil
}

type metricWrapper struct {
	metric    snitch.Metric
	container *MetricsContainer
}

func (m *metricWrapper) Description() *snitch.Description {
	return m.container.description(m.metric.Description())
}

func (m *metricWrapper) Measure() (*snitch.MeasureValue, error) {
	return m.metric.Measure()
}
