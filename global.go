package wtsc

import (
	"github.com/alitto/pond"
	"github.com/hashicorp/go-retryablehttp"
	pocketGoProvider "github.com/pokt-foundation/pocket-go/provider"
	pocketGoSigner "github.com/pokt-foundation/pocket-go/signer"
	pocketCoreCodec "github.com/pokt-network/pocket-core/codec"
	"github.com/puzpuzpuz/xsync"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"time"
)

var (
	HttpRetryableClient *retryablehttp.Client
	PocketRpcProvider   *pocketGoProvider.Provider
	PocketCoreCodec     *pocketCoreCodec.Codec
	WorkerPool          *pond.WorkerPool
	CronJob             *cron.Cron
	AppConfig           *Config
	ServicersMap        *xsync.MapOf[string, *pocketGoSigner.Signer]
)

// UpdateServicers adds new servicers to the servicers map and removes orphaned servicers from the map.
func UpdateServicers(servicers []string) {
	// add new
	for i, key := range servicers {
		signer, err := pocketGoSigner.NewSignerFromPrivateKey(key)
		if err != nil {
			log.Fatal().Int("index", i).Msg("unable to extract read signer")
			continue
		}

		_, found := ServicersMap.LoadOrStore(signer.GetAddress(), signer)

		if !found {
			log.Debug().Str("address", signer.GetAddress()).Msg("signer added")
		}
	}

	// remove orphans
	toRemove := make([]string, 0)
	ServicersMap.Range(func(key string, value *pocketGoSigner.Signer) bool {
		found := FindStringInSlice(servicers, value.GetPrivateKey())
		if !found {
			toRemove = append(toRemove, value.GetAddress())
		}
		return true
	})
	for _, address := range toRemove {
		ServicersMap.Delete(address)
	}
}

func NewWorker(maxWorkers, maxCapacity uint64) {
	WorkerPool = pond.New(int(maxWorkers), int(maxCapacity))
}

func NewSignerMap(servicerKeys []string) {
	ServicersMap = xsync.NewMapOf[*pocketGoSigner.Signer]()
	UpdateServicers(servicerKeys)
}

func NewHttpClient(maxRetries uint64) *retryablehttp.Client {
	HttpRetryableClient = retryablehttp.NewClient()
	HttpRetryableClient.RetryMax = int(maxRetries)
	HttpRetryableClient.Logger = NewZerologLeveledLogger()
	HttpRetryableClient.HTTPClient.Timeout = time.Duration(AppConfig.MaxTimeout) * time.Millisecond
	return HttpRetryableClient
}
