# cards
Cardinality Scorer, for Prometheus Exporters

## Why does this exist?
Cards exists to fulfill a dream I once had to add automation to spot-checking Prometheus exporters for (frankly) shitty cardinality or other problems. This project isn't there yet and may never be, but I hope to one day be able to run an ad-hoc scrape against an exporter, run several "scoring" functions against each set of metrics returned from the exporter, and then give my human brain a nice little table or JSON output (to consume elsewhere) that flags potentially problematic metric sets.

tl;dr: cards exists as an attempt to add an independent linter for Prometheus text-format exporters, with an eye towards mitigating bad-quality metrics from entering a Prometheus instance.

## How will this work?
Right now, my idea is to:
- Perform an ad-hoc scrape (works, but is very brittle)
- Parse the response as Prometheus MetricFamily data (works)
- Run any number of configured scoring functions against the same entire metric set ('works', but read on)
- Tidy up results, send the output back to the user


## Scoring Functions
Presently, I haven't done much work to achieve my long-term goal for this mechanism. What I'd _love_ is to have a configuration-first template representing each scoring function. Scoring functions can be anything - either built-in "common" functions representing questions like "are these metrics typed?", "do these metrics have help text?", or more *generic* abstractions like "does this metric have a label value that is a UUID (potentially high-cardinality)?", "has this metric historically performed poorly in my Prometheus environment?", and so on. The former feels easier, while the latter involves potentially creating a DSL - something that appears intimidating to me as of 2021.

Another short-term problem to solve here is that right now, cards will just run every scorer defined in `cmd/realcmd/scrape.go`- there's no ability for users to select _which_ scoring functions to run. Hopefully as I make progress on the other goals with scoring functions, I'll naturally fix this up too.

## Usage
```
./bin/cards --help
cards is a utility meant to help gauge the quality of a Prometheus text-format exporter

Usage:
  cards [flags]
  cards [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  debug       Scrape Prometheus exporter
  help        Help about any command

Flags:
  -h, --help            help for cards
  -m, --mult int        Set to maximum estimated hosts being scraped by Prometheus (default 1)
  -t, --target string   HTTP endpoint to evaluate (default "http://localhost:9100/metrics")

Use "cards [command] --help" for more information about a command.
```

### Example Usage:
```
michael@plank:~/Documents/development/src/github.com/ful09003/cards$ ./bin/cards 
INFO[0000] attempting scrape                             endpoint="http://localhost:9100/metrics"
INFO[0000] received response                             res_code=200 res_content_len_bytes=-1
                                      Metric               Score
               node_network_name_assign_type                   3
                node_power_supply_cyclecount                   1
                          node_procs_running                   0
        node_schedstat_running_seconds_total                   8
                node_softnet_processed_total                   8
            go_memstats_last_gc_time_seconds                   0
             node_memory_FileHugePages_bytes                   0
                      node_thermal_zone_temp                  22
...
```

```
michael@plank:~/Documents/development/src/github.com/ful09003/cards$ ./bin/cards debug | grep node_thermal_zone_temp
INFO[0000] attempting scrape                             endpoint="http://localhost:9100/metrics"
INFO[0000] received response                             res_code=200 res_content_len_bytes=-1
# HELP node_thermal_zone_temp Zone temperature in Celsius
# TYPE node_thermal_zone_temp gauge
node_thermal_zone_temp{type="B0D4",zone="9"} 39.05
node_thermal_zone_temp{type="INT3400 Thermal",zone="1"} 20
node_thermal_zone_temp{type="SEN1",zone="2"} 35.05
node_thermal_zone_temp{type="SEN2",zone="3"} 31.05
...
```