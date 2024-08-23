# What To Stake?

This tool provides optimal* deployments of a series of nodes in the POKT Network.

### Optimization Process

Before using the provided recommendations, it's important to understand several key points regarding the optimization process:

- The optimization is based on a basic model of the network, as some data (e.g., gateway selection mechanisms) isn't available.
- Nodes' QoS (Quality of Service) is not taken into account; all nodes are treated as "equals". Good nodes will yield higher rewards than predicted, while poor nodes will yield lower rewards.
- We optimize for "round-robin" relay distribution.
- We don't know where your nodes or the traffic are located, which significantly impacts results. We cannot advise on the deployment locations of your chain nodes.
- The optimization is **domain-based**; it provides optimal deployment for the **current** network, not a balanced network.
- Given the dynamic nature of the network and other participants' reactions to your changes, be prepared to adjust your stakes periodically.
- We expect network state convergence if enough node runners use this tool, but **convergence is not guaranteed**.
- Recommendations may oscillate (i.e., a proposed strategy leads to a previous stake strategy). We don't expect this, but it could happen due to adversarial strategies among node runners.
- By default, the tool considers relay averages over the past 24 hours. The chosen period for calculating average relays greatly affects outcomes. For example, if a chain had no relays in the last 24 hours, it won't be considered during optimization (no relays = no rewards).

**The optimization results are optimal given the aforementioned simplifications.**

### Requirements

1. POKTscan API Token. Get yours at: [POKTscan API Token](https://poktscan.com/teams?tab=api_token)

### How To Use

1. Create a copy of `config.json.sample` and name it `config.json`.
2. Edit `config.json` with your desired values.
3. Generate the GraphQL types by running: `go generate`.
4. Start the service with Docker: `docker compose up --build`.

### Understanding `config.json` Values

```json
{
  "dry_mode": true, // If true, calls to "What to Stake" will write results to `results_path` or print them if `results_path` is empty.
  "poktscan_api": "https://api.poktscan.com/poktscan/api/graphql", // POKTscan API endpoint
  "poktscan_api_token": "", // POKTscan API token
  "network_id": "testnet", // Can be "mainnet" or "testnet"
  "tx_memo": "wtsc", // Transaction memo
  "tx_fee": 10000, // Transaction fee
  "domain": "example.com", // Your nodes domain (e.g., "poktscan.cloud" or "c0d3r.org")
  "chain_pool": ["0021", "0003"], // Chain IDs available in your fleet
  "servicer_keys": ["PRIVATE_KEY"], // List of private keys for signing stake transactions
  "stake_weight": 4, // Used by the "What to Stake" service to estimate potential rewards (1-4)
  "min_increase_percent": 5, // Minimum percentage increase expected to process stakes
  "min_service_stake": [ // Minimum number of nodes for specific services
    {
      "service": "0021",
      "min_node": 1
    }
  ],
  "time_period": 24, // Time in hours to consider for relay averages
  "results_path": "", // Path to save "What to Stake" results (empty to disable)
  "pocket_rpc": "http://localhost:8081", // Pocket node or load balancer URL
  "log_level": "debug", // Log level
  "log_format": "json", // Log format (json or colorized text)
  "schedule": "@every 5m", // Frequency of the "What to Stake" service calls
  "max_workers": 1, // Number of workers to process stake transactions in parallel
  "max_retries": 1, // Number of retries for HTTP calls (POKTscan API or Pocket RPC)
  "max_timeout": 15000 // Timeout in milliseconds for HTTP calls
}
```

### FAQ

#### Do I need to restart the WTSC if I change `config.json`?

No, the service supports hot reload, configurable via the `RELOAD_SECONDS` environment variable.