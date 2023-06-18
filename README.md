# Dunwich

A data-centre focused load tester for saturating high-throughput networks. On 20G+ networks, standard tools such as iperf can become CPU bound. Dunwich aims to leverage the concurrent nature of Go to bypass these limitations.

## Why "Dunwich"?
Most load testing apps simulate traffic using generated content of uniform size. Dunwich, on the other hands, aim to offer a more realistic load simulation by sharing paragraphs of H.P. Lovecraft's work _The Dunwich Horror_ (pulled from [Project Gutenberg](https://www.gutenberg.org/ebooks/50133)).

## Configuration
| Environment Variable | Description                                          | Default |
|----------------------|------------------------------------------------------|---------|
| DUNWICH_LOGGER_LEVEL | The log level to use (e.g. Info, Debug, Trace, etc.) | Info    |

## Development Setup
1. Make sure Go 1.20+ is installed.
