module go.opentelemetry.io/collector

go 1.18

require (
	contrib.go.opencensus.io/exporter/prometheus v0.4.2
	github.com/cenkalti/backoff/v4 v4.1.3
	github.com/f5/otel-arrow-adapter v0.0.0-20221111033209-0140de5d0a01
	github.com/gogo/protobuf v1.3.2
	github.com/golang/mock v1.6.0
	github.com/golang/snappy v0.0.4
	github.com/google/go-cmp v0.5.9
	github.com/google/uuid v1.3.0
	github.com/klauspost/compress v1.15.11
	github.com/knadh/koanf v1.4.4
	github.com/magiconair/properties v1.8.6
	github.com/mitchellh/mapstructure v1.5.0
	github.com/mostynb/go-grpc-compression v1.1.17
	github.com/prometheus/client_golang v1.13.0
	github.com/prometheus/client_model v0.3.0
	github.com/prometheus/common v0.37.0
	github.com/rs/cors v1.8.2
	github.com/shirou/gopsutil/v3 v3.22.9
	github.com/spf13/cobra v1.6.1
	github.com/stretchr/testify v1.8.1
	go.opencensus.io v0.23.0
	go.opentelemetry.io/collector/pdata v0.63.1
	go.opentelemetry.io/collector/semconv v0.63.1
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.36.4
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.36.4
	go.opentelemetry.io/contrib/propagators/b3 v1.11.1
	go.opentelemetry.io/contrib/zpages v0.36.4
	go.opentelemetry.io/otel v1.11.1
	go.opentelemetry.io/otel/exporters/prometheus v0.33.0
	go.opentelemetry.io/otel/metric v0.33.0
	go.opentelemetry.io/otel/sdk v1.11.1
	go.opentelemetry.io/otel/sdk/metric v0.33.0
	go.opentelemetry.io/otel/trace v1.11.1
	go.uber.org/atomic v1.10.0
	go.uber.org/multierr v1.8.0
	go.uber.org/zap v1.23.0
	golang.org/x/net v0.1.0
	golang.org/x/sys v0.1.0
	google.golang.org/genproto v0.0.0-20220126215142-9970aeb2e350
	google.golang.org/grpc v1.50.1
	google.golang.org/protobuf v1.28.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/apache/arrow/go/v11 v11.0.0-20221108144011-89a65b06a199 // indirect
	github.com/apache/thrift v0.16.0 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/goccy/go-json v0.9.11 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/flatbuffers v2.0.8+incompatible // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/asmfmt v1.3.2 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/minio/asm2plan9s v0.0.0-20200509001527-cdd76441f9d8 // indirect
	github.com/minio/c2goasm v0.0.0-20190812172519-36a3d3bbc4f3 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/prometheus/statsd_exporter v0.22.7 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/tklauser/numcpus v0.4.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	golang.org/x/mod v0.6.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	golang.org/x/tools v0.2.0 // indirect
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace go.opentelemetry.io/collector/semconv => ./semconv

replace go.opentelemetry.io/collector/pdata => ./pdata

retract (
	v0.57.1 // Release failed, use v0.57.2
	v0.57.0 // Release failed, use v0.57.2
	v0.32.0 // Contains incomplete metrics transition to proto 0.9.0, random components are not working.
)
