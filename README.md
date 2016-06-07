## sensupluginses

## Commands
 * handlerElasticsearchStatus

## Usage

### handlerElasticsearchStatus
This contains a single handler that will dump the running status or every check into elasticsearch allowing developers the freedom
to create their own dashboards using Kibana or any other tool they wish.

Check specific stats for checkChronyStats

Ex. `./sensupluginses handlerElasticsearchStatus --port --host --index`

## Installation

1. godep go build -o bin/sensupluginses
1. chmod +x sensupluginses
1. cp sensupluginses /usr/local/bin

## Notes
