# What To Stake ?

This tool is designed to provide optimal* deployments of a series of nodes in the POKT Network.

### On the optimization process

The optimization process has **many** important details that the user should be aware of before using the provided recommendations:

- The optimization process is based on basic model of the network since not all data is available (i.e. the gateways selection mechanisms).
- We do not take into account nodes QoS, we treat all nodes as "equals", if your nodes are good, your rewards will be higher than predicted, if your nodes are bad, the other way around.
- We optimize for "round-robin" relay distribution.
- We don't know were your nodes are located and were the traffic is. This affects your results greatly. We cannot tell you were to deploy your chain nodes.
- The optimization is **domain based**, we are not going to give you your optimal deployment for a balanced network, we will tell you the optimal deployment for the **current** network.
- Since the network changes all the time and other participants will act based on your changes, be prepared to change your stakes periodically.
- We expect that if enough node runners use this tool the network state will converge (no more updates needed), but **we cannot guarantee convergence**.
- It can happen that recommendations begin to oscillate (a proposed strategy leads to a previous stake strategy). We don't expect this to happen, but if so, we told you so. This could be due to the adversarial strategy of different node runners.
- By default we look back 24Hs to get averages of relays in a service. The time used to calculate the average relays in a service greatly affects outcomes. For example, if a chain had no relays in the last 24Hs it won't be used during optimization (no relays = no rewards).

**The results of our optimizations are optimal given the simplifications we commented on above.** 