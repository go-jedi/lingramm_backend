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
	URL                   string `yaml:"url"`
	TimeoutCheckConnect   int    `yaml:"timeout_check_connect"`
	TimeoutPublishMessage int    `yaml:"timeout_publish_message"`
	Amqp                  struct {
		Dialer struct {
			Timeout   int `yaml:"timeout"`
			KeepAlive int `yaml:"keep_alive"`
		} `yaml:"dialer"`
		SASL struct {
			IsPlainAuth bool `yaml:"is_plain_auth"`
			PlainAuth   struct {
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			} `yaml:"plain_auth"`
		} `yaml:"sasl"`
		VHost      string `yaml:"v_host"`
		ChannelMax uint16 `yaml:"channel_max"`
		FrameSize  int    `yaml:"frame_size"`
		HeartBeat  int    `yaml:"heart_beat"`
		Properties struct {
			ConnectionName string `yaml:"connection_name"`
		} `yaml:"properties"`
		Locale string `yaml:"locale"`
		IsDial bool   `yaml:"is_dial"`
	} `yaml:"amqp"`
	Exchange struct {
		Name       string `yaml:"name"`
		Kind       string `yaml:"kind"`
		Durable    bool   `yaml:"durable"`
		AutoDelete bool   `yaml:"auto_delete"`
		Internal   bool   `yaml:"internal"`
		NoWait     bool   `yaml:"no_wait"`
	} `yaml:"exchange"`
	Queue struct {
		Name       string `yaml:"name"`
		Durable    bool   `yaml:"durable"`
		AutoDelete bool   `yaml:"auto_delete"`
		Exclusive  bool   `yaml:"exclusive"`
		NoWait     bool   `yaml:"no_wait"`
		Args       struct {
			XExpires    int32 `yaml:"x_expires"`
			XMessageTTL int32 `yaml:"x_message_ttl"`
		} `yaml:"args"`
	} `yaml:"queue"`
	QueueBind struct {
		Exchange string `yaml:"exchange"`
		NoWait   bool   `yaml:"no_wait"`
	} `yaml:"queue_bind"`
	Consume struct {
		Consumer  string `yaml:"consumer"`
		AutoAck   bool   `yaml:"auto_ack"`
		Exclusive bool   `yaml:"exclusive"`
		NoLocal   bool   `yaml:"no_local"`
		NoWait    bool   `yaml:"no_wait"`
	} `yaml:"consume"`
}

type PublisherConfig struct {
	URL                   string `yaml:"url"`
	TimeoutCheckConnect   int    `yaml:"timeout_check_connect"`
	TimeoutPublishMessage int    `yaml:"timeout_publish_message"`
	QueryTimeout          int64  `yaml:"query_timeout"`
	Amqp                  struct {
		Dialer struct {
			Timeout   int `yaml:"timeout"`
			KeepAlive int `yaml:"keep_alive"`
		} `yaml:"dialer"`
		SASL struct {
			IsPlainAuth bool `yaml:"is_plain_auth"`
			PlainAuth   struct {
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			} `yaml:"plain_auth"`
		} `yaml:"sasl"`
		VHost      string `yaml:"v_host"`
		ChannelMax uint16 `yaml:"channel_max"`
		FrameSize  int    `yaml:"frame_size"`
		HeartBeat  int    `yaml:"heart_beat"`
		Properties struct {
			ConnectionName string `yaml:"connection_name"`
		} `yaml:"properties"`
		Locale string `yaml:"locale"`
		IsDial bool   `yaml:"is_dial"`
	} `yaml:"amqp"`
	Exchange struct {
		Name       string `yaml:"name"`
		Kind       string `yaml:"kind"`
		Durable    bool   `yaml:"durable"`
		AutoDelete bool   `yaml:"auto_delete"`
		Internal   bool   `yaml:"internal"`
		NoWait     bool   `yaml:"no_wait"`
	} `yaml:"exchange"`
	Publish struct {
		Exchange   string `yaml:"exchange"`
		Mandatory  bool   `yaml:"mandatory"`
		Immediate  bool   `yaml:"immediate"`
		Publishing struct {
			ContentType string `yaml:"content_type"`
			Timestamp   string `yaml:"timestamp"`
			Type        string `yaml:"type"`
		} `yaml:"publishing"`
	} `yaml:"publish"`
}

type RabbitMQNotificationConfig struct {
	Consumer  ConsumerConfig  `yaml:"consumer"`
	Publisher PublisherConfig `yaml:"publisher"`
}

type RabbitMQConfig struct {
	Notification RabbitMQNotificationConfig `yaml:"notification"`
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
	RabbitMQ      RabbitMQConfig      `yaml:"rabbitmq"`
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
