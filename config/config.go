package config

import (
	"flag"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type LoggerConfig struct {
	IsJSON     bool   `yaml:"is_json"`
	AddSource  bool   `yaml:"add_source"`
	Level      string `yaml:"level"`
	SetFile    bool   `yaml:"set_file"`
	FileName   string `yaml:"file_name"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

type JWTConfig struct {
	SecretPath    string `yaml:"secret_path"`
	SecretHashLen int    `yaml:"secret_hash_len"`
	AccessExpAt   int    `yaml:"access_exp_at"`
	RefreshExpAt  int    `yaml:"refresh_exp_at"`
}

type PostgresConfig struct {
	Host          string `yaml:"host"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	DBName        string `yaml:"dbname"`
	Port          int    `yaml:"port"`
	SSLMode       string `yaml:"sslmode"`
	PoolMaxConns  int    `yaml:"pool_max_conns"`
	MigrationsDir string `yaml:"migrations_dir"`
	QueryTimeout  int64  `yaml:"query_timeout"`
}

type BigCacheConfig struct {
	Shards               int  `yaml:"shards"`
	LifeWindow           int  `yaml:"life_window"`
	CleanWindow          int  `yaml:"clean_window"`
	MaxEntriesInWindow   int  `yaml:"max_entries_in_window"`
	MaxEntrySize         int  `yaml:"max_entry_size"`
	HardMaxCacheSize     int  `yaml:"hard_max_cache_size"`
	Verbose              bool `yaml:"verbose"`
	IsOnRemoveWithReason bool `yaml:"is_on_remove_with_reason"`
}

type RefreshTokenConfig struct {
	QueryTimeout int64 `yaml:"query_timeout"`
	Expiration   int64 `yaml:"expiration"`
}

type RedisConfig struct {
	Addr            string             `yaml:"addr"`
	Password        string             `yaml:"password"`
	DB              int                `yaml:"db"`
	DialTimeout     int                `yaml:"dial_timeout"`
	ReadTimeout     int                `yaml:"read_timeout"`
	WriteTimeout    int                `yaml:"write_timeout"`
	PoolSize        int                `yaml:"pool_size"`
	MinIdleConns    int                `yaml:"min_idle_conns"`
	PoolTimeout     int                `yaml:"pool_timeout"`
	MaxRetries      int                `yaml:"max_retries"`
	MinRetryBackoff int                `yaml:"min_retry_backoff"`
	MaxRetryBackoff int                `yaml:"max_retry_backoff"`
	RefreshToken    RefreshTokenConfig `yaml:"refresh_token"`
}

type ClientAssets struct {
	Path               string `yaml:"path"`
	URL                string `yaml:"url"`
	Dir                string `yaml:"dir"`
	Browse             bool   `yaml:"browse"`
	Compress           bool   `yaml:"compress"`
	MaxFileSize        int64  `yaml:"max_file_size"`
	ImageQuality       int    `yaml:"image_quality"`
	IsNext             bool   `yaml:"is_next"`
	IsNextIgnoreFormat string `yaml:"is_next_ignore_format"`
}

type FileServerConfig struct {
	ClientAssets ClientAssets `yaml:"client_assets"`
	DirPerm      uint32       `yaml:"dir_perm"`
	FilePerm     uint32       `yaml:"file_perm"`
}

type CorsConfig struct {
	AllowOrigins        []string `yaml:"allow_origins"`
	AllowMethods        []string `yaml:"allow_methods"`
	AllowHeaders        []string `yaml:"allow_headers"`
	ExposeHeaders       []string `yaml:"expose_headers"`
	MaxAge              int      `yaml:"max_age"`
	AllowCredentials    bool     `yaml:"allow_credentials"`
	AllowPrivateNetwork bool     `yaml:"allow_private_network"`
}

type HTTPServerConfig struct {
	Host              string     `yaml:"host"`
	Port              int        `yaml:"port"`
	EnablePrefork     bool       `yaml:"enable_prefork"`
	EnablePrintRoutes bool       `yaml:"enable_print_routes"`
	Cors              CorsConfig `yaml:"cors"`
}

type Config struct {
	Logger     LoggerConfig     `yaml:"logger"`
	JWT        JWTConfig        `yaml:"jwt"`
	Postgres   PostgresConfig   `yaml:"postgres"`
	BigCache   BigCacheConfig   `yaml:"big_cache"`
	Redis      RedisConfig      `yaml:"redis"`
	FileServer FileServerConfig `yaml:"file_server"`
	HTTPServer HTTPServerConfig `yaml:"httpserver"`
}

// LoadConfig load config file.
func LoadConfig() string {
	var cf string

	flag.StringVar(&cf, "config", "config.yaml", "config file path")
	flag.Parse()

	return cf
}

// ParseConfig parse config file.
func ParseConfig(configFile string) (config Config, err error) {
	f, err := os.Open(configFile)
	if err != nil {
		return config, err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Printf("error closing the file: %v", err)
		}
	}(f)

	err = yaml.NewDecoder(f).Decode(&config)

	return config, err
}

// GetConfig get config.
func GetConfig() (config Config, err error) {
	cf := LoadConfig()

	return ParseConfig(cf)
}
