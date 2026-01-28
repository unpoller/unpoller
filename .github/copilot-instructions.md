# UnPoller - GitHub Copilot Instructions

## Project Overview

UnPoller is a Go application that collects metrics and events from UniFi network controllers and exports them to monitoring backends (InfluxDB, Prometheus, Loki, DataDog, MySQL). The application uses a plugin-based architecture with a generic core.

## Architecture

### Plugin System
- **Core Library** (`pkg/poller/`): Generic plugin system with input/output interfaces
- **Input Plugins**: Collect data (e.g., `pkg/inputunifi/`)
- **Output Plugins**: Export data to backends (e.g., `pkg/influxunifi/`, `pkg/promunifi/`)
- **Plugin Discovery**: Automatic via blank imports in `main.go`

### Key Interfaces

**Input Interface**:
```go
type Input interface {
    GetMetrics() (*Metrics, error)
    GetEvents() (*Events, error)
}
```

**Output Interface**:
```go
type Output interface {
    WriteMetrics(*Metrics) error
    WriteEvents(*Events) error
}
```

## Code Style

### Go Conventions
- Go 1.25.5+
- Use `gofmt` formatting
- Follow `.golangci.yaml` linting rules
- Enabled linters: `nlreturn`, `revive`, `tagalign`, `testpackage`, `wsl_v5`

### Naming
- Packages: `lowercase`
- Exported: `PascalCase`
- Unexported: `camelCase`
- Errors: Always `err`, check immediately

### Error Handling
```go
if err != nil {
    return fmt.Errorf("context: %w", err)
}
```

### Configuration
All config structs must include format tags:
```go
type Config struct {
    Field string `json:"field" toml:"field" xml:"field" yaml:"field"`
}
```

Environment variables use `UP_` prefix.

## Common Patterns

### Plugin Registration
```go
func init() {
    poller.RegisterInput(&Input{})
    // or
    poller.RegisterOutput(&Output{})
}
```

### Device Type Reporting
Each output plugin has device-specific functions:
- `UAP()` - Access Points
- `USG()` - Security Gateways
- `USW()` - Switches
- `UDM()` - Dream Machines
- `UXG()` - Next-Gen Gateways
- `UBB()` - Building Bridges
- `UCI()` - Industrial devices
- `PDU()` - Power Distribution Units

### Adding New Plugin
1. Create package in `pkg/`
2. Implement `Input` or `Output` interface
3. Register in `init()`
4. Add blank import to `main.go`
5. Add config struct with tags
6. Create `README.md`

## Dependencies

### Core
- `github.com/unpoller/unifi/v5` - UniFi API (local: `/Users/briangates/unifi`)
- `golift.io/cnfg` - Configuration
- `golift.io/cnfgfile` - Config parsing
- `github.com/spf13/pflag` - CLI flags

### Outputs
- InfluxDB: `github.com/influxdata/influxdb1-client`, `github.com/influxdata/influxdb-client-go/v2`
- Prometheus: `github.com/prometheus/client_golang`
- DataDog: `github.com/DataDog/datadog-go/v5`
- Loki: Custom HTTP client

## Important Notes

- `pkg/poller` core is generic - no UniFi/backend knowledge
- Support Windows, macOS, Linux, BSD
- Configuration: TOML (default), JSON, YAML, env vars
- Always check and return errors
- Use `context.Context` for cancellable operations
- Tests in `_test` packages
- Use `github.com/stretchr/testify` for assertions

## When Writing Code

1. Follow existing plugin patterns
2. Keep core library generic
3. Check all errors
4. Document exported functions
5. Write tests in `_test` packages
6. Consider cross-platform compatibility
7. Use structured logging
8. Respect timeouts and deadlines
