module github.com/unpoller/unpoller

go 1.24

toolchain go1.24.2

require (
	github.com/DataDog/datadog-go/v5 v5.6.0
	github.com/gorilla/mux v1.8.1
	github.com/influxdata/influxdb1-client v0.0.0-20220302092344-a9ab5670611c
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.22.0
	github.com/prometheus/common v0.63.0
	github.com/spf13/pflag v1.0.6
	github.com/stretchr/testify v1.10.0
	github.com/unpoller/unifi/v5 v5.1.0
	golang.org/x/crypto v0.38.0
	golang.org/x/net v0.40.0
	golang.org/x/term v0.32.0
	golift.io/cnfg v0.2.3
	golift.io/cnfgfile v0.0.0-20240713024420-a5436d84eb48
	golift.io/version v0.0.2
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/brianvoe/gofakeit/v6 v6.28.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/influxdata/influxdb-client-go/v2 v2.14.0
	github.com/influxdata/line-protocol v0.0.0-20210922203350-b1ad95c89adf // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/oapi-codegen/runtime v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/procfs v0.16.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

// for local iterative development only
// replace github.com/unpoller/unifi/v5 => ../unifi
