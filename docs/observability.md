# Observability

Ghostify provides metrics + health/status commands:

* `ghostify health -c <config>` starts the Sing-box flow for 5 seconds to exercise the runtime and record success/failure.
* `ghostify status` prints the metrics snapshot plus active protocols and bridge runners.
* `ghostify status-server -addr 127.0.0.1:9111` exposes `/status` over HTTP so external monitors can poll realtime state.
* `ghostify run` already logs structured lifecycle events like `ghostify runtime running` and `shutdown signal received`.

Metrics capture run/health outcomes in `internal/observability/metrics.go` for future aggregation.
