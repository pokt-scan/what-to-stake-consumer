# What To Stake?

This tool provides optimal* deployments of a series of nodes in the POKT Network.

### Table of Contents

- [Optimization Process](#optimization-process)
- [Important Warning](#important-warning)
- [Requirements](#requirements)
- [How To Use](#how-to-use)
- [Using the Makefile](#using-the-makefile)
- [Understanding `config.json` Values](#understanding-configjson-values)
- [Available Environment Variables](#available-environment-variables)
- [FAQ](#faq)

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

### IMPORTANT WARNING

⚠️ **Recommendation:** Before adding keys, we strongly advise you to first use the tool with `dry_mode` set to `true`. This will allow you to understand what the tool is going to do without making any actual changes. ⚠️

### Requirements

1. POKTscan API Token. Get yours at: [POKTscan API Token](https://poktscan.com/teams?tab=api_token)

### How To Use

1. Create a copy of `config.json.sample` and name it `config.json`.
2. Edit `config.json` with your desired values. **Make sure to set `dry_mode` to `true` initially to understand the tool's behavior without making any changes.**
3. Generate the GraphQL types by running: `make generate`.
4. Start the service with Docker:
    ```sh
    make start_docker
    ```
   An existing Docker image build is available at [Docker Hub](https://hub.docker.com/repository/docker/poktscan/wtsc/general).

If you want/need to modify the path and name of the config file, please use `CONFIG_FILE` to override the default `./config.json`.

### Using the Makefile

The provided `Makefile` includes several targets that help manage the project lifecycle, including generating code, building the project, and managing Docker containers. Below are the available targets and how to use them:

#### Targets

- **generate**: Generates POKTscan GraphQL schema types.
  ```sh
  make generate
  ```

- **build**: Builds the WTSC project after generating the necessary types.
  ```sh
  make build
  ```

- **build_docker**: Builds the WTSC Docker image using Docker Compose.
  ```sh
  make build_docker
  ```

- **build_docker_no_cache**: Builds the WTSC Docker image without using the cache.
  ```sh
  make build_docker_no_cache
  ```

- **start**: Builds and starts the WTSC project on the host.
  ```sh
  make start
  ```

- **start_docker**: Starts the WTSC project in a Docker container using Docker Compose (without detaching the terminal).
  ```sh
  make start_docker
  ```

- **start_as_daemon**: Starts the WTSC project in a Docker container using Docker Compose (detached mode).
  ```sh
  make start_as_daemon
  ```

- **stop_docker**: Stops the WTSC Docker container without destroying it.
  ```sh
  make stop_docker
  ```

- **down**: Stops and removes the WTSC Docker container and associated volumes.
  ```sh
  make down
  ```

### Understanding `config.json` Values

| Parameter            | Type             | Description                                                                                                          |
|----------------------|------------------|----------------------------------------------------------------------------------------------------------------------|
| dry_mode             | boolean          | If true, calls to "What to Stake" will write results to `results_path` or print them if `results_path` is empty.     |
| poktscan_api         | string           | POKTscan API endpoint                                                                                                |
| poktscan_api_token   | string           | POKTscan API token                                                                                                   |
| network_id           | string           | Can be "mainnet" or "testnet"                                                                                        |
| tx_memo              | string           | Transaction memo                                                                                                     |
| tx_fee               | integer          | Transaction fee                                                                                                      |
| domain               | string           | Your node's domain (e.g., "poktscan.cloud" or "c0d3r.org")                                                           |
| service_pool         | array of strings | Service IDs (aka chain on morse) available in your fleet                                                             |
| servicer_keys        | array of strings | List of private keys for signing stake transactions                                                                  |
| stake_weight         | integer          | Used by the "What to Stake" service to estimate potential rewards (1-4)                                              |
| min_increase_percent | integer          | Minimum percentage increase expected to process stakes                                                               |
| min_service_stake    | array of objects | Minimum number of nodes for specific services `{"service":"<service_id>", "min_node": <int>}`. Empty is allowed      |
| time_period          | integer          | Time in hours to consider for relay averages                                                                         |
| results_path         | string           | Path to save "What to Stake" results (empty to disable)                                                              |
| pocket_rpc           | string           | Pocket node or load balancer URL                                                                                     |
| log_level            | string           | Log level                                                                                                            |
| log_format           | string           | Log format (json or colorized text)                                                                                  |
| schedule             | string           | Frequency of the "What to Stake" service calls                                                                       |
| run_once_at_start    | boolean          | If true, the WTSC evaluation job runs immediately at startup, then follows the schedule. If false, the first job runs according to the schedule parameter. |
| max_workers          | integer          | Number of workers to process stake transactions in parallel                                                          |
| max_retries          | integer          | Number of retries for HTTP calls (POKTscan API or Pocket RPC)                                                        |
| max_timeout          | integer          | Timeout in milliseconds for HTTP calls                                                                               |

#### Environment Variables

To customize the behavior of the Makefile commands, you can set the following environment variables:

- **PROJECT_ROOT**: Override the current working directory.
- **CONFIG_FILE**: Override the default config file name `config.json`.
- **VERSION**: Specify the project version.

For example, to use a custom config file:

```sh
CONFIG_FILE=custom_config.json make build
```

### FAQ

#### Do I need to restart the WTSC if I change `config.json`?

No, the service supports hot reload, configurable via the `RELOAD_SECONDS` environment variable.

#### What happens in `dry_mode`?

When `dry_mode` is set to `true`, the tool simulates its operations and outputs the potential changes without making any actual stakes. This allows you to preview the recommended adjustments and understand their impact without modifying your node configurations.

#### How often should I adjust my stakes?

This depends on the network dynamics and the changes in relay patterns. Given that the network and other participants' behaviors are constantly evolving, it is advisable to review the recommendations periodically and adjust your stakes accordingly.

#### Can I customize the logging format?

Yes, you can customize the log format by setting the `log_format` parameter in `config.json`. Available options are `json` for JSON formatted logs and `text` for colorized text logs.

#### What does the `schedule` parameter do?

The `schedule` parameter specifies how often the "What to Stake" service should run. This is useful for automating the regular review and adjustment of your node stakes based on up-to-date network data.

#### How do I override the default configuration file path?

You play on this using `PROJECT_ROOT` and `CONFIG_FILE` to override working directory and config file name.

#### What should I do if the tool is not producing expected results?

Ensure that your configuration is correct and up-to-date. Also, consider checking the logs for any errors or warnings that might indicate issues. If the problem persists, you can reach out to the support or community forums for assistance.