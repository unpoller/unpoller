// Package otelunifi provides the methods to turn UniFi measurements into
// OpenTelemetry metrics and export them via OTLP.
package otelunifi

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"golift.io/cnfg"

	"github.com/unpoller/unpoller/pkg/poller"
	"github.com/unpoller/unpoller/pkg/webserver"
)

// PluginName is the name of this plugin.
const PluginName = "otel"

const (
	defaultInterval    = 30 * time.Second
	minimumInterval    = 10 * time.Second
	defaultOTLPHTTPURL = "http://localhost:4318"
	defaultOTLPGRPCURL = "localhost:4317"
	protoHTTP          = "http"
	protoGRPC          = "grpc"
)

// Config defines the data needed to export metrics via OpenTelemetry.
type Config struct {
	// URL is the OTLP endpoint to send metrics to.
	// For HTTP: http://localhost:4318
	// For gRPC: localhost:4317
	URL string `json:"url,omitempty" toml:"url,omitempty" xml:"url" yaml:"url"`

	// APIKey is an optional bearer token / API key for authentication.
	// Sent as the "Authorization: Bearer <key>" header.
	APIKey string `json:"api_key,omitempty" toml:"api_key,omitempty" xml:"api_key" yaml:"api_key"`

	// Interval controls the push interval for sending metrics to the OTLP endpoint.
	Interval cnfg.Duration `json:"interval,omitempty" toml:"interval,omitempty" xml:"interval" yaml:"interval"`

	// Timeout is the per-export deadline.
	Timeout cnfg.Duration `json:"timeout,omitempty" toml:"timeout,omitempty" xml:"timeout" yaml:"timeout"`

	// Protocol selects the OTLP transport protocol: "http" (default) or "grpc".
	Protocol string `json:"protocol,omitempty" toml:"protocol,omitempty" xml:"protocol" yaml:"protocol"`

	// Disable when true disables this output plugin.
	Disable bool `json:"disable" toml:"disable" xml:"disable,attr" yaml:"disable"`

	// DeadPorts when true will save data for dead ports, for example ports that are down or disabled.
	DeadPorts bool `json:"dead_ports" toml:"dead_ports" xml:"dead_ports" yaml:"dead_ports"`
}

// OtelUnifi wraps the config for nested TOML/JSON/YAML config file support.
type OtelUnifi struct {
	*Config `json:"otel" toml:"otel" xml:"otel" yaml:"otel"`
}

// OtelOutput is the working struct for this plugin.
type OtelOutput struct {
	Collector  poller.Collect
	LastCheck  time.Time
	provider   *sdkmetric.MeterProvider
	*OtelUnifi
}

var _ poller.OutputPlugin = &OtelOutput{}

func init() { //nolint:gochecknoinits
	u := &OtelOutput{OtelUnifi: &OtelUnifi{Config: &Config{}}, LastCheck: time.Now()}

	poller.NewOutput(&poller.Output{
		Name:         PluginName,
		Config:       u.OtelUnifi,
		OutputPlugin: u,
	})
}

// Enabled returns true when the plugin is configured and not disabled.
func (u *OtelOutput) Enabled() bool {
	if u == nil {
		return false
	}

	if u.Config == nil {
		return false
	}

	return !u.Disable
}

// DebugOutput validates the plugin configuration without starting the run loop.
func (u *OtelOutput) DebugOutput() (bool, error) {
	if u == nil {
		return true, nil
	}

	if !u.Enabled() {
		return true, nil
	}

	u.setConfigDefaults()

	if u.URL == "" {
		return false, fmt.Errorf("otel: URL must be set")
	}

	proto := u.Protocol
	if proto != protoHTTP && proto != protoGRPC {
		return false, fmt.Errorf("otel: protocol must be %q or %q, got %q", protoHTTP, protoGRPC, proto)
	}

	return true, nil
}

// Run is the main loop called by the poller core.
func (u *OtelOutput) Run(c poller.Collect) error {
	u.Collector = c

	if !u.Enabled() {
		u.LogDebugf("OTel config missing (or disabled), OTel output disabled!")

		return nil
	}

	u.Logf("OpenTelemetry (OTel) output plugin enabled")
	u.setConfigDefaults()

	if err := u.setupProvider(); err != nil {
		return fmt.Errorf("otel: setup provider: %w", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := u.provider.Shutdown(ctx); err != nil {
			u.LogErrorf("otel: shutdown provider: %v", err)
		}
	}()

	webserver.UpdateOutput(&webserver.Output{Name: PluginName, Config: u.Config})
	u.pollController()

	return nil
}

// pollController runs the ticker loop, pushing metrics on each tick.
func (u *OtelOutput) pollController() {
	interval := u.Interval.Duration.Round(time.Second)
	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	u.Logf("OTel->OTLP started, protocol: %s, interval: %v, url: %s",
		u.Protocol, interval, u.URL)

	for u.LastCheck = range ticker.C {
		u.poll(interval)
	}
}

// poll fetches metrics once and sends them to the OTLP endpoint.
func (u *OtelOutput) poll(interval time.Duration) {
	metrics, err := u.Collector.Metrics(&poller.Filter{Name: "unifi"})
	if err != nil {
		u.LogErrorf("metric fetch for OTel failed: %v", err)

		return
	}

	events, err := u.Collector.Events(&poller.Filter{Name: "unifi", Dur: interval})
	if err != nil {
		u.LogErrorf("event fetch for OTel failed: %v", err)

		return
	}

	report, err := u.reportMetrics(metrics, events)
	if err != nil {
		u.LogErrorf("otel report: %v", err)

		return
	}

	u.Logf("OTel Metrics Exported. %v", report)
}

// setupProvider creates and registers the OTel MeterProvider with an OTLP exporter.
func (u *OtelOutput) setupProvider() error {
	ctx := context.Background()

	exp, err := u.buildExporter(ctx)
	if err != nil {
		return fmt.Errorf("building exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(poller.AppName),
		),
	)
	if err != nil {
		return fmt.Errorf("building resource: %w", err)
	}

	u.provider = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(exp,
				sdkmetric.WithInterval(u.Interval.Duration),
				sdkmetric.WithTimeout(u.Timeout.Duration),
			),
		),
		sdkmetric.WithResource(res),
	)

	otel.SetMeterProvider(u.provider)

	return nil
}

// buildExporter creates either an HTTP or gRPC OTLP exporter.
func (u *OtelOutput) buildExporter(ctx context.Context) (sdkmetric.Exporter, error) {
	switch u.Protocol {
	case protoGRPC:
		opts := []otlpmetricgrpc.Option{
			otlpmetricgrpc.WithEndpoint(u.URL),
			otlpmetricgrpc.WithInsecure(),
		}

		if u.APIKey != "" {
			opts = append(opts, otlpmetricgrpc.WithHeaders(map[string]string{
				"Authorization": "Bearer " + u.APIKey,
			}))
		}

		exp, err := otlpmetricgrpc.New(ctx, opts...)
		if err != nil {
			return nil, fmt.Errorf("grpc exporter: %w", err)
		}

		return exp, nil

	default: // http
		opts := []otlpmetrichttp.Option{
			otlpmetrichttp.WithEndpoint(u.URL),
			otlpmetrichttp.WithInsecure(),
		}

		if u.APIKey != "" {
			opts = append(opts, otlpmetrichttp.WithHeaders(map[string]string{
				"Authorization": "Bearer " + u.APIKey,
			}))
		}

		exp, err := otlpmetrichttp.New(ctx, opts...)
		if err != nil {
			return nil, fmt.Errorf("http exporter: %w", err)
		}

		return exp, nil
	}
}

// setConfigDefaults fills in zero-value fields with sensible defaults.
func (u *OtelOutput) setConfigDefaults() {
	if u.Protocol == "" {
		u.Protocol = protoHTTP
	}

	if u.URL == "" {
		switch u.Protocol {
		case protoGRPC:
			u.URL = defaultOTLPGRPCURL
		default:
			u.URL = defaultOTLPHTTPURL
		}
	}

	if u.Interval.Duration == 0 {
		u.Interval = cnfg.Duration{Duration: defaultInterval}
	} else if u.Interval.Duration < minimumInterval {
		u.Interval = cnfg.Duration{Duration: minimumInterval}
	}

	u.Interval = cnfg.Duration{Duration: u.Interval.Duration.Round(time.Second)}

	if u.Timeout.Duration == 0 {
		u.Timeout = cnfg.Duration{Duration: 10 * time.Second}
	}
}
