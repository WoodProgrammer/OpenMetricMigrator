<div><img src="docs/img/image.png" width="100"/><h1>Open Metric Converter Prometheus</h1></div>

OpenMetricMigrator is a tool designed to migrate and transform OpenMetrics data efficiently. It simplifies the process of converting metric formats and moving data between monitoring systems.

## Features

- Convert OpenMetrics data to various formats
- Seamlessly migrate metrics between monitoring platforms
- Easy configuration with minimal setup
- High performance and scalable

## Installation

To install OpenMetricMigrator, clone the repository and build the project:

```sh
# Clone the repository
git clone https://github.com/WoodProgrammer/OpenMetricMigrator.git
cd OpenMetricMigrator

go build -o opm .

mv opm /usr/local/bin
```

Alternatively, you can install it using Go:

```sh
go install github.com/WoodProgrammer/OpenMetricMigrator@latest
```

## Usage

Run the tool with the required options:

```sh
./openmetricmigrator --input input-file.prom --output output-file.prom
```

### Available Flags


```sh

CLI tool to export Prometheus data in OpenMetrics format 

Usage:
  promcli [flags]

Flags:
  -d, --directory string   Data directory to export (default "data")
  -e, --end string         End timestamp (epoch) (default "0")
  -h, --help               help for promcli
  -H, --host string        Prometheus host (default "localhost")
  -P, --port string        Prometheus port (default "9090")
  -q, --query string       PromQL query
  -s, --start string       Start timestamp (epoch) (default "0")
  -t, --step string        Query step (default "15s")

```

## Example

Convert an OpenMetrics file to a Prometheus-compatible format:

```sh
./prom-migrator -H localhost -P 9090 -s 1741484483 -e 1741488083 -q 'up{job="prometheus"}'
```

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests to improve the tool.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

For any questions or feedback, feel free to open an issue on GitHub.