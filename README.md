# proxl
[pro]metheus e[x]porter [l]inter

## Why does Proxl exist?
Proxl exists to fit a (perhaps niche) dream I have about better identification of Prometheus-format exporter metrics that don't do well in a given Prometheus environment. Prometheus is very much a system where identifying possible cardinality or quality problems in metrics _before_ they are ingested, is the right move, but not always possible.

Moreso, what works in one environment may not make sense in another. Thus, by allowing a declarative approach to "what is a problem in our Prometheus ecosystem?", Proxl aims to be a helpful tool to mitigate possible issues, before they happen. 

## How does Proxl work?
Right now, Proxl essentially does this:
- Loads a `cards.yml` file defining `scorers` and optionally, an output type (table|json)  
- Performs a (very brittle) HTTP request to a desired metrics endpoint  
- Sends each resulting MetricFamily to all defined `scorers`  
- For any scorers which evaluate the MetricFamily as problematic, output the MetricFamily name, scorer name, scorer description, and scorer criticality.  

## Key Concepts

- **MetricFamily**: `MetricFamily` is a Prometheus concept found [in their Golang client model](github.com/prometheus/client_model/go). In essence, a MetricFamily is "the combined set of labels, values, help text, and metadata for a specific metric".  
- **Scorer**: `Scorer` is a Proxl concept scattered throughout the `internal/` package. Scorers may be thought of as a construct outlining some concern that Proxl should have, along with the necessary inputs to evaluate a MetricFamily against the concern. As an example, a Proxl scorer could be of a `family_name_scorer` type, with an input of "some_metric_name". When this scorer receives a MetricFamily, if the MetricFamily name == "some_metric_name", the scorer flags it.  
- **Criteria**: `Criteria` is treated slightly differently within the `internal/` package. This is a holdover from previous iterations on the idea of Proxl, but Criteria are the heart of what computation each Scorer performs. Very confusing, yes, and I apologize <3  

One mental model which might help:

Prometheus exporter -> (is scraped by) Proxl -> (sends each MetricFamily to) a slice of Scorers -> (which evaluate) each MetricFamily against exactly one Criteria


## Usage
```
Usage:
  cards [flags]
  cards [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  debug       Scrape Prometheus exporter
  help        Help about any command

Flags:
  -h, --help              help for cards
      --loglevel string   log level to set (default "info")
  -t, --target string     HTTP endpoint to evaluate (default "http://localhost:9100/metrics")

Use "cards [command] --help" for more information about a command.
```

### Example Usage:
```
michael@plank:~/Documents/development/src/github.com/ful09003/proxl$ ./bin/cards
INFO[0000] setting log level                             fields.level=info
{"endpoint":"http://localhost:9100/metrics","level":"info","msg":"attempting scrape","time":"2021-10-26T18:28:44-07:00"}
{"level":"info","msg":"received response","res_code":200,"res_content_len_bytes":-1,"time":"2021-10-26T18:28:44-07:00"}
                               Metric                                        Scorer                                                                                                       Scorer Description         Criticality
                    node_network_info                 MAC addresses are discouraged                     MAC addresses tend to be high cardinality in our environment; if not strictly needed drop this label                   3
                    node_network_info Exceptionally lengthy metrics are discouraged                          Metrics with more than 4 labels tend to not do well in our environment; drop unnecessary labels                   3
                node_hwmon_chip_names           Exclude problematic metric families Certain metrics are known to cause confusion to users; these metrics should only be accepted if they are well-understood                   4
                      node_uname_info Exceptionally lengthy metrics are discouraged                          Metrics with more than 4 labels tend to not do well in our environment; drop unnecessary labels                   3
             node_exporter_build_info Exceptionally lengthy metrics are discouraged                          Metrics with more than 4 labels tend to not do well in our environment; drop unnecessary labels                   3
 node_cpu_scaling_frequency_max_hertz           Exclude problematic metric families Certain metrics are known to cause confusion to users; these metrics should only be accepted if they are well-understood                   4
```

## Proxl Future

Where to begin? I'd argue that there are _several_ areas where Proxl can be more robust and turn into a solid engine/library to build on top of. Right now, there's a very tight coupling between the concept of a scorer as defined by a YAML configuration, and what the core code runs - it would be lovely to decouple that. Another possibility is enabling the core code to take advantage of Golang's concurrency primatives; right now the core code executes all Scorers sequentially which is fine for a 1:1 invocation:target usage pattern, but fits well into the "what ifs" of decoupling the core code from the calling code. I'd love nothing more than some random internet stranger to also give feedback on this project and what it could/should do, so feel welcome to suggest :)
