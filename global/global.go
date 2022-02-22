package global

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/larspensjo/config"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"strconv"
	"strings"
)

const (
	// Exporter Namespace.
	Namespace = "mysql_exporter"
)

func NewDesc(subsystem, name, help string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, name),
		help, nil, nil,
	)
}

// --------------------------------------------------
var (
	err        error
	configFile = flag.String("config", "config.ini", "General configuration file")
	Option     = make(map[string]string)
)

func InitConnect() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=5s&readTimeout=10s", Option["mysql_user"],
		Option["mysql_password"], Option["mysql_host"], Option["mysql_port"], Option["mysql_database"]))
	if err != nil {
		panic(fmt.Sprintln("Init mysql connect err,", err))
	}
	if err := db.Ping(); err != nil {
		panic(fmt.Sprintln("Init mysql connect err,", err))
	}
	return db
}

func Connect(host, port, username, password, database string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=3s&readTimeout=5s", username, password, host, port, database))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func Execute(db *sql.DB, sql string) (err error) {
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}
	return
}

func QueryOne(db *sql.DB, sql string) (data string, err error) {
	row := db.QueryRow(sql)
	if err := row.Scan(); err != nil {
		return "", nil
	}
	return
}

func QueryAll(db *sql.DB, sql string) ([]map[string]interface{}, error) {
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var list []map[string]interface{}
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			continue
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			v := values[i]
			b, ok := v.([]byte)
			if ok {
				entry[col] = string(b)
				//entry[col] = b
			} else {
				entry[col] = v
			}
		}
		list = append(list, entry)
	}
	return list, nil
}

func init() {
	flag.Parse()
	cfg, err := config.ReadDefault(*configFile)
	if err != nil {
		panic(err)
		os.Exit(0)
	}
	if cfg.HasSection("main") {
		section, err := cfg.SectionOptions("main")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("main", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}
	if cfg.HasSection("mysql") {
		section, err := cfg.SectionOptions("mysql")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("mysql", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}
	if cfg.HasSection("mongodb") {
		section, err := cfg.SectionOptions("mongodb")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("mongodb", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}
	if cfg.HasSection("redis") {
		section, err := cfg.SectionOptions("redis")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("redis", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}
	if cfg.HasSection("kafka") {
		section, err := cfg.SectionOptions("kafka")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("kafka", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}
	if cfg.HasSection("mail") {
		section, err := cfg.SectionOptions("mail")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("mail", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}
	if cfg.HasSection("task") {
		section, err := cfg.SectionOptions("task")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("task", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}

	if cfg.HasSection("event") {
		section, err := cfg.SectionOptions("event")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("event", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}

	if cfg.HasSection("nsq") {
		section, err := cfg.SectionOptions("nsq")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("nsq", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}

}

func StrToInt(str string) int {
	nonFractionalPart := strings.Split(str, ".")
	result, _ := strconv.Atoi(nonFractionalPart[0])
	return result
}

func StrToFloat(str string) float64 {
	result, _ := strconv.ParseFloat(str, 32)
	return result
}
