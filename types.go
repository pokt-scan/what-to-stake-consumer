package wtsc

import "encoding/json"

// Application Types

type ServiceStake struct {
	// Service aka Chain ID
	Service string `json:"service"`
	// MinNode is the minimum amount of nodes (default 0 on what-to-stake)
	MinNode uint `json:"min_node"`
}

type Config struct {
	// WhatToStakeService service url
	WhatToStakeService string `json:"what_to_stake_service"`
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
	MinIncreasePercent float32 `json:"min_increase_percent"`
	// MinServiceStake (optional) just in case the need to pass a minimum amount of services on a specific chain
	MinServiceStake []ServiceStake `json:"min_service_stake"`
	// LogLevel is the level of logging
	LogLevel string `json:"log_level"`
	// Schedule is the cron schedule. Allow @every 1m or cron text.
	// Refer to: https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format
	Schedule string `json:"schedule"`
	// Worker Pool config
	MaxWorkers uint64 `json:"max_workers"`
	// PocketRPC url to call /v1/query/nodes & /v1/client/rawTx
	PocketRPC string `json:"pocket_rpc"`
	// MaxRetries for PocketRpc and WhatToStake call
	MaxRetries uint64 `json:"max_retries"`
	// MaxTimeout for PocketRpc and WhatToStake call (milliseconds)
	MaxTimeout uint64 `json:"max_timeout"`
}

// What To Stake Types

type WTSParams struct {
	Domain             string         `json:"domain"`
	ServicePool        []string       `json:"service_pool"`
	MinIncreasePercent float32        `json:"min_increase_percent"`
	StakeWeight        uint           `json:"stake_weight"`
	MinServiceStake    []ServiceStake `json:"min_service_stake"`
}

type WTSService struct {
	Address string   `json:"address"`
	Chains  []string `json:"chains"`
}

type WTSResponse struct {
	DoUpdate  bool          `json:"do_update"`
	Reason    string        `json:"reason"`
	Servicers []*WTSService `json:"servicers"`
}
