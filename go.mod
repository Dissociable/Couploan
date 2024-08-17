module github.com/Dissociable/Couploan

go 1.22.5

// Replace github.com/justhyped/gocaptcha with github.com/Dissociable/gocaptcha
replace github.com/justhyped/gocaptcha v1.0.4 => github.com/Dissociable/gocaptcha v1.0.5

replace github.com/bogdanfinn/fhttp v0.5.28 => github.com/Dissociable/fhttp v0.5.29

require (
	ariga.io/atlas-go-sdk v0.5.6
	entgo.io/ent v0.14.0
	github.com/Dissociable/persistent-cookiejar v0.0.0-20240309151603-d1ea45f219a9
	github.com/adrg/strutil v0.3.1
	github.com/bogdanfinn/fhttp v0.5.28
	github.com/bogdanfinn/tls-client v1.7.5
	github.com/brianvoe/gofakeit/v7 v7.0.4
	github.com/discomco/go-status v0.0.3
	github.com/eko/gocache/lib/v4 v4.1.6
	github.com/eko/gocache/store/redis/v4 v4.2.2
	github.com/elliotchance/orderedmap/v2 v2.2.0
	github.com/go-playground/validator/v10 v10.22.0
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572
	github.com/goccy/go-json v0.10.3
	github.com/gofiber/fiber/v2 v2.52.4
	github.com/gofiber/fiber/v3 v3.0.0-beta.2
	github.com/google/uuid v1.6.0
	github.com/hibiken/asynq v0.24.1
	github.com/hibiken/asynq/x v0.0.0-20240506061152-d04888e74845
	github.com/jackc/pgx/v5 v5.6.0
	github.com/justhyped/gocaptcha v1.0.4
	github.com/klauspost/compress v1.17.9
	github.com/labstack/echo/v4 v4.12.0
	github.com/microcosm-cc/bluemonday v1.0.26
	github.com/mwitkow/go-http-dialer v0.0.0-20161116154839-378f744fb2b8
	github.com/phuslu/shardmap v0.0.0-20230929024548-c0f3d8a4fccd
	github.com/pkg/errors v0.9.1
	github.com/redis/go-redis/v9 v9.5.3
	github.com/sourcegraph/conc v0.3.0
	github.com/spf13/viper v1.19.0
	github.com/stretchr/testify v1.9.0
	github.com/teris-io/shortid v0.0.0-20220617161101-71ec9f2aa569
	github.com/tidwall/gjson v1.17.1
	github.com/txthinking/socks5 v0.0.0-20230325130024-4230056ae301
	github.com/wk8/go-ordered-map/v2 v2.1.8
	go.uber.org/zap v1.21.0
	golang.org/x/net v0.27.0
	golang.org/x/text v0.16.0
	golang.org/x/time v0.5.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	ariga.io/atlas v0.25.0 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bogdanfinn/utls v1.6.1 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudflare/circl v1.3.6 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/go-openapi/inflect v0.21.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/gofiber/utils/v2 v2.0.0-beta.4 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/hcl/v2 v2.21.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/juju/go4 v0.0.0-20160222163258-40d72ab9641a // indirect
	github.com/justlovediaodiao/udp-over-tcp v0.0.0-20230616061358-151b7c7a401d // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/quic-go/quic-go v0.37.4 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/cobra v1.7.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/tam7t/hpkp v0.0.0-20160821193359-2b70b4024ed5 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tinylib/msgp v1.1.8 // indirect
	github.com/txthinking/runnergroup v0.0.0-20210608031112-152c7c4432bf // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.55.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/zclconf/go-cty v1.14.4 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/exp v0.0.0-20240416160154-fe59bbe5cc7f // indirect
	golang.org/x/mod v0.19.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/tools v0.23.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/errgo.v1 v1.0.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/retry.v1 v1.0.3 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
