package scrape

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
)

type Scraper interface {
	// Name of the Scraper.
	Name() string

	// Help describes the role of the Scraper.
	Help() string

	// Scrape collects data from database connection and sends it over channel as prometheus metric.
	Scrape(ctx context.Context, ch chan<- prometheus.Metric) error
}
