# UnPoller - Claude Code Context

## Project Overview

UnPoller is a Go application that collects metrics and events from UniFi network controllers and exports them to various time-series databases and monitoring systems. The application uses a plugin-based architecture where input plugins collect data and output plugins export it to different backends.

## Architecture

### Core Design Principles

1. **Plugin-Based Architecture**: The core `pkg/poller` library is generic and provides interfaces for input and output plugins. It has no knowledge of UniFi, InfluxDB, Prometheus, etc.

2. **Automatic Plugin Discovery**: Plugins are loaded via blank imports in `main.go`. The poller automatically discovers and initializes all imported plugins.

3. **Flexible Configuration**: Supports TOML (default), JSON, and YAML configuration files, plus environment variables with `UP_` prefix.

4. **Multiple Backends**: Supports InfluxDB, Prometheus, Loki, DataDog, and MySQL as output backends.

### Package Structure

```
pkg/
├── poller/          # Core plugin system - generic, no UniFi/backend knowledge
├── inputunifi/      # UniFi controller input plugin
├── influxunifi/     # InfluxDB output plugin
├── promunifi/       # Prometheus output plugin
├── lokiunifi/       # Loki output plugin
├── datadogunifi/    # DataDog output plugin
├── mysqlunifi/      # MySQL output plugin
└── webserver/       # Web server for health checks and metrics
```

### Key Interfaces

**Input Interface** (`pkg/poller/inputs.go`):
- `GetMetrics()` - Returns aggregated metrics
- `GetEvents()` - Returns events/logs
- Plugins must implement this to provide data

**Output Interface** (`pkg/poller/outputs.go`):
- `WriteMetrics()` - Writes metrics to backend
- `WriteEvents()` - Writes events to backend
- Plugins must implement this to export data

## Code Style Guidelines

### Go Conventions

- **Go Version**: 1.25.5+
- **Formatting**: Use `gofmt` standard formatting
- **Linting**: Follow `.golangci.yaml` configuration
  - Enabled: `nlreturn`, `revive`, `tagalign`, `testpackage`, `wsl_v5`
  - Use `//nolint:gci` for import ordering exceptions

### Naming

- **Packages**: Lowercase, single word
- **Exported**: PascalCase (types, functions, constants)
- **Unexported**: camelCase
- **Errors**: Always named `err`, checked immediately

### Error Handling

```go
// Good: Check and return errors
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// Bad: Ignoring errors
_ = someFunction()
```

### Configuration Pattern

```go
type Config struct {
    URL      string `json:"url" toml:"url" yaml:"url"`
    Interval string `json:"interval" toml:"interval" yaml:"interval"`
    Timeout  string `json:"timeout" toml:"timeout" yaml:"timeout"`
}
```

- Always include `json`, `toml`, `xml`, `yaml` tags
- Environment variables use `UP_` prefix automatically
- Use `golift.io/cnfg` for configuration management

### Testing

- Tests in separate `_test` packages (enforced by `testpackage` linter)
- Use `github.com/stretchr/testify` for assertions
- Integration tests use `integration_test_expectations.yaml`

## Common Tasks

### Adding a New Output Plugin

1. **Create Package Structure**:
   ```go
   package myoutput
   
   import (
       "github.com/unpoller/unpoller/pkg/poller"
   )
   
   type Output struct {
       // Configuration fields
   }
   
   func (o *Output) WriteMetrics(m *poller.Metrics) error {
       // Implementation
   }
   
   func (o *Output) WriteEvents(e *poller.Events) error {
       // Implementation
   }
   
   func init() {
       poller.RegisterOutput(&Output{})
   }
   ```

2. **Add Blank Import** to `main.go`:
   ```go
   _ "github.com/unpoller/unpoller/pkg/myoutput"
   ```

3. **Add Configuration** to config structs with proper tags

4. **Create README.md** documenting usage

### Adding Device Type Support

Each output plugin has device-specific files:
- `uap.go` - Access Points
- `usg.go` - Security Gateways
- `usw.go` - Switches
- `udm.go` - Dream Machines
- `uxg.go` - Next-Gen Gateways
- `ubb.go` - Building Bridges
- `uci.go` - Industrial devices
- `pdu.go` - Power Distribution Units

Follow the pattern in existing files:
```go
func (r *Report) UAP(uaps []*unifi.UAP) {
    // Convert UniFi device data to output format
}
```

### Working with UniFi Data

- UniFi library: `github.com/unpoller/unifi/v5` (local replace at `/Users/briangates/unifi`)
- Device types come from the UniFi controller API
- Data includes: sites, clients, devices, DPI data, speed tests, country traffic
- Events include: system logs, alarms, IDS events, anomalies

## Dependencies

### Core
- `github.com/unpoller/unifi/v5` - UniFi API client
- `golift.io/cnfg` - Configuration management
- `golift.io/cnfgfile` - Config file parsing
- `github.com/spf13/pflag` - CLI flags

### Output Backends
- InfluxDB: `github.com/influxdata/influxdb1-client` (v1) and `github.com/influxdata/influxdb-client-go/v2` (v2)
- Prometheus: `github.com/prometheus/client_golang`
- DataDog: `github.com/DataDog/datadog-go/v5`
- Loki: Custom HTTP client implementation

## Important Patterns

### Plugin Registration
```go
func init() {
    poller.RegisterInput(&Input{})
    // or
    poller.RegisterOutput(&Output{})
}
```

### Configuration Loading
Configuration is automatically loaded from:
1. Config file (TOML/JSON/YAML) - path specified via `--config` flag or defaults
2. Environment variables with `UP_` prefix
3. CLI flags

### Metrics Structure
```go
type Metrics struct {
    TS             time.Time
    Sites          []any
    Clients        []any
    SitesDPI       []any
    ClientsDPI     []any
    Devices        []any
    RogueAPs       []any
    SpeedTests     []any
    CountryTraffic []any
}
```

### Events Structure
```go
type Events struct {
    Logs []any
}
```

## Build & Deployment

- **CI/CD**: GitHub Actions
- **Build Tool**: `goreleaser` for multi-platform builds
- **Platforms**: Linux, macOS, Windows, FreeBSD
- **Docker**: Images built automatically to `ghcr.io`
- **Homebrew**: Formula for macOS users
- **Packages**: Debian, RedHat packages via goreleaser

## When Writing Code

1. **Keep it Generic**: The `pkg/poller` core should remain generic
2. **Follow Patterns**: Look at existing plugins for examples
3. **Error Handling**: Always check and return errors
4. **Cross-Platform**: Consider Windows, macOS, Linux, BSD differences
5. **Context Usage**: Use `context.Context` for cancellable operations
6. **Timeouts**: Respect timeouts and deadlines
7. **Documentation**: Document exported functions and types
8. **Testing**: Write tests in separate `_test` packages

## Configuration Examples

Configuration files support multiple formats. See `examples/` directory:
- `up.conf.example` - TOML format
- `up.json.example` - JSON format
- `up.yaml.example` - YAML format

Environment variables use `UP_` prefix:
- `UP_INFLUX_URL=http://localhost:8086`
- `UP_UNIFI_DEFAULT_USER=admin`
- `UP_POLLER_DEBUG=true`
