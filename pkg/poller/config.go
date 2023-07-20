package poller

import (
	"fmt"
	"os"
	"path"
	"plugin"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"golift.io/cnfg"
	"golift.io/cnfgfile"
)

const (
	// AppName is the name of the application.
	AppName = "unpoller"
	// ENVConfigPrefix is the prefix appended to an env variable tag name.
	ENVConfigPrefix = "UP"
)

// DefaultConfFile is where to find config if --config is not provided.
func DefaultConfFile() string {
	switch runtime.GOOS {
	case "windows":
		return `C:\ProgramData\unifi-poller\up.conf`
	case "darwin":
		fallthrough
	case "freebsd":
		fallthrough
	case "netbsd":
		fallthrough
	case "openbsd":
		return "/etc/unpoller/up.conf,/etc/unifi-poller/up.conf,/usr/local/etc/unifi-poller/up.conf"
	default:
		// linux and everything else
		return "/etc/unpoller/up.conf,/config/unifi-poller.conf,/etc/unifi-poller/up.conf"
	}
}

// DefaultObjPath is the path to look for shared object libraries (plugins).
func DefaultObjPath() string {
	switch runtime.GOOS {
	case "windows":
		// DefaultObjPath is useless in this context. Bummer.
		return "PLUGINS_DO_NOT_WORK_ON_WINDOWS_SOWWWWWY"
	case "darwin":
		fallthrough
	case "freebsd":
		fallthrough
	case "netbsd":
		fallthrough
	case "openbsd":
		return "/usr/local/lib/unpoller"
	default:
		// linux and everything else
		return "/usr/local/lib/unpoller"
	}
}

// UnifiPoller contains the application startup data, and auth info for UniFi & Influx.
type UnifiPoller struct {
	Flags *Flags
	*Config
}

// Flags represents the CLI args available and their settings.
type Flags struct {
	ConfigFile string
	DumpJSON   string
	HashPW     string
	ShowVer    bool
	DebugIO    bool
	*pflag.FlagSet
}

// Metrics is a type shared by the exporting and reporting packages.
type Metrics struct {
	TS         time.Time
	Sites      []any
	Clients    []any
	SitesDPI   []any
	ClientsDPI []any
	Devices    []any
	RogueAPs   []any
}

// Events defines the type for log entries.
type Events struct {
	Logs []any
}

// Config represents the core library input data.
type Config struct {
	*Poller `json:"poller" toml:"poller" xml:"poller" yaml:"poller"`
}

// Poller is the global config values.
type Poller struct {
	Plugins []string `json:"plugins" toml:"plugins" xml:"plugin"     yaml:"plugins"`
	Debug   bool     `json:"debug"   toml:"debug"   xml:"debug,attr" yaml:"debug"`
	Quiet   bool     `json:"quiet"   toml:"quiet"   xml:"quiet,attr" yaml:"quiet"`
}

// LoadPlugins reads-in dynamic shared libraries.
// Not used very often, if at all.
func (u *UnifiPoller) LoadPlugins() error {
	for _, p := range u.Plugins {
		name := strings.TrimSuffix(p, ".so") + ".so"

		if name == ".so" {
			continue // Just ignore it. uhg.
		}

		if _, err := os.Stat(name); os.IsNotExist(err) {
			name = path.Join(DefaultObjPath(), name)
		}

		u.Logf("Loading Dynamic Plugin: %s", name)

		if _, err := plugin.Open(name); err != nil {
			return fmt.Errorf("opening plugin: %w", err)
		}
	}

	return nil
}

// ParseConfigs parses the poller config and the config for each registered output plugin.
func (u *UnifiPoller) ParseConfigs() error {
	// Parse core config.
	if err := u.parseInterface(u.Config); err != nil {
		return err
	}

	// Load dynamic plugins.
	if err := u.LoadPlugins(); err != nil {
		return err
	}

	if err := u.parseInputs(); err != nil {
		return err
	}

	return u.parseOutputs()
}

// getFirstFile returns the first file that exists and is "reachable".
func getFirstFile(files []string) (string, error) {
	var err error

	for _, f := range files {
		if _, err = os.Stat(f); err == nil {
			return f, nil
		}
	}

	return "", fmt.Errorf("finding file: %w", err)
}

// parseInterface parses the config file and environment variables into the provided interface.
func (u *UnifiPoller) parseInterface(i any) error {
	// Parse config file into provided interface.
	if err := cnfgfile.Unmarshal(i, u.Flags.ConfigFile); err != nil {
		return fmt.Errorf("cnfg unmarshal: %w", err)
	}

	// Parse environment variables into provided interface.
	if _, err := cnfg.UnmarshalENV(i, ENVConfigPrefix); err != nil {
		return fmt.Errorf("env unmarshal: %w", err)
	}

	return nil
}

// Parse input plugin configs.
func (u *UnifiPoller) parseInputs() error {
	inputSync.Lock()
	defer inputSync.Unlock()

	for _, i := range inputs {
		if err := u.parseInterface(i.Config); err != nil {
			return err
		}
	}

	return nil
}

// Parse output plugin configs.
func (u *UnifiPoller) parseOutputs() error {
	outputSync.Lock()
	defer outputSync.Unlock()

	for _, o := range outputs {
		if err := u.parseInterface(o.Config); err != nil {
			return err
		}
	}

	return nil
}
