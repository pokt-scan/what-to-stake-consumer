# What To Stake?

This tool provides optimal* deployments of a series of nodes in the POKT Network.

### Optimization Process

Before using the provided recommendations, it's important to understand several key points regarding the optimization process:

- **Network Model Limitations:** The optimization is based on a basic model of the network, as some data (e.g., gateway selection mechanisms) isn't available.
- **Quality of Service:** Nodes' QoS (Quality of Service) is not taken into account; all nodes are treated as "equals". Good nodes will yield higher rewards than predicted, while poor nodes will yield lower rewards.
- **Relay Distribution:** We optimize for "round-robin" relay distribution.
- **Deployment Locations:** We don't know where your nodes or the traffic are located. This significantly impacts results, and we cannot advise on the deployment locations of your chain nodes.
- **Current Network Focus:** The optimization is **domain-based** and provides optimal deployment for the **current** network, not a balanced network.
- **Dynamic Adjustments:** Given the dynamic nature of the network and other participants' reactions to your changes, be prepared to adjust your stakes periodically.
- **Network Convergence:** We expect network state convergence if enough node runners use this tool. However, **convergence is not guaranteed**.
- **Oscillating Recommendations:** Recommendations may oscillate, leading to a previous stake strategy due to adversarial strategies among node runners.
- **Relay Averages:** By default, the tool considers relay averages over the past 24 hours. The chosen period for calculating average relays greatly affects outcomes. For example, if a chain had no relays in the last 24 hours, it won't be considered during optimization (no relays = no rewards).

**The optimization results are optimal given the aforementioned simplifications.**

### Requirements

1. POKTscan API Token. Get yours at: [POKTscan API Token](https://poktscan.com/teams?tab=api_token)

### How To Use

1. Create a copy of `config.json.sample` and name it `config.json`.
2. Edit `config.json` with your desired values.
3. Generate the GraphQL types by running: `go generate`.
4. Start the service with Docker:
```sh
docker compose up --build
```
An existing Docker image build is available at [Docker Hub](https://hub.docker.com/repository/docker/poktscan/wtsc/general).

If you want/need to modify the path and name of the config file, please use `CONFIG_FILE` to override the default `./config.json`

### Understanding `config.json` Values

| Parameter             | Type                | Description                                                                                                      |
|-----------------------|---------------------|------------------------------------------------------------------------------------------------------------------|
| dry_mode              | boolean             | If true, calls to "What to Stake" will write results to `results_path` or print them if `results_path` is empty. |
| poktscan_api          | string              | POKTscan API endpoint                                                                                            |
| poktscan_api_token    | string              | POKTscan API token                                                                                               |
| network_id            | string              | Can be "mainnet" or "testnet"                                                                                    |
| tx_memo               | string              | Transaction memo                                                                                                 |
| tx_fee                | integer             | Transaction fee                                                                                                  |
| domain                | string              | Your node's domain (e.g., "poktscan.cloud" or "c0d3r.org")                                                       |
| chain_pool            | array of strings    | Chain IDs available in your fleet                                                                                |
| servicer_keys         | array of strings    | List of private keys for signing stake transactions                                                              |
| stake_weight          | integer             | Used by the "What to Stake" service to estimate potential rewards (1-4)                                          |
| min_increase_percent  | integer             | Minimum percentage increase expected to process stakes                                                           |
| min_service_stake     | array of objects    | Minimum number of nodes for specific services `{"service":"<service_id>", "min_node": <int>}`. Empty is allowed  |
| time_period           | integer             | Time in hours to consider for relay averages                                                                     |
| results_path          | string              | Path to save "What to Stake" results (empty to disable)                                                          |
| pocket_rpc            | string              | Pocket node or load balancer URL                                                                                 |
| log_level             | string              | Log level                                                                                                        |
| log_format            | string              | Log format (json or colorized text)                                                                              |
| schedule              | string              | Frequency of the "What to Stake" service calls                                                                   |
| max_workers           | integer             | Number of workers to process stake transactions in parallel                                                      |
| max_retries           | integer             | Number of retries for HTTP calls (POKTscan API or Pocket RPC)                                                    |
| max_timeout           | integer             | Timeout in milliseconds for HTTP calls                                                                           |

### FAQ

#### Do I need to restart the WTSC if I change `config.json`?

No, the service supports hot reload, configurable via the `RELOAD_SECONDS` environment variable.