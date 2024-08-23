package wtsc

//go:generate go run github.com/pokt-scan/wtsc/cmd/schema

import (
	"encoding/json"
	"github.com/pokt-scan/wtsc/generated"
	"net/http"
)

// Application Types

type ServiceStake struct {
	// Service aka Chain ID
	Service string `json:"service"`
	// MinNode is the minimum amount of nodes (default 0 on what-to-stake)
	MinNode uint `json:"min_node"`
}

type MinServiceStake []ServiceStake

func (ss MinServiceStake) CastToGqlType() (r []generated.WtsMinServiceStakeInput) {
	for _, v := range ss {
		r = append(r, generated.WtsMinServiceStakeInput{
			Service:   v.Service,
			Min_nodes: int(v.MinNode),
		})
	}
	return
}

type Config struct {
	// DryMode allows you to run the service without impact your stake, this will just print logs and save results
	// if ResultsPath has a value.
	DryMode bool `json:"dry_mode"`
	// POKTscanApi api url
	POKTscanApi string `json:"poktscan_api"`
	// POKTscanApiToken access token
	POKTscanApiToken string `json:"poktscan_api_token"`
	// Pocket Network ID for tx (mainnet or testnet)
	// Refers to: https://github.com/pokt-foundation/pocket-go/blob/master/transaction-builder/transaction_builder.go#L36
	NetworkID string `json:"network_id"`
	// TxMemo optional tx memo
	TxMemo string `json:"tx_memo"`
	// TxFee optional tx fee
	TxFee json.Number `json:"tx_fee"`
	// Domain is Servicer domain to query what-to-stake
	Domain string `json:"domain"`
	// ChainPool list of the chains that are supported by the node runner.
	ChainPool []string `json:"chain_pool"`
	// ServicerKeys list of the private keys to use for sign nodes
	ServicerKeys []string `json:"servicer_keys"`
	// StakeWeight the stake weight that will be sent to what-to-stake service
	StakeWeight uint `json:"stake_weight"`
	// MinIncreasePercent the amount of change percent sent to what-to-stake service
	MinIncreasePercent float64 `json:"min_increase_percent"`
	// MinServiceStake (optional) just in case the need to pass a minimum amount of services on a specific chain
	MinServiceStake MinServiceStake `json:"min_service_stake"`
	// TimePeriod in hours used to get the amount of relays
	TimePeriod uint `json:"time_period"`
	// ResultPath allows you to save wts results, mostly for debug or if you think something is going wrong.
	// This will also allow you to share with POKTscan in case you think something is wrong.
	// Empty value disable this.
	ResultsPath string `json:"results_path"`
	// LogLevel is the level of logging
	LogLevel string `json:"log_level"`
	// LogFormat allows to use JSON(optimal) or ColorizedText(slower). Values allowed: json|text
	LogFormat string `json:"log_format"`
	// Schedule is the cron schedule. Allow @every 1m or cron text.
	// Refer to: https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format
	Schedule string `json:"schedule"`
	// Worker Pool config
	MaxWorkers uint `json:"max_workers"`
	// PocketRPC url to call /v1/query/nodes & /v1/client/rawTx
	PocketRPC string `json:"pocket_rpc"`
	// MaxRetries for PocketRpc and WhatToStake call
	MaxRetries uint `json:"max_retries"`
	// MaxTimeout for PocketRpc and WhatToStake call (milliseconds)
	MaxTimeout uint `json:"max_timeout"`
}

type AuthedTransport struct {
	token   string
	wrapped http.RoundTripper
}

func (t *AuthedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", t.token)
	return t.wrapped.RoundTrip(req)
}
