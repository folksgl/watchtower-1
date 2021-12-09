package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace = "watchtower"

var configPath = flag.String("config", "config.yaml", "Path to configuration file.")
var validationInterval = flag.Int("interval", int(DetectionInterval.Seconds()), "The interval (in seconds) that Watchtower will run validation checks and update exported metrics")
var configString = ""

var (
	// Counters for failed/successful validation checks
	failedAppChecks = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "app_checks",
		Name:      "failed_total",
		Help:      "Number of times the config refresh for V3Apps has failed for any reason",
	})
	successfulAppChecks = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "app_checks",
		Name:      "success_total",
		Help:      "Number of times the config refresh for V3Apps has succeeded",
	})
	failedSpaceChecks = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "space_checks",
		Name:      "failed_total",
		Help:      "Number of times the config check for Spaces has failed for any reason",
	})
	successfulSpaceChecks = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "space_checks",
		Name:      "success_total",
		Help:      "Number of times the config check for Spaces has succeeded",
	})

	// Counters for unknown/missing/misconfigured resources
	totalUnknownApps = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "unknown",
		Name:      "apps_total",
		Help:      "Number of Apps deployed that are not in the allowed config file (config.yaml)",
	})
	totalMissingApps = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "missing",
		Name:      "apps_total",
		Help:      "Number of Apps in the provided config file that are not deployed",
	})

	totalSpaceSSHViolations = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "ssh",
		Name:      "space_misconfiguration_total",
		Help:      "Number of Spaces that have misconfigured SSH access settings",
	})
)

func configHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, configString)
}

func main() {
	flag.Parse()
	NewDetector(configPath, *validationInterval)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/config", configHandler)
	log.Fatal(http.ListenAndServe(":"+ReadPortFromEnv(), nil))
}
