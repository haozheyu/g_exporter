package global

import (
	"flag"
	"github.com/larspensjo/config"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"strconv"
	"strings"
)

const (
	Namespace = "general"
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

func StrToFloat64(str string) float64 {
	result, _ := strconv.ParseFloat(str, 32)
	return result
}
