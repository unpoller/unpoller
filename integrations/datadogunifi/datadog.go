// Package datadogunifi provides the methods to turn UniFi measurements into Datadog
// data points with appropriate tags and fields.
package datadogunifi

import (
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/unpoller/poller"
	"github.com/unpoller/unifi"
	"go.uber.org/zap"
	"golift.io/cnfg"
)

const (
	defaultInterval = 30 * time.Second
	minimumInterval = 10 * time.Second
)

// Config defines the data needed to store metrics in Datadog.
type Config struct {
	// Required Config

	// Interval controls the collection and reporting interval
	Interval cnfg.Duration `json:"interval,omitempty" toml:"interval,omitempty" xml:"interval,omitempty" yaml:"interval,omitempty"`

	// Save data for dead ports? ie. ports that are down or disabled.
	DeadPorts bool `json:"dead_ports,omitempty" toml:"dead_ports,omitempty" xml:"dead_ports,omitempty" yaml:"dead_ports,omitempty"`

	// Disable when true disables this output plugin
	Disable bool `json:"disable" toml:"disable" xml:"disable,attr" yaml:"disable"`
	// Address determines how to talk to the Datadog agent
	Address string `json:"address" toml:"address" xml:"address,attr" yaml:"address"`

	// Optional Statsd Options - mirrored from statsd.Options

	// Namespace to prepend to all metrics, events and service checks name.
	Namespace *string `json:"namespace" toml:"namespace" xml:"namespace,attr" yaml:"namespace"`

	// Tags are global tags to be applied to every metrics, events and service checks.
	Tags []string `json:"tags" toml:"tags" xml:"tags,attr" yaml:"tags"`

	// MaxBytesPerPayload is the maximum number of bytes a single payload will contain.
	// The magic value 0 will set the option to the optimal size for the transport
	// protocol used when creating the client: 1432 for UDP and 8192 for UDS.
	MaxBytesPerPayload *int `json:"max_bytes_per_payload" toml:"max_bytes_per_payload" xml:"max_bytes_per_payload,attr" yaml:"max_bytes_per_payload"`

	// MaxMessagesPerPayload is the maximum number of metrics, events and/or service checks a single payload will contain.
	// This option can be set to `1` to create an unbuffered client.
	MaxMessagesPerPayload *int `json:"max_messages_per_payload" toml:"max_messages_per_payload" xml:"max_messages_per_payload,attr" yaml:"max_messages_per_payload"`

	// BufferPoolSize is the size of the pool of buffers in number of buffers.
	// The magic value 0 will set the option to the optimal size for the transport
	// protocol used when creating the client: 2048 for UDP and 512 for UDS.
	BufferPoolSize *int `json:"buffer_pool_size" toml:"buffer_pool_size" xml:"buffer_pool_size,attr" yaml:"buffer_pool_size"`

	// BufferFlushInterval is the interval after which the current buffer will get flushed.
	BufferFlushInterval *cnfg.Duration `json:"buffer_flush_interval" toml:"buffer_flush_interval" xml:"buffer_flush_interval,attr" yaml:"buffer_flush_interval"`

	// BufferShardCount is the number of buffer "shards" that will be used.
	// Those shards allows the use of multiple buffers at the same time to reduce
	// lock contention.
	BufferShardCount *int `json:"buffer_shard_count" toml:"buffer_shard_count" xml:"buffer_shard_count,attr" yaml:"buffer_shard_count"`

	// SenderQueueSize is the size of the sender queue in number of buffers.
	// The magic value 0 will set the option to the optimal size for the transport
	// protocol used when creating the client: 2048 for UDP and 512 for UDS.
	SenderQueueSize *int `json:"sender_queue_size" toml:"sender_queue_size" xml:"sender_queue_size,attr" yaml:"sender_queue_size"`

	// WriteTimeoutUDS is the timeout after which a UDS packet is dropped.
	WriteTimeoutUDS *cnfg.Duration `json:"write_timeout_uds" toml:"write_timeout_uds" xml:"write_timeout_uds,attr" yaml:"write_timeout_uds"`

	// ReceiveMode determines the behavior of the client when receiving to many
	// metrics. The client will either drop the metrics if its buffers are
	// full (ChannelMode mode) or block the caller until the metric can be
	// handled (MutexMode mode). By default the client will MutexMode. This
	// option should be set to ChannelMode only when use under very high
	// load.
	//
	// MutexMode uses a mutex internally which is much faster than
	// channel but causes some lock contention when used with a high number
	// of threads. Mutex are sharded based on the metrics name which
	// limit mutex contention when goroutines send different metrics.
	//
	// ChannelMode: uses channel (of ChannelModeBufferSize size) to send
	// metrics and drop metrics if the channel is full. Sending metrics in
	// this mode is slower that MutexMode (because of the channel), but
	// will not block the application. This mode is made for application
	// using many goroutines, sending the same metrics at a very high
	// volume. The goal is to not slow down the application at the cost of
	// dropping metrics and having a lower max throughput.
	ReceiveMode *statsd.ReceivingMode `json:"receive_mode" toml:"receive_mode" xml:"receive_mode,attr" yaml:"receive_mode"`

	// ChannelModeBufferSize is the size of the channel holding incoming metrics
	ChannelModeBufferSize *int `json:"channel_mode_buffer_size" toml:"channel_mode_buffer_size" xml:"channel_mode_buffer_size,attr" yaml:"channel_mode_buffer_size"`

	// AggregationFlushInterval is the interval for the aggregator to flush metrics
	AggregationFlushInterval *time.Duration `json:"aggregation_flush_interval" toml:"aggregation_flush_interval" xml:"aggregation_flush_interval,attr" yaml:"aggregation_flush_interval"`
}

// Datadog allows the data to be context aware with configuration
type Datadog struct {
	*Config `json:"datadog" toml:"datadog" xml:"datadog" yaml:"datadog"`
	options []statsd.Option
}

// DatadogUnifi is returned by New() after you provide a Config.
type DatadogUnifi struct {
	Collector poller.Collect
	datadog   statsd.ClientInterface
	LastCheck time.Time
	Logger    *zap.SugaredLogger
	*Datadog
}

func init() { // nolint: gochecknoinits
	l, _ := zap.NewProduction()
	u := &DatadogUnifi{Datadog: &Datadog{}, LastCheck: time.Now(), Logger: l.Sugar()}

	poller.NewOutput(&poller.Output{
		Name:   "datadog",
		Config: u.Datadog,
		Method: u.Run,
	})
}

func (u *DatadogUnifi) setConfigDefaults() {
	if u.Interval.Duration == 0 {
		u.Interval = cnfg.Duration{Duration: defaultInterval}
	} else if u.Interval.Duration < minimumInterval {
		u.Interval = cnfg.Duration{Duration: minimumInterval}
	}

	u.Interval = cnfg.Duration{Duration: u.Interval.Duration.Round(time.Second)}

	u.options = make([]statsd.Option, 0)

	if u.Namespace != nil {
		u.options = append(u.options, statsd.WithNamespace(*u.Namespace))
	}

	if u.Tags != nil && len(u.Tags) > 0 {
		u.options = append(u.options, statsd.WithTags(u.Tags))
	}

	if u.MaxBytesPerPayload != nil {
		u.options = append(u.options, statsd.WithMaxBytesPerPayload(*u.MaxBytesPerPayload))
	}

	if u.MaxMessagesPerPayload != nil {
		u.options = append(u.options, statsd.WithMaxMessagesPerPayload(*u.MaxMessagesPerPayload))
	}

	if u.BufferPoolSize != nil {
		u.options = append(u.options, statsd.WithBufferPoolSize(*u.BufferPoolSize))
	}

	if u.BufferFlushInterval != nil {
		u.options = append(u.options, statsd.WithBufferFlushInterval((*u.BufferFlushInterval).Duration))
	}

	if u.BufferShardCount != nil {
		u.options = append(u.options, statsd.WithBufferShardCount(*u.BufferShardCount))
	}

	if u.SenderQueueSize != nil {
		u.options = append(u.options, statsd.WithSenderQueueSize(*u.SenderQueueSize))
	}

	if u.WriteTimeoutUDS != nil {
		u.options = append(u.options, statsd.WithWriteTimeoutUDS((*u.WriteTimeoutUDS).Duration))
	}

	if u.ReceiveMode != nil {
		switch *u.ReceiveMode {
		case statsd.ChannelMode:
			u.options = append(u.options, statsd.WithChannelMode())
		case statsd.MutexMode:
			u.options = append(u.options, statsd.WithMutexMode())
		}
	}

	if u.ChannelModeBufferSize != nil {
		u.options = append(u.options, statsd.WithChannelModeBufferSize(*u.ChannelModeBufferSize))
	}

	if u.AggregationFlushInterval != nil {
		u.options = append(u.options, statsd.WithAggregationInterval(*u.AggregationFlushInterval))
	}

}

// Run runs a ticker to poll the unifi server and update Datadog.
func (u *DatadogUnifi) Run(c poller.Collect) error {
	defer u.Logger.Sync()
	if u.Disable {
		u.Logger.Debug("Datadog config is disabled, output is disabled.")
		return nil
	}
	if u.Config == nil && !u.Disable {
		u.Logger.Error("DataDog config is missing and is not disabled: Datadog output is disabled!")
		return nil
	}

	u.Collector = c
	u.setConfigDefaults()

	var err error
	u.datadog, err = statsd.New(u.Address, u.options...)
	if err != nil {
		return err
	}

	u.PollController()

	return nil
}

// PollController runs forever, polling UniFi and pushing to Datadog
// This is started by Run() or RunBoth() after everything is validated.
func (u *DatadogUnifi) PollController() {
	interval := u.Interval.Round(time.Second)
	ticker := time.NewTicker(interval)
	u.Logger.Info("Everything checks out! Poller started", zap.Duration("interval", interval))

	for u.LastCheck = range ticker.C {
		metrics, err := u.Collector.Metrics(&poller.Filter{Name: "unifi"})
		if err != nil {
			u.Logger.Error("metric fetch for Datadog failed", zap.Error(err))
			continue
		}

		events, err := u.Collector.Events(&poller.Filter{Name: "unifi", Dur: interval})
		if err != nil {
			u.Logger.Error("event fetch for Datadog failed", zap.Error(err))
			continue
		}

		report, err := u.ReportMetrics(metrics, events)
		if err != nil {
			// Is the agent down?
			u.Logger.Error("unable to report metrics and events", zap.Error(err))
			report.reportCount("unifi.collect.errors", 1, []string{})
			continue
		}
		report.reportCount("unifi.collect.success", 1, []string{})
		u.LogDatadogReport(report)
	}
}

// ReportMetrics batches all device and client data into datadog data points.
// Call this after you've collected all the data you care about.
// Returns an error if datadog statsd calls fail, otherwise returns a report.
func (u *DatadogUnifi) ReportMetrics(m *poller.Metrics, e *poller.Events) (*Report, error) {
	r := &Report{
		Metrics: m,
		Events:  e,
		Start:   time.Now(),
		Counts:  &Counts{Val: make(map[item]int)},
		Logger:  u.Logger,
	}
	// batch all the points.
	u.loopPoints(r)
	r.End = time.Now()
	r.Elapsed = r.End.Sub(r.Start)
	r.reportTiming("unifi.collector_timing", r.Elapsed, []string{})
	return r, nil
}

// loopPoints collects all the data to immediately report to Datadog.
func (u *DatadogUnifi) loopPoints(r report) {
	m := r.metrics()

	for _, s := range m.RogueAPs {
		u.switchExport(r, s)
	}

	for _, s := range m.Sites {
		u.switchExport(r, s)
	}

	for _, s := range m.SitesDPI {
		u.reportSiteDPI(r, s.(*unifi.DPITable))
	}

	for _, s := range m.Clients {
		u.switchExport(r, s)
	}

	for _, s := range m.Devices {
		u.switchExport(r, s)
	}

	for _, s := range r.events().Logs {
		u.switchExport(r, s)
	}

	appTotal := make(totalsDPImap)
	catTotal := make(totalsDPImap)

	for _, s := range m.ClientsDPI {
		u.batchClientDPI(r, s, appTotal, catTotal)
	}

	reportClientDPItotals(r, appTotal, catTotal)
}

func (u *DatadogUnifi) switchExport(r report, v interface{}) { //nolint:cyclop
	switch v := v.(type) {
	case *unifi.RogueAP:
		u.batchRogueAP(r, v)
	case *unifi.UAP:
		u.batchUAP(r, v)
	case *unifi.USW:
		u.batchUSW(r, v)
	case *unifi.USG:
		u.batchUSG(r, v)
	case *unifi.UXG:
		u.batchUXG(r, v)
	case *unifi.UDM:
		u.batchUDM(r, v)
	case *unifi.Site:
		u.reportSite(r, v)
	case *unifi.Client:
		u.batchClient(r, v)
	case *unifi.Event:
		u.batchEvent(r, v)
	case *unifi.IDS:
		u.batchIDS(r, v)
	case *unifi.Alarm:
		u.batchAlarms(r, v)
	case *unifi.Anomaly:
		u.batchAnomaly(r, v)
	default:
		u.Logger.Error("invalid export", zap.Reflect("type", v))
	}
}

// LogDatadogReport writes a log message after exporting to Datadog.
func (u *DatadogUnifi) LogDatadogReport(r *Report) {
	m := r.Metrics
	u.Logger.Info("UniFi Metrics Recorded",
		zap.Int("num_sites", len(m.Sites)),
		zap.Int("num_sites_dpi", len(m.SitesDPI)),
		zap.Int("num_clients", len(m.Clients)),
		zap.Int("num_clients_dpi", len(m.ClientsDPI)),
		zap.Int("num_rogue_ap", len(m.RogueAPs)),
		zap.Int("num_devices", len(m.Devices)),
		zap.Errors("errors", r.Errors),
		zap.Duration("elapsed", r.Elapsed),
	)
}
