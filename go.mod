module github.com/pokt-scan/wtsc

go 1.21

require (
	github.com/Khan/genqlient v0.7.0
	github.com/alitto/pond v1.9.1
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-retryablehttp v0.7.7
	github.com/pokt-foundation/pocket-go v0.21.0
	github.com/pokt-network/pocket-core v0.0.0-20240814175146-7f936ff73532
	github.com/puzpuzpuz/xsync v1.5.2
	github.com/robfig/cron/v3 v3.0.1
	github.com/rs/zerolog v1.33.0
	github.com/suessflorian/gqlfetch v0.6.0
	github.com/tendermint/tendermint v0.33.7
)

replace github.com/tendermint/tendermint => github.com/pokt-network/tendermint v0.32.11-0.20230426215212-59310158d3e9

replace github.com/tendermint/tm-db => github.com/pokt-network/tm-db v0.5.2-0.20220118210553-9b2300f289ba

// Added this replace because current pocket-go does not have rewards delegators on his rpc output.
replace github.com/pokt-foundation/pocket-go => github.com/pokt-scan/pocket-go v0.0.0-20240824192825-4342359fb21c

require (
	github.com/ChainSafe/go-schnorrkel v0.0.0-20200405005733-88cbf1b4c40d // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/alexflint/go-arg v1.4.2 // indirect
	github.com/alexflint/go-scalar v1.0.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/btcsuite/btcd v0.20.1-beta // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d // indirect
	github.com/cosmos/gogoproto v1.4.10 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-kit/kit v0.12.0 // indirect
	github.com/go-kit/log v0.2.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/gtank/merlin v0.1.1 // indirect
	github.com/gtank/ristretto255 v0.1.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/jmhodges/levigo v1.0.0 // indirect
	github.com/jordanorelli/lexnum v0.0.0-20141216151731-460eeb125754 // indirect
	github.com/libp2p/go-buffer-pool v0.0.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mimoo/StrobeGo v0.0.0-20181016162300-f8f6d4d2b643 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/pokt-foundation/utils-go v0.7.0 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.30.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/regen-network/cosmos-proto v0.3.0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7 // indirect
	github.com/tecbot/gorocksdb v0.0.0-20191217155057-f0fad39f321c // indirect
	github.com/tendermint/go-amino v0.15.1 // indirect
	github.com/tendermint/tm-db v0.5.1 // indirect
	github.com/vektah/gqlparser v1.3.1 // indirect
	github.com/vektah/gqlparser/v2 v2.5.11 // indirect
	github.com/willf/bitset v1.1.10 // indirect
	github.com/willf/bloom v2.0.3+incompatible // indirect
	go.etcd.io/bbolt v1.3.3 // indirect
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/exp v0.0.0-20230131160201-f062dba9d201 // indirect
	golang.org/x/mod v0.15.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.18.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
	google.golang.org/grpc v1.55.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
