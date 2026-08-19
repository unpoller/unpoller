package inputunas

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/unpoller/unifi/v5"
	"github.com/unpoller/unpoller/pkg/poller"
	"github.com/unpoller/unpoller/pkg/webserver"
)

// ErrNoDevices is returned by DebugInput when the plugin is enabled but nothing is configured.
var ErrNoDevices = errors.New("no UNAS devices configured")

/* This file contains the poller.Input interface methods. */

// Initialize gets called one time when starting up.
// Satisfies poller.Input interface.
func (u *InputUNAS) Initialize(l poller.Logger) error {
	if u.Config == nil {
		u.Config = &Config{}
	}

	u.Logger = l

	// enable defaults to false, so the plugin is off until an operator asks for it. Warn when
	// consoles are configured but the switch was never flipped -- that combination is always a
	// mistake, and silence would leave the operator hunting for missing metrics. An operator
	// who has never heard of UNAS has no devices either, and still sees nothing.
	if !u.Enable {
		if len(u.configuredDevices()) > 0 {
			u.LogErrorf("UNAS devices are configured but unas.enable is false; not polling them")
		}

		return nil
	}

	u.Devices = u.configuredDevices()

	// Enabled with nothing to poll: nothing to say, and nothing to do.
	if len(u.Devices) == 0 {
		return nil
	}

	for i, d := range u.Devices {
		if err := u.login(d); err != nil {
			// Not fatal: a console that is down at startup gets another attempt every poll.
			u.LogErrorf("UNAS console %d of %d auth or connection error, retrying next poll: %v",
				i+1, len(u.Devices), err)

			continue
		}

		u.Logf("Configured UNAS console %d of %d:", i+1, len(u.Devices))
		u.logDevice(d)
	}

	webserver.UpdateInput(&webserver.Input{Name: PluginName, Config: formatConfig(u.Config)})

	return nil
}

// configuredDevices returns the devices with defaults applied, dropping any without a URL.
// A device with no URL is a config typo (or an empty [[unas.device]] table); polling it
// would only produce a confusing error against the empty string every cycle.
func (u *InputUNAS) configuredDevices() []*Device {
	devices := make([]*Device, 0, len(u.Devices))

	for _, d := range u.Devices {
		if d == nil {
			continue
		}

		if d.URL == "" {
			u.LogErrorf("Ignoring UNAS device with no url configured.")

			continue
		}

		devices = append(devices, u.setDefaults(d))
	}

	return devices
}

func (u *InputUNAS) logDevice(d *Device) {
	u.Logf("   => URL: %s (verify SSL: %v, timeout: %v)", d.URL, *d.VerifySSL, d.Timeout.Duration)

	if len(d.CertPaths) > 0 {
		u.Logf("   => Cert Files: %s", strings.Join(d.CertPaths, ", "))
	}

	u.Logf("   => Username: %s (has password: %v)", d.User, d.Pass != "")
}

// Metrics polls every configured UNAS console and returns the aggregate.
//
// A per-console failure is logged and skipped rather than returned, and an error comes back
// only when nothing at all was collected. This is not politeness: poller.collectMetrics
// discards the metrics entirely when an input returns both a result and an error, so
// returning `(metrics, err)` after one console of three fails would throw away the two that
// worked. Satisfies poller.Input interface.
func (u *InputUNAS) Metrics(filter *poller.Filter) (*poller.Metrics, error) {
	if !u.Enable {
		return nil, nil
	}

	metrics := &poller.Metrics{TS: time.Now()}

	if filter == nil {
		filter = &poller.Filter{}
	}

	var errs []error

	for _, d := range u.Devices {
		if filter.Path != "" && !strings.EqualFold(d.URL, filter.Path) {
			continue
		}

		device, err := u.collectDevice(d)
		if err != nil {
			u.LogErrorf("Failed to collect metrics from UNAS console %s: %v", d.URL, err)
			errs = append(errs, fmt.Errorf("%s: %w", d.URL, err))

			continue
		}

		metrics.UNASDevices = append(metrics.UNASDevices, device)
	}

	if len(metrics.UNASDevices) == 0 && len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return metrics, nil
}

// collectDevice fetches one console, re-authenticating once if the whole fetch failed.
//
// UNAS session cookies expire after roughly two hours, so a long-running poller must expect
// to be logged out mid-life and must re-attach fresh cookies rather than reuse the dead jar.
// We retry on any total failure rather than on a 401 specifically: a mid-session 401 from
// GetData surfaces as ErrInvalidStatusCode, not ErrAuthenticationFailed, so there is no
// sentinel to match on. Session expiry fails every endpoint at once, which is exactly the
// total-failure case, and one wasted re-login against a console that is simply unreachable
// is cheaper than never recovering from an expired session. Partial failures never reach
// here: GetUNASDevice logs those endpoints and leaves their fields nil.
func (u *InputUNAS) collectDevice(d *Device) (*unifi.UNASDevice, error) {
	device, err := u.getUNASDevice(d)
	if err == nil {
		return device, nil
	}

	u.Logf("Re-authenticating to UNAS console: %s", d.URL)

	if lerr := u.login(d); lerr != nil {
		// Report the original failure too, so a console that is merely down does not look
		// like a credentials problem.
		return nil, fmt.Errorf("re-authenticating after %v: %w", err, lerr)
	}

	return u.getUNASDevice(d)
}

// getUNASDevice reads the client under the read lock and releases it before polling, so a
// slow console cannot block a concurrent re-auth (which takes the write lock).
func (u *InputUNAS) getUNASDevice(d *Device) (*unifi.UNASDevice, error) {
	u.RLock()

	client := d.unas

	u.RUnlock()

	if client == nil {
		return nil, unifi.ErrNilUnifi
	}

	return client.GetUNASDevice()
}

// Events is a no-op: a UNAS console exposes no event or log endpoints in this version.
// Satisfies poller.Input interface.
func (u *InputUNAS) Events(_ *poller.Filter) (*poller.Events, error) {
	return &poller.Events{}, nil
}

// RawMetrics returns the raw JSON from one UNAS endpoint, selected by filter.Kind.
// Adjust filter.Unit to pull from a console other than the first.
// Satisfies poller.Input interface.
func (u *InputUNAS) RawMetrics(filter *poller.Filter) ([]byte, error) {
	if filter == nil {
		filter = &poller.Filter{}
	}

	if l := len(u.Devices); filter.Unit >= l {
		return nil, fmt.Errorf("%d UNAS console(s) configured, '%d': %w", l, filter.Unit, ErrNoDevices)
	}

	d := u.Devices[filter.Unit]

	u.RLock()

	client := d.unas

	u.RUnlock()

	if client == nil {
		if err := u.login(d); err != nil {
			return nil, err
		}

		u.RLock()

		client = d.unas

		u.RUnlock()
	}

	var path string

	switch filter.Kind {
	case "", "device", "device-info", "d":
		path = unifi.APIUNASDeviceInfoPath
	case "storage", "s":
		path = unifi.APIUNASStoragePath
	case "drives", "dr":
		path = unifi.APIUNASDrivesPath
	case "network-io", "n":
		path = unifi.APIUNASNetworkIOPath
	default:
		return nil, fmt.Errorf("must provide filter: device-info, storage, drives, network-io: %w",
			unifi.ErrEndpointNotFound)
	}

	return client.GetJSON(path)
}

// DebugInput checks that every configured console can be reached and authenticated against.
// Satisfies poller.Input interface.
func (u *InputUNAS) DebugInput() (bool, error) {
	if u == nil || u.Config == nil || !u.Enable {
		return true, nil
	}

	// Safe to assign without holding the lock: --debugio is a one-shot mode that exits before
	// Run() ever starts polling (see poller.Start), so this cannot race a Metrics call. It is
	// not redundant with Initialize either -- that never runs in this mode.
	u.Devices = u.configuredDevices()

	if len(u.Devices) == 0 {
		return true, nil
	}

	allOK := true

	var errs []error

	for i, d := range u.Devices {
		if err := u.login(d); err != nil {
			u.LogErrorf("UNAS console %d of %d auth or connection error: %v", i+1, len(u.Devices), err)

			allOK = false

			errs = append(errs, err)

			continue
		}

		u.Logf("Valid UNAS console %d of %d:", i+1, len(u.Devices))
		u.logDevice(d)
	}

	return allOK, errors.Join(errs...)
}

// formatConfig copies the config for display on the web interface, replacing the password
// with whether one is set.
func formatConfig(config *Config) *Config {
	devices := make([]*Device, len(config.Devices))

	for i, d := range config.Devices {
		devices[i] = &Device{
			URL:       d.URL,
			User:      d.User,
			Pass:      strconv.FormatBool(d.Pass != ""),
			VerifySSL: d.VerifySSL,
			CertPaths: d.CertPaths,
			Timeout:   d.Timeout,
		}
	}

	return &Config{
		Default: Device{
			URL:       config.Default.URL,
			User:      config.Default.User,
			Pass:      strconv.FormatBool(config.Default.Pass != ""),
			VerifySSL: config.Default.VerifySSL,
			CertPaths: config.Default.CertPaths,
			Timeout:   config.Default.Timeout,
		},
		Enable:  config.Enable,
		Devices: devices,
	}
}

// Logf logs a message.
func (u *InputUNAS) Logf(msg string, v ...any) {
	webserver.NewInputEvent(PluginName, PluginName, &webserver.Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "info"},
	})

	if u.Logger != nil {
		u.Logger.Logf(msg, v...)
	}
}

// LogErrorf logs an error message.
func (u *InputUNAS) LogErrorf(msg string, v ...any) {
	webserver.NewInputEvent(PluginName, PluginName, &webserver.Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "error"},
	})

	if u.Logger != nil {
		u.Logger.LogErrorf(msg, v...)
	}
}

// LogDebugf logs a debug message.
func (u *InputUNAS) LogDebugf(msg string, v ...any) {
	webserver.NewInputEvent(PluginName, PluginName, &webserver.Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "debug"},
	})

	if u.Logger != nil {
		u.Logger.LogDebugf(msg, v...)
	}
}
