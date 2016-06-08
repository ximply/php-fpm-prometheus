# php-fpm-prometheus

Simple [PHP-FPM](http://php.net/manual/en/install.fpm.php) status exporter for [Prometheus](https://prometheus.io/).

## Installation

If you are using Go 1.6+ (or 1.5 with the `GO15VENDOREXPERIMENT=1` environment variable), you can install `php-fpm-prometheus` with the following command:

```bash
$ go get -u github.com/peakgames/php-fpm-prometheus
```

## Usage

```bash
$ ./php-fpm-prometheus --help
Usage of ./php-fpm-prometheus:
  -addr string
        IP/port for the HTTP server (default "0.0.0.0:8080")
  -status-url string
        PHP-FPM status URL

$ ./php-fpm-prometheus -status-url "http://example.com/status" -addr "127.0.0.1:8080"
```

Finally, point Prometheus to `http://127.0.0.1:8080/metrics`.

## Contributing

All contributions are welcome, but if you are considering significant changes, please open issue beforehand and discuss it with us.

## License

MIT. See the `LICENSE` file for more information.