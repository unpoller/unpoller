# UnPoller - AI Agent Context

This file provides comprehensive context for AI coding assistants working on the UnPoller project.

## Project Identity

**Name**: UnPoller (UniFi Poller)  
**Language**: Go 1.25.5+  
**Purpose**: Collect metrics and events from UniFi network controllers and export to monitoring backends  
**Repository**: https://github.com/unpoller/unpoller  
**Documentation**: https://unpoller.com

## Architecture Overview

UnPoller uses a **plugin-based architecture** with a generic core that provides input/output interfaces. The core library (`pkg/poller`) has no knowledge of UniFi or specific backends - it's a generic plugin system.

### Core Components

1. **Plugin System** (`pkg/poller/`):
   - Generic input/output interfaces
   - Automatic plugin discovery via blank imports
   - Configuration management (TOML/JSON/YAML + env vars)
   - Metrics and events aggregation

2. **Input Plugins**:
   - `pkg/inputunifi/` - Collects data from UniFi controllers

3. **Output Plugins**:
   - `pkg/influxunifi/` - InfluxDB (v1 and v2)
   - `pkg/promunifi/` - Prometheus
   - `pkg/lokiunifi/` - Loki
   - `pkg/datadogunifi/` - DataDog
   - `pkg/mysqlunifi/` - MySQL

4. **Web Server** (`pkg/webserver/`):
   - Health checks
   - Metrics endpoint
   - Plugin information

## Code Organization

### Directory Structure

```
unpoller/
├── main.go                    # Entry point, loads plugins
├── pkg/
│   ├── poller/               # Core plugin system (generic)
│   │   ├── config.go         # Configuration structures
│   │   ├── inputs.go         # Input plugin interface
│   │   ├── outputs.go        # Output plugin interface
│   │   ├── start.go          # Application startup
│   │   └── commands.go       # CLI commands
│   ├── inputunifi/           # UniFi input plugin
│   ├── influxunifi/          # InfluxDB output
│   ├── promunifi/            # Prometheus output
│   ├── lokiunifi/            # Loki output
│   ├── datadogunifi/         # DataDog output
│   ├── mysqlunifi/           # MySQL output
│   └── webserver/            # Web server
├── examples/                 # Configuration examples
├── init/                     # Init scripts (systemd, docker, etc.)
└── scripts/                  # Build scripts
```

### Plugin Interface Pattern

**Input Plugin**:
```go
type Input interface {
    GetMetrics() (*Metrics, error)
    GetEvents() (*Events, error)
}
```

**Output Plugin**:
```go
type Output interface {
    WriteMetrics(*Metrics) error
    WriteEvents(*Events) error
}
```

**Registration**:
```go
func init() {
    poller.RegisterInput(&MyInput{})
    // or
    poller.RegisterOutput(&MyOutput{})
}
```

## Coding Standards

### Go Style

- **Formatting**: Standard `gofmt`
- **Linting**: `.golangci.yaml` configuration
  - Enabled: `nlreturn`, `revive`, `tagalign`, `testpackage`, `wsl_v5`
- **Imports**: Use `//nolint:gci` when import ordering needs exception

### Naming Conventions

- Packages: `lowercase`, single word
- Exported: `PascalCase`
- Unexported: `camelCase`
- Constants: `PascalCase` with descriptive names
- Errors: Always `err`, checked immediately

### Error Handling

```go
// ✅ Good
if err != nil {
    return fmt.Errorf("context: %w", err)
}

// ❌ Bad
_ = functionThatReturnsError()
```

### Configuration Pattern

All configuration structs must include:
```go
type Config struct {
    Field string `json:"field" toml:"field" xml:"field" yaml:"field"`
}
```

- Environment variables use `UP_` prefix (defined in `ENVConfigPrefix`)
- Automatic unmarshaling via `golift.io/cnfg` and `golift.io/cnfgfile`

### Testing

- Tests in `_test` packages (enforced by `testpackage` linter)
- Use `github.com/stretchr/testify` for assertions
- Integration tests: `integration_test.go` with `integration_test_expectations.yaml`

## Key Dependencies

### Core Libraries
- `github.com/unpoller/unifi/v5` - UniFi API client (local: `/Users/briangates/unifi`)
- `golift.io/cnfg` - Configuration management
- `golift.io/cnfgfile` - Config file parsing
- `github.com/spf13/pflag` - CLI flags
- `github.com/gorilla/mux` - HTTP router

### Output Backends
- **InfluxDB**: `github.com/influxdata/influxdb1-client` (v1), `github.com/influxdata/influxdb-client-go/v2` (v2)
- **Prometheus**: `github.com/prometheus/client_golang`
- **DataDog**: `github.com/DataDog/datadog-go/v5`
- **Loki**: Custom HTTP implementation

## Data Structures

### Metrics
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

### Events
```go
type Events struct {
    Logs []any
}
```

### Device Types
UniFi devices supported:
- **UAP** - Access Points (`uap.go`)
- **USG** - Security Gateways (`usg.go`)
- **USW** - Switches (`usw.go`)
- **UDM** - Dream Machines (`udm.go`)
- **UXG** - Next-Gen Gateways (`uxg.go`)
- **UBB** - Building Bridges (`ubb.go`)
- **UCI** - Industrial devices (`uci.go`)
- **PDU** - Power Distribution Units (`pdu.go`)

Each output plugin has device-specific reporting functions following the pattern:
```go
func (r *Report) UAP(uaps []*unifi.UAP) {
    // Convert and export UAP data
}
```

## Common Development Tasks

### Adding a New Output Plugin

1. Create package in `pkg/myoutput/`
2. Implement `Output` interface
3. Register in `init()` function
4. Add blank import to `main.go`: `_ "github.com/unpoller/unpoller/pkg/myoutput"`
5. Add configuration struct with tags
6. Create `README.md` with usage examples

### Adding Device Type Support

1. Add device-specific file (e.g., `mydevice.go`)
2. Follow pattern from existing device files
3. Convert UniFi device struct to output format
4. Handle errors appropriately

### Modifying Configuration

1. Add fields to appropriate config struct
2. Include all format tags: `json`, `toml`, `xml`, `yaml`
3. Update example configs in `examples/`
4. Document in package `README.md`

## Build & Deployment

- **CI/CD**: GitHub Actions (`.github/workflows/`)
- **Build**: `goreleaser` (`.goreleaser.yaml`)
- **Platforms**: Linux, macOS, Windows, FreeBSD
- **Docker**: Auto-built to `ghcr.io`
- **Packages**: Debian, RedHat via goreleaser
- **Homebrew**: Formula for macOS

## Important Constraints

1. **Generic Core**: `pkg/poller` must remain generic - no UniFi/backend knowledge
2. **Cross-Platform**: Support Windows, macOS, Linux, BSD
3. **Configuration**: Support TOML (default), JSON, YAML, and environment variables
4. **Error Handling**: Always check and return errors
5. **Context**: Use `context.Context` for cancellable operations
6. **Timeouts**: Respect timeouts and deadlines
7. **Logging**: Use structured logging via logger interfaces

## Configuration Examples

**File Formats**: See `examples/up.conf.example`, `examples/up.json.example`, `examples/up.yaml.example`

**Environment Variables**:
```bash
UP_INFLUX_URL=http://localhost:8086
UP_UNIFI_DEFAULT_USER=admin
UP_UNIFI_DEFAULT_PASS=password
UP_POLLER_DEBUG=true
UP_POLLER_INTERVAL=30s
```

## When Writing Code

1. **Follow Existing Patterns**: Look at similar plugins for examples
2. **Keep Core Generic**: Don't add UniFi/backend-specific code to `pkg/poller`
3. **Error Handling**: Check all errors, return descriptive messages
4. **Documentation**: Document exported functions and types
5. **Testing**: Write tests in `_test` packages
6. **Cross-Platform**: Test on multiple platforms when possible
7. **Performance**: Consider polling intervals and data volume
8. **Security**: Don't log passwords or sensitive data

## Resources

- **Documentation**: https://unpoller.com
- **UniFi Library**: https://github.com/unpoller/unifi
- **Grafana Dashboards**: https://grafana.com/dashboards?search=unifi-poller
- **Discord**: https://golift.io/discord (#unpoller channel)
