package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	MigrateSource string        `yaml:"migrate_source"`
	RedisSource   string        `yaml:"redis_source"`
	LimitCategory int           `yaml:"limit_category"`
	CacheTTL      time.Duration `yaml:"cache_ttl"`
	Database      `yaml:"database"`
	HTTPServer    `yaml:"http_server"`
	Logs          `yaml:"logs"`
}

type HTTPServer struct {
	Address        string        `yaml:"address" env-default:"0.0.0.0:8080"`
	MaxHeaderBytes int           `yaml:"max_header_bytes"`
	ReadTimeOut    time.Duration `yaml:"read_timeout"`
	WriteTimeOut   time.Duration `yaml:"write_timeout"`
	IdleTimeout    time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Logs struct {
	Env        string `yaml:"env" env-default:"local"`
	InfoPath   string `yaml:"info_path"`
	DebugPath  string `yaml:"debug_path"`
	ErrorPath  string `yaml:"error_path"`
	WarnPath   string `yaml:"warn_path"`
	AccessPath string `yaml:"access_path"`
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
}

type Database struct {
	Source            string        `yaml:"source"`
	MaxConns          int32         `yaml:"max_conns" env-default:"10"`
	MinConns          int32         `yaml:"min_conns" env-default:"0"`
	MaxConnLifeTime   time.Duration `yaml:"max_conn_life_time" env-default:"30m"`
	MaxConnIdleTime   time.Duration `yaml:"max_conn_idle_time" env-default:"15m"`
	HealthCheckPeriod time.Duration `yaml:"health_check_period" env-default:"1m"`
	ConnectTimeout    time.Duration `yaml:"connect_timeout" env-default:"5m"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}

	configPath := os.Getenv("CONFIG_PATH")

	var cfg Config

	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err.Error())
	}

	return &cfg
}
