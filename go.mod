module github.com/traefik/grpc-web

go 1.19

require (
	github.com/cenkalti/backoff/v4 v4.2.1
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f
	github.com/golang/protobuf v1.5.3
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f
	github.com/mwitkow/grpc-proxy v0.0.0-20181017164139-0f1106ef9c76
	github.com/prometheus/client_golang v1.12.1
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.4
	golang.org/x/net v0.12.0
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.27.1
	nhooyr.io/websocket v1.8.7
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/genproto v0.0.0-20210126160654-44e461bb6506 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Replace to handle https://coder.com/blog/websocket
replace nhooyr.io/websocket => github.com/coder/websocket v1.8.7
