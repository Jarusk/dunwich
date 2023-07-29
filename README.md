# Dunwich (Work In Progress)

A data-centre focused load tester for saturating high-throughput networks. On 20G+ networks, standard tools such as iperf can become CPU bound. Dunwich aims to leverage the concurrent nature of Go to bypass these limitations.

**NOT IN A FUNCTIONAL STATE YET**

## Why "Dunwich"?
Most load testing apps simulate traffic using generated content of uniform size. Dunwich, on the other hands, aim to offer a more realistic load simulation by sharing paragraphs of H.P. Lovecraft's work _The Dunwich Horror_ (pulled from [Project Gutenberg](https://www.gutenberg.org/ebooks/50133)).

## Configuration
| Environment Variable | Description                                          | Required | Default |
|----------------------|------------------------------------------------------|----------|---------|
| DUNWICH_LOG | The log level to use (e.g. Info, Debug, Trace, etc.) | False | Info |
| DUNWICH_CLUSTER_PORT | The port to use for gossip between nodes | False | 7946 |
| DUNWICH_JOIN | The address:port list of nodes to bootstrap to | False | N/A |
| DUNWICH_PPROF | Enable the pprof server | False | N/A |


## Development Setup
1. Make sure Go 1.20+ is installed.
