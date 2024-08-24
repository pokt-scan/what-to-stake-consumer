package wtsc

import (
	"context"
	"errors"
	"github.com/Khan/genqlient/graphql"
	"github.com/alitto/pond"
	"github.com/hashicorp/go-retryablehttp"
	pocketGoProvider "github.com/pokt-foundation/pocket-go/provider"
	pocketGoSigner "github.com/pokt-foundation/pocket-go/signer"
	pocketCoreCodec "github.com/pokt-network/pocket-core/codec"
	"github.com/puzpuzpuz/xsync"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"net/http"
	"runtime/debug"
	"time"
)

var (
	Logger            zerolog.Logger
	HttpClient        *retryablehttp.Client
	POKTscanApiClient *graphql.Client
	PocketRpcProvider *pocketGoProvider.Provider
	PocketCoreCodec   *pocketCoreCodec.Codec
	WorkerPool        *pond.WorkerPool
	CronJob           *cron.Cron
	AppConfig         *Config
	ServicersMap      *xsync.MapOf[string, *pocketGoSigner.Signer]
)

// UpdateServicers adds new servicers to the servicers map and removes orphaned servicers from the map.
func UpdateServicers(servicers []string) (int, error) {
	Logger.Info().Msg("updating signer map")
	currentSignerMapSize := ServicersMap.Size()

	// add new
	signers := make([]*pocketGoSigner.Signer, 0)
	for i, key := range servicers {
		signer, e := pocketGoSigner.NewSignerFromPrivateKey(key)
		if e != nil {
			return i, e
		}
		signers = append(signers, signer)
	}

	// do a second round of iter here to prevent update the map if we get error on any signer.
	added := 0
	for _, signer := range signers {
		_, found := ServicersMap.LoadOrStore(signer.GetAddress(), signer)

		if !found {
			added++
			Logger.Debug().Str("address", signer.GetAddress()).Msg("signer added")
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

	if currentSignerMapSize > 0 {
		Logger.Debug().
			Int("added", added).
			Int("removed", len(toRemove)).
			Msg("signer map updated")
	} else {
		Logger.Debug().Int("added", added).Msg("signer map loaded")
	}

	return -1, nil
}

func NewWorker(maxWorkers, maxCapacity uint) {
	Logger.Info().Msg("preparing worker pool")
	WorkerPool = pond.New(
		int(maxWorkers),                        // max amount of parallel workers (configurable by config.json)
		int(maxCapacity),                       // max amount of tasks in queue before block, it will be the amount of servicers
		pond.IdleTimeout(100*time.Millisecond), // remove unused workers after 100ms
		pond.PanicHandler(func(panic interface{}) {
			// enhance log and use logger instead of fmt (default)
			Logger.Error().
				Str("stack", string(debug.Stack())).
				Interface("panic", panic).
				Msg("Worker exits from a panic")
		}),
	)
}

func NewSignerMap(servicerKeys []string) {
	Logger.Info().Msg("preparing signer map")
	ServicersMap = xsync.NewMapOf[*pocketGoSigner.Signer]()
	if index, err := UpdateServicers(servicerKeys); err != nil {
		// stop everything
		Logger.Fatal().Err(err).Int("index", index).Msg("unable to update servicers")
	}
}

func NewHttpClient(token string, maxRetries, maxTimeout uint) {
	Logger.Info().Msg("preparing http client")
	HttpClient = retryablehttp.NewClient()
	HttpClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if err != nil || resp.StatusCode == http.StatusOK {
			return false, err
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			// Retry-After: Remaining seconds to get the limit of the burst rate limiter reset.
			// X-RateLimit-Limit: The total amount of credits available in the API token for the burst rate limiter.
			// X-RateLimit-Remaining: The amount of spendable credits remaining in the API token for the current minute.
			// X-RateLimit-Reset: The approximate time and date when the amount of spendable credits is renewed for the API token for the burst rate limiter.
			// X-Long-RateLimit-Limit: The total amount of credits available in the API token for the long rate limiter.
			// X-Long-RateLimit-Remaining: The amount of spendable credits remaining in the API token for the current month.
			// X-Long-RateLimit-Consumed-Points: The amount of credits consumed in the request for the long rate limiter.
			Logger.Warn().
				Str("Retry-After", resp.Header.Get("Retry-After")).
				Str("X-RateLimit-Limit", resp.Header.Get("X-RateLimit-Limit")).
				Str("X-RateLimit-Remaining", resp.Header.Get("X-RateLimit-Remaining")).
				Str("X-RateLimit-Reset", resp.Header.Get(" X-RateLimit-Reset")).
				Str("X-Long-RateLimit-Limit", resp.Header.Get("X-Long-RateLimit-Limit")).
				Str("X-Long-RateLimit-Remaining", resp.Header.Get("X-Long-RateLimit-Remaining")).
				Str("X-Long-RateLimit-Consumed-Points", resp.Header.Get("X-Long-RateLimit-Consumed-Points")).
				Msg("We are hitting POKTscan API Rate Limit. You may need to consider lower the rate of the queries to WTS.")
		}

		// if rate-limit is hit we will log it too, but do not repeat it, or if the payload is wrong for x or y reason
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusBadRequest {
			return false, errors.New(resp.Status)
		}

		return true, nil
	}
	HttpClient.RetryMax = int(maxRetries)
	HttpClient.Logger = NewZerologLeveledLogger(Logger)
	HttpClient.HTTPClient.Timeout = time.Duration(maxTimeout) * time.Millisecond
	HttpClient.HTTPClient.Transport = &AuthedTransport{
		token: token,
		// wrap it to add authorization header on each request to poktscan api
		wrapped: HttpClient.HTTPClient.Transport,
	}
}

func NewPOKTscanClient(url string) {
	Logger.Info().Msg("preparing poktscan api client")
	gClient := graphql.NewClient(url, HttpClient.StandardClient())
	POKTscanApiClient = &gClient
}
