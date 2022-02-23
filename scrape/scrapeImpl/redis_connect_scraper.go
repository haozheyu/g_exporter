package scrapeImpl

import (
	"bufio"
	"context"
	"g_exporter/global"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

const (
	redisConnSubs = "redis_connect"
)


// Exporter implements the prometheus.Exporter interface, and exports Redis metrics.
type RedisExporter struct {
	metricMapCounters map[string]string
	metricMapGauges   map[string]string
}

func (RedisExporter) Name() string {
	return redisConnSubs
}

func (RedisExporter) Help() string {
	return "redis connect info"
}

func NewRedisExporter() *RedisExporter{
	return &RedisExporter{
		metricMapCounters: make(map[string]string),
		metricMapGauges:   make(map[string]string),
	}
}

func (RedisExporter) Scrape(ctx context.Context, ch chan<- prometheus.Metric) error {
	ip := global.Option["redis_host"]
	port := global.Option["redis_port"]
	pass := global.Option["redis_password"]
	client := global.RedisClient(ip, port, pass)
	res, err := client.Do(ctx, "info", "all").Text()
	infoArgs := parseKeyArg(res, ":")
    // server
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(redisConnSubs, "server_time_usec", "server_time_usec"),
		prometheus.CounterValue,
		global.StrToFloat64(infoArgs["server_time_usec"]),
	)
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(redisConnSubs, "uptime_in_seconds", "uptime_in_seconds"),
		prometheus.CounterValue,
		global.StrToFloat64(infoArgs["uptime_in_seconds"]),
	)
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(redisConnSubs, "uptime_in_days", "uptime_in_days"),
		prometheus.CounterValue,
		global.StrToFloat64(infoArgs["uptime_in_days"]),
	)
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(redisConnSubs, "io_threads_active", "io_threads_active"),
		prometheus.GaugeValue,
		global.StrToFloat64(infoArgs["io_threads_active"]),
	)
    // clients
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(redisConnSubs, "connected_clients", "connected_clients"),
		prometheus.GaugeValue,
		global.StrToFloat64(infoArgs["connected_clients"]),
	)
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(redisConnSubs, "maxclients", "maxclients"),
		prometheus.CounterValue,
		global.StrToFloat64(infoArgs["maxclients"]),
	)
	// memory
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(redisConnSubs, "used_memory", "used_memory bytes"),
		prometheus.GaugeValue,
		global.StrToFloat64(infoArgs["used_memory"]),
	)
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(redisConnSubs, "total_system_memory", "total_system_memory bytes"),
		prometheus.GaugeValue,
		global.StrToFloat64(infoArgs["total_system_memory"]),
	)
	// Persistence
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(redisConnSubs, "rdb_last_cow_size", "rdb_last_cow_size bytes"),
		prometheus.GaugeValue,
		global.StrToFloat64(infoArgs["rdb_last_cow_size"]),
	)
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(redisConnSubs, "aof_rewrite_in_progress", "aof_rewrite_in_progress bytes"),
		prometheus.GaugeValue,
		global.StrToFloat64(infoArgs["aof_rewrite_in_progress"]),
	)
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(redisConnSubs, "current_cow_size", "current_cow_size bytes"),
		prometheus.GaugeValue,
		global.StrToFloat64(infoArgs["current_cow_size"]),
	)
	// Stats
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(redisConnSubs, "total_connections_received", "total_connections_received bytes"),
		prometheus.GaugeValue,
		global.StrToFloat64(infoArgs["total_connections_received"]),
	)
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(redisConnSubs, "total_commands_processed", "total_commands_processed bytes"),
		prometheus.GaugeValue,
		global.StrToFloat64(infoArgs["total_commands_processed"]),
	)
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(redisConnSubs, "total_net_input_bytes", "total_net_input_bytes bytes"),
		prometheus.GaugeValue,
		global.StrToFloat64(infoArgs["total_net_input_bytes"]),
	)
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(redisConnSubs, "total_net_output_bytes", "total_net_output_bytes bytes"),
		prometheus.GaugeValue,
		global.StrToFloat64(infoArgs["total_net_output_bytes"]),
	)
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(redisConnSubs, "expired_keys", "expired_keys"),
		prometheus.GaugeValue,
		global.StrToFloat64(infoArgs["expired_keys"]),
	)

	// cpu
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(redisConnSubs, "used_cpu_sys", "used_cpu_sys 1c=1000"),
		prometheus.GaugeValue,
		global.StrToFloat64(infoArgs["used_cpu_sys"]),
	)
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(redisConnSubs, "used_cpu_user", "used_cpu_user 1c=1000"),
		prometheus.GaugeValue,
		global.StrToFloat64(infoArgs["used_cpu_user"]),
	)
	// cluster
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(redisConnSubs, "cluster_enabled", "cluster_enabled"),
		prometheus.CounterValue,
		global.StrToFloat64(infoArgs["cluster_enabled"]),
	)
	// Keyspace
	//ch <- prometheus.MustNewConstMetric(
	//	global.NewDesc(redisConnSubs, "db0", "db0 keys number"),
	//	prometheus.CounterValue,
	//	global.StrToFloat64(infoArgs["db0"]),
	//)

	return err
}

func parseKeyArg(keysArgString,split string) map[string]string {
	scanner := bufio.NewScanner(strings.NewReader(keysArgString))
	item := make(map[string]string)
	for scanner.Scan(){
		strarr := strings.Split(scanner.Text(), split)
		if len(strarr)>1 {
			item[strarr[0]] = strarr[1]
		}
	}
	return item
}