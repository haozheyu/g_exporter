package main

import (
	"g_exporter/collector"
	"g_exporter/global"
	"g_exporter/scrape"
	"g_exporter/scrape/scrapeImpl"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"net/http"
)

func init() {
	prometheus.MustRegister(version.NewCollector(global.Namespace + "_exporter"))

}
func main() {

	enabledScrapers := []scrape.Scraper{
		scrapeImpl.MysqlConnectScraper{},
	}
	handlerFunc := newHandler(collector.NewMetrics(), enabledScrapers)
	http.Handle("/metrics", promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, handlerFunc))
	srv := &http.Server{Addr: ":8000"}
	web.ListenAndServe(srv, "", promlog.New(&promlog.Config{}))

}

func newHandler(metrics collector.Metrics, scrapers []scrape.Scraper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mydsn := "ip:pord?username&password/url"
		ctx := r.Context()
		registry := prometheus.NewRegistry()
		registry.MustRegister(collector.New(ctx, mydsn, metrics, scrapers))
		gatherers := prometheus.Gatherers{
			prometheus.DefaultGatherer,
			registry,
		}
		h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	}
}
