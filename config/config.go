package config

import (
	"flag"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type LoggerConfig struct {
	Level      string `yaml:"level"`
	FileName   string `yaml:"file_name"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	IsJSON     bool   `yaml:"is_json"`
	AddSource  bool   `yaml:"add_source"`
	SetFile    bool   `yaml:"set_file"`
}

type JWTConfig struct {
	SecretHashLen int    `yaml:"secret_hash_len"`
	AccessExpAt   int    `yaml:"access_exp_at"`
	RefreshExpAt  int    `yaml:"refresh_exp_at"`
	SecretPath    string `yaml:"secret_path"`
}

type ConsumerConfig struct {
	StreamName   string `yaml:"stream_name"`
	Subject      string `yaml:"subject"`
	Timeout      int    `yaml:"timeout"`
	StreamOption struct {
		Name      string `yaml:"name"`
		Storage   int    `yaml:"storage"`
		Retention int    `yaml:"retention"`
		MaxAge    int    `yaml:"max_age"`
		MaxMsgs   int64  `yaml:"max_msgs"`
		MaxBytes  int64  `yaml:"max_bytes"`
		Discard   int    `yaml:"discard"`
	} `yaml:"stream_option"`
	NatsOption struct {
		MaxReconnects        int    `yaml:"max_reconnects"`
		ReconnectWait        int    `yaml:"reconnect_wait"`
		Timeout              int    `yaml:"timeout"`
		DrainTimeout         int    `yaml:"drain_timeout"`
		PingInterval         int    `yaml:"ping_interval"`
		MaxPingsOutstanding  int    `yaml:"max_pings_outstanding"`
		ClosedHandler        bool   `yaml:"closed_handler"`
		ReconnectHandler     bool   `yaml:"reconnect_handler"`
		DisconnectErrHandler bool   `yaml:"disconnect_err_handler"`
		ErrorHandler         bool   `yaml:"error_handler"`
		Name                 string `yaml:"name"`
		ReconnectJitter      struct {
			Jitter       int `yaml:"jitter"`
			JitterForTLS int `yaml:"jitter_for_tls"`
		} `yaml:"reconnect_jitter"`
	} `yaml:"nats_option"`
	SubscribeOption struct {
		AckWait       int    `yaml:"ack_wait"`
		MaxAckPending int    `yaml:"max_ack_pending"`
		ManualAck     bool   `yaml:"manual_ack"`
		DeliverAll    bool   `yaml:"deliver_all"`
		Durable       string `yaml:"durable"`
	} `yaml:"subscribe_option"`
}

type PublisherConfig struct {
	StreamName   string `yaml:"stream_name"`
	Subject      string `yaml:"subject"`
	Timeout      int    `yaml:"timeout"`
	StreamOption struct {
		Name        string `yaml:"name"`
		Storage     int    `yaml:"storage"`
		Retention   int    `yaml:"retention"`
		MaxAge      int    `yaml:"max_age"`
		MaxMsgs     int64  `yaml:"max_msgs"`
		MaxBytes    int64  `yaml:"max_bytes"`
		Discard     int    `yaml:"discard"`
		Replicas    int    `yaml:"replicas"`
		AllowDirect bool   `yaml:"allow_direct"`
		DenyDelete  bool   `yaml:"deny_delete"`
		DenyPurge   bool   `yaml:"deny_purge"`
	} `yaml:"stream_option"`
	NatsOption struct {
		MaxReconnects   int    `yaml:"max_reconnects"`
		ReconnectWait   int    `yaml:"reconnect_wait"`
		ErrorHandler    bool   `yaml:"error_handler"`
		Name            string `yaml:"name"`
		ReconnectJitter struct {
			Jitter       int `yaml:"jitter"`
			JitterForTLS int `yaml:"jitter_for_tls"`
		} `yaml:"reconnect_jitter"`
	} `yaml:"nats_option"`
	JetStreamOption struct {
		PublishAsyncMaxPending int  `yaml:"publish_async_max_pending"`
		PublishAsyncErrHandler bool `yaml:"publish_async_err_handler"`
	} `yaml:"jet_stream_option"`
}

type NatsNotificationConfig struct {
	NatsURL   string          `yaml:"nats_url"`
	Consumer  ConsumerConfig  `yaml:"consumer"`
	Publisher PublisherConfig `yaml:"publisher"`
}

type NatsConfig struct {
	Notification NatsNotificationConfig `yaml:"notification"`
}

type PostgresConfig struct {
	QueryTimeout  int64  `yaml:"query_timeout"`
	Port          int    `yaml:"port"`
	PoolMaxConns  int    `yaml:"pool_max_conns"`
	Host          string `yaml:"host"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	DBName        string `yaml:"dbname"`
	SSLMode       string `yaml:"sslmode"`
	MigrationsDir string `yaml:"migrations_dir"`
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

type UnDeleteFileClientConfig struct {
	QueryTimeout int64 `yaml:"query_timeout"`
	Expiration   int64 `yaml:"expiration"`
}

type UnDeleteFileAchievementConfig struct {
	QueryTimeout int64 `yaml:"query_timeout"`
	Expiration   int64 `yaml:"expiration"`
}

type UnDeleteFileAwardConfig struct {
	QueryTimeout int64 `yaml:"query_timeout"`
	Expiration   int64 `yaml:"expiration"`
}

type RedisConfig struct {
	Addr                    string                        `yaml:"addr"`
	Password                string                        `yaml:"password"`
	DB                      int                           `yaml:"db"`
	DialTimeout             int                           `yaml:"dial_timeout"`
	ReadTimeout             int                           `yaml:"read_timeout"`
	WriteTimeout            int                           `yaml:"write_timeout"`
	PoolSize                int                           `yaml:"pool_size"`
	MinIdleConns            int                           `yaml:"min_idle_conns"`
	PoolTimeout             int                           `yaml:"pool_timeout"`
	MaxRetries              int                           `yaml:"max_retries"`
	MinRetryBackoff         int                           `yaml:"min_retry_backoff"`
	MaxRetryBackoff         int                           `yaml:"max_retry_backoff"`
	RefreshToken            RefreshTokenConfig            `yaml:"refresh_token"`
	UnDeleteFileClient      UnDeleteFileClientConfig      `yaml:"un_delete_file_client"`
	UnDeleteFileAchievement UnDeleteFileAchievementConfig `yaml:"un_delete_file_achievement"`
	UnDeleteFileAward       UnDeleteFileAwardConfig       `yaml:"un_delete_file_award"`
}

type ClientAssets struct {
	Path               string `yaml:"path"`
	URL                string `yaml:"url"`
	Dir                string `yaml:"dir"`
	IsNextIgnoreFormat string `yaml:"is_next_ignore_format"`
	MaxFileSize        int64  `yaml:"max_file_size"`
	ImageQuality       int    `yaml:"image_quality"`
	Browse             bool   `yaml:"browse"`
	Compress           bool   `yaml:"compress"`
	IsNext             bool   `yaml:"is_next"`
}

type AchievementAssets struct {
	Path               string `yaml:"path"`
	URL                string `yaml:"url"`
	Dir                string `yaml:"dir"`
	IsNextIgnoreFormat string `yaml:"is_next_ignore_format"`
	MaxFileSize        int64  `yaml:"max_file_size"`
	ImageQuality       int    `yaml:"image_quality"`
	Browse             bool   `yaml:"browse"`
	Compress           bool   `yaml:"compress"`
	IsNext             bool   `yaml:"is_next"`
}

type AwardAssets struct {
	Path               string `yaml:"path"`
	URL                string `yaml:"url"`
	Dir                string `yaml:"dir"`
	IsNextIgnoreFormat string `yaml:"is_next_ignore_format"`
	MaxFileSize        int64  `yaml:"max_file_size"`
	ImageQuality       int    `yaml:"image_quality"`
	Browse             bool   `yaml:"browse"`
	Compress           bool   `yaml:"compress"`
	IsNext             bool   `yaml:"is_next"`
}

type FileServerConfig struct {
	ClientAssets      ClientAssets      `yaml:"client_assets"`
	AchievementAssets AchievementAssets `yaml:"achievement_assets"`
	AwardAssets       AwardAssets       `yaml:"award_assets"`
	DirPerm           uint32            `yaml:"dir_perm"`
	FilePerm          uint32            `yaml:"file_perm"`
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

type CronConfig struct {
	UnDeleteFileClientCleaner struct {
		SleepDuration int `yaml:"sleep_duration"`
		Timeout       int `yaml:"timeout"`
	} `yaml:"un_delete_file_client_cleaner"`
	UnDeleteFileAchievementCleaner struct {
		SleepDuration int `yaml:"sleep_duration"`
		Timeout       int `yaml:"timeout"`
	} `yaml:"un_delete_file_achievement_cleaner"`
	UnDeleteFileAwardCleaner struct {
		SleepDuration int `yaml:"sleep_duration"`
		Timeout       int `yaml:"timeout"`
	} `yaml:"un_delete_file_award_cleaner"`
}

type MiddlewareConfig struct {
	ContentLengthLimiter struct {
		MaxBodySize int `yaml:"max_body_size"`
	} `yaml:"content_length_limiter"`
}

type CookieConfig struct {
	Refresh struct {
		MaxAge      int    `yaml:"max_age"`
		Name        string `yaml:"name"`
		Path        string `yaml:"path"`
		Domain      string `yaml:"domain"`
		SameSite    string `yaml:"same_site"`
		Secure      bool   `yaml:"secure"`
		HTTPOnly    bool   `yaml:"http_only"`
		Partitioned bool   `yaml:"partitioned"`
		SessionOnly bool   `yaml:"session_only"`
	} `yaml:"refresh"`
}

type IPsConfig struct {
	Allowed []string `yaml:"allowed"`
}

type SwaggerServerConfig struct {
	Host            string     `yaml:"host"`
	Cors            CorsConfig `yaml:"cors"`
	ShutdownTimeout int64      `yaml:"shutdown_timeout"`
	Port            int        `yaml:"port"`
}

type HTTPServerConfig struct {
	Host                  string     `yaml:"host"`
	Cors                  CorsConfig `yaml:"cors"`
	ShutdownTimeout       int64      `yaml:"shutdown_timeout"`
	Port                  int        `yaml:"port"`
	DisableStartupMessage bool       `yaml:"disable_startup_message"`
	EnablePrefork         bool       `yaml:"enable_prefork"`
	EnablePrintRoutes     bool       `yaml:"enable_print_routes"`
}

type Config struct {
	Logger        LoggerConfig        `yaml:"logger"`
	JWT           JWTConfig           `yaml:"jwt"`
	Nats          NatsConfig          `yaml:"nats"`
	Postgres      PostgresConfig      `yaml:"postgres"`
	BigCache      BigCacheConfig      `yaml:"big_cache"`
	Redis         RedisConfig         `yaml:"redis"`
	FileServer    FileServerConfig    `yaml:"file_server"`
	Cron          CronConfig          `yaml:"cron"`
	Middleware    MiddlewareConfig    `yaml:"middleware"`
	Cookie        CookieConfig        `yaml:"cookie"`
	IPs           IPsConfig           `yaml:"ips"`
	SwaggerServer SwaggerServerConfig `yaml:"swagger_server"`
	HTTPServer    HTTPServerConfig    `yaml:"httpserver"`
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
	f, err := os.Open(configFile) // #nosec G304
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
