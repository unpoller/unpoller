module github.com/unpoller/unpoller

go 1.20

require (
	github.com/DataDog/datadog-go v4.8.3+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/influxdata/influxdb1-client v0.0.0-20200827194710-b269163b24ab
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.17.0
	github.com/prometheus/common v0.45.0
	github.com/spf13/pflag v1.0.6-0.20201009195203-85dd5c8bc61c
	github.com/stretchr/testify v1.8.4
	golang.org/x/crypto v0.14.0
	golang.org/x/net v0.17.0
	golang.org/x/term v0.13.0
	golift.io/cnfg v0.2.2
	golift.io/cnfgfile v0.0.0-20220509075834-08755d9ef3f5
	golift.io/version v0.0.2
)

require (
	github.com/Microsoft/go-winio v0.6.0 // indirect
	github.com/brianvoe/gofakeit/v6 v6.23.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deepmap/oapi-codegen v1.8.2 // indirect
	github.com/influxdata/line-protocol v0.0.0-20200327222509-2487e7298839 // indirect
	github.com/matttproud/golang_protobuf_extensions/v2 v2.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/tools v0.1.12 // indirect
)

require (
	github.com/BurntSushi/toml v1.0.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/influxdata/influxdb-client-go/v2 v2.12.3
	github.com/prometheus/client_model v0.4.1-0.20230718164431-9a2bf3000d16 // indirect
	github.com/prometheus/procfs v0.11.1 // indirect
	github.com/unpoller/unifi v0.3.15
	golang.org/x/sys v0.13.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1
)

// for local iterative development only
// replace github.com/unpoller/unifi => ../unifi
