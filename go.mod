module github.com/unpoller/unpoller

go 1.21

toolchain go1.22.1

require (
	github.com/DataDog/datadog-go v4.8.3+incompatible
	github.com/gorilla/mux v1.8.1
	github.com/influxdata/influxdb1-client v0.0.0-20200827194710-b269163b24ab
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.19.0
	github.com/prometheus/common v0.52.2
	github.com/spf13/pflag v1.0.6-0.20201009195203-85dd5c8bc61c
	github.com/stretchr/testify v1.9.0
	golang.org/x/crypto v0.22.0
	golang.org/x/net v0.24.0
	golang.org/x/term v0.19.0
	golift.io/cnfg v0.2.3
	golift.io/cnfgfile v0.0.0-20230531075023-f880041cc0a0
	golift.io/version v0.0.2
)

require (
	github.com/Microsoft/go-winio v0.6.0 // indirect
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/brianvoe/gofakeit/v6 v6.28.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/influxdata/line-protocol v0.0.0-20200327222509-2487e7298839 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/oapi-codegen/runtime v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/unpoller/unifi v0.4.2
	golang.org/x/mod v0.7.0 // indirect
	golang.org/x/tools v0.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/influxdata/influxdb-client-go/v2 v2.13.0
	github.com/prometheus/client_model v0.6.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

// for local iterative development only
// replace github.com/unpoller/unifi => ../unifi
