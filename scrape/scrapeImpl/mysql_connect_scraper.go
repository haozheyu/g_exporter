package scrapeImpl

import (
	"context"
	"fmt"
	"g_exporter/global"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"time"
)

const (
	// Subsystem.
	mysqlConnSubs = "mysql_connect"
)

var (
	mysqlConnDesc = prometheus.NewDesc(
		prometheus.BuildFQName(global.Namespace, mysqlConnSubs, "info"),
		"mysql connect info .",
		[]string{""}, nil,
	)
	dbClient          = global.InitConnect()
	queryVersionSQL   = "select version() as version limit 1"
	queryStatusSQL    = "show global status"
	queryVariablesSQL = "show global variables"
)

type MysqlConnectScraper struct{}

// Name of the Scraper. Should be unique.
func (MysqlConnectScraper) Name() string {
	return mysqlConnSubs
}

// Help describes the role of the Scraper.
func (MysqlConnectScraper) Help() string {
	return "mysql connect info"
}
func (MysqlConnectScraper) Scrape(ctx context.Context, dc string, ch chan<- prometheus.Metric) error {
	ip := global.Option["mysql_host"]
	port := global.Option["mysql_port"]
	pass := global.Option["mysql_password"]
	user := global.Option["mysql_user"]
	mydb, err := global.Connect(ip, port, user, pass, "information_schema")
	row := mydb.QueryRow(queryVersionSQL)
	var version string
	if err := row.Scan(&version); err != nil {
		log.Error(fmt.Sprintf("Can't scan mysql version on %s:%d, %s", ip, port, err))
	}
	fmt.Println("==============" + version + "==============")
	rows, err := mydb.Query(queryStatusSQL)
	if err != nil {
		log.Error(fmt.Sprintf("Can't query mysql status on %s:%d, %s", ip, port, err))
	}
	defer rows.Close()
	var key, value string
	globalStatusPrev := make(map[string]string)
	for rows.Next() {
		err := rows.Scan(&key, &value)
		if err != nil {
			log.Error(fmt.Sprintf("Can't scan mysql status on %s:%d, %s", ip, port, err))
		}
		globalStatusPrev[key] = value
	}

	time.Sleep(time.Duration(1) * time.Second)

	rows, err = mydb.Query(queryStatusSQL)
	if err != nil {
		log.Error(fmt.Sprintf("Can't query mysql status on %s:%d, %s", ip, port, err))
	}
	defer rows.Close()
	globalStatus := make(map[string]string)
	for rows.Next() {
		err := rows.Scan(&key, &value)
		if err != nil {
			log.Error(fmt.Sprintf("Can't scan mysql status on %s:%s, %s", ip, port, err))
		}
		globalStatus[key] = value
	}

	rows, err = mydb.Query(queryVariablesSQL)
	if err != nil {
		log.Error(fmt.Sprintf("Can't query mysql variables on %s:%s, %s", ip, port, err))
	}
	defer rows.Close()
	globalVariables := make(map[string]string)
	for rows.Next() {
		err := rows.Scan(&key, &value)
		if err != nil {
			log.Error(fmt.Sprintf("Can't scan mysql variables on %s:%s, %s", ip, port, err))
		}
		globalVariables[key] = value
	}

	//connect := 1
	//variables
	timezone := globalVariables["time_zone"]
	readonly := globalVariables["read_only"]
	hostname := globalVariables["hostname"]
	keyBufferSize := globalVariables["key_buffer_size"]
	sortBufferSize := globalVariables["sort_buffer_size"]
	joinBufferSize := globalVariables["join_buffer_size"]
	maxConnections := globalVariables["max_connections"]
	openFilesLimit := globalVariables["open_files_limit"]
	tableOpenCache := globalVariables["table_open_cache"]
	fmt.Println(timezone, readonly, hostname, keyBufferSize, sortBufferSize, joinBufferSize,
		maxConnections, openFilesLimit, tableOpenCache,
	)

	//uptime := global.StrToInt(globalStatus["Uptime"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "uptime", "uptime second"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Uptime"])),
	)
	//openFiles := global.StrToInt(globalStatus["open_files"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "open_files", "open_files number"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["open_files"])),
	)
	//openTables := global.StrToInt(globalStatus["Open_tables"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "Open_tables", "Open_tables number"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Open_tables"])),
	)
	//threadsConnected := global.StrToInt(globalStatus["Threads_connected"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "Threads_connected", "Threads_connected number"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Threads_connected"])),
	)
	//threadsRunning := global.StrToInt(globalStatus["Threads_running"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "Threads_running", "Threads_running number"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Threads_running"])),
	)
	//threadsCreated := global.StrToInt(globalStatus["Threads_created"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "Threads_created", "Threads_created number"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Threads_created"])),
	)
	//threadsCached := global.StrToInt(globalStatus["Threads_cached"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "Threads_cached", "Threads_cached number"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Threads_cached"])),
	)
	//connections := global.StrToInt(globalStatus["Connections"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "Connections", "Connections number"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Connections"])),
	)
	//abortedClients := global.StrToInt(globalStatus["Aborted_clients"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "Aborted_clients", "Aborted_clients number"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Aborted_clients"])),
	)
	//abortedConnects := global.StrToInt(globalStatus["Aborted_connects"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "Aborted_connects", "Aborted_connects number"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Aborted_connects"])),
	)
	//
	//bytesReceived := global.StrToInt(globalStatus["Bytes_received"]) - global.StrToInt(globalStatusPrev["Bytes_received"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "bytesReceived", "bytesReceived"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["bytesReceived"])-global.StrToInt(globalStatusPrev["Bytes_received"])),
	)
	//bytesSent := global.StrToInt(globalStatus["Bytes_sent"]) - global.StrToInt(globalStatusPrev["Bytes_sent"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "bytesSent", "bytesSent"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Bytes_sent"])-global.StrToInt(globalStatusPrev["Bytes_sent"])),
	)
	//comSelect := global.StrToInt(globalStatus["Com_select"]) - global.StrToInt(globalStatusPrev["Com_select"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "comSelect", "comSelect"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Com_select"])-global.StrToInt(globalStatusPrev["Com_select"])),
	)
	//comInsert := global.StrToInt(globalStatus["Com_insert"]) - global.StrToInt(globalStatusPrev["Com_insert"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "comInsert", "comInsert"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Com_insert"])-global.StrToInt(globalStatusPrev["Com_insert"])),
	)
	//comUpdate := global.StrToInt(globalStatus["Com_update"]) - global.StrToInt(globalStatusPrev["Com_update"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "comUpdate", "comUpdate"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Com_update"])-global.StrToInt(globalStatusPrev["Com_update"])),
	)
	//comDelete := global.StrToInt(globalStatus["Com_delete"]) - global.StrToInt(globalStatusPrev["Com_delete"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "comDelete", "comDelete"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Com_delete"])-global.StrToInt(globalStatusPrev["Com_delete"])),
	)
	//comCommit := global.StrToInt(globalStatus["Com_commit"]) - global.StrToInt(globalStatusPrev["Com_commit"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "comCommit", "comCommit"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Com_commit"])-global.StrToInt(globalStatusPrev["Com_commit"])),
	)
	//comRollback := global.StrToInt(globalStatus["Com_rollback"]) - global.StrToInt(globalStatusPrev["Com_rollback"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "comRollback", "comRollback"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Com_rollback"])-global.StrToInt(globalStatusPrev["Com_rollback"])),
	)
	//questions := global.StrToInt(globalStatus["Questions"]) - global.StrToInt(globalStatusPrev["Questions"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "questions", "questions"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Questions"])-global.StrToInt(globalStatusPrev["Questions"])),
	)
	//queries := global.StrToInt(globalStatus["Queries"]) - global.StrToInt(globalStatusPrev["Queries"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "queries", "queries"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Queries"])-global.StrToInt(globalStatusPrev["Queries"])),
	)
	//slowQueries := global.StrToInt(globalStatus["Slow_queries"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "slowQueries", "slowQueries"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Slow_queries"])),
	)
	//
	////innodb status
	//innodbPagesCreated := global.StrToInt(globalStatus["Innodb_pages_created"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "innodbPagesCreated", "innodbPagesCreated"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Innodb_pages_created"])),
	)
	//innodbPagesRead := global.StrToInt(globalStatus["Innodb_pages_read"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "innodbPagesRead", "innodbPagesRead"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Innodb_pages_read"])),
	)
	//innodbPagesWritten := global.StrToInt(globalStatus["Innodb_pages_written"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "innodbPagesWritten", "innodbPagesWritten"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Innodb_pages_written"])),
	)
	//innodbRowLockCurrentWaits := global.StrToInt(globalStatus["Innodb_row_lock_current_waits"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "innodbRowLockCurrentWaits", "innodbRowLockCurrentWaits"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Innodb_row_lock_current_waits"])),
	)
	//innodbBufferPoolReadRequests := global.StrToInt(globalStatus["Innodb_buffer_pool_read_requests"]) - global.StrToInt(globalStatusPrev["Innodb_buffer_pool_read_requests"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "innodbBufferPoolReadRequests", "innodbBufferPoolReadRequests"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Innodb_buffer_pool_read_requests"])-global.StrToInt(globalStatusPrev["Innodb_buffer_pool_read_requests"])),
	)
	//innodbBufferPoolWriteRequests := global.StrToInt(globalStatus["Innodb_buffer_pool_write_requests"]) - global.StrToInt(globalStatusPrev["Innodb_buffer_pool_write_requests"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "innodbBufferPoolWriteRequests", "innodbBufferPoolWriteRequests"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Innodb_buffer_pool_write_requests"])-global.StrToInt(globalStatusPrev["Innodb_buffer_pool_write_requests"])),
	)
	//innodbRowsDeleted := global.StrToInt(globalStatus["Innodb_rows_deleted"]) - global.StrToInt(globalStatusPrev["Innodb_rows_deleted"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "innodbRowsDeleted", "innodbRowsDeleted"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Innodb_rows_deleted"])-global.StrToInt(globalStatusPrev["Innodb_rows_deleted"])),
	)
	//innodbRowsInserted := global.StrToInt(globalStatus["Innodb_rows_inserted"]) - global.StrToInt(globalStatusPrev["Innodb_rows_inserted"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "innodbRowsInserted", "innodbRowsInserted"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Innodb_rows_inserted"])-global.StrToInt(globalStatusPrev["Innodb_rows_inserted"])),
	)
	//innodbRowsRead := global.StrToInt(globalStatus["Innodb_rows_read"]) - global.StrToInt(globalStatusPrev["Innodb_rows_read"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "innodbRowsRead", "innodbRowsRead"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Innodb_rows_read"])-global.StrToInt(globalStatusPrev["Innodb_rows_read"])),
	)
	//innodbRowsUpdated := global.StrToInt(globalStatus["Innodb_rows_updated"]) - global.StrToInt(globalStatusPrev["Innodb_rows_updated"])
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(mysqlConnSubs, "innodbRowsUpdated", "innodbRowsUpdated"),
		prometheus.UntypedValue,
		float64(global.StrToInt(globalStatus["Innodb_rows_updated"])-global.StrToInt(globalStatusPrev["Innodb_rows_updated"])),
	)
	return nil
}
