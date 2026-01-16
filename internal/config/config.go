package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Storage  StorageConfig  `yaml:"storage"`
	Email    EmailConfig    `yaml:"email"`
	App      AppConfig      `yaml:"app"`
}

type ServerConfig struct {
	Port         string        `yaml:"port"`
	Host         string        `yaml:"host"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type DatabaseConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	DBName          string        `yaml:"dbname"`
	SSLMode         string        `yaml:"sslmode"`
	URL             string        `yaml:"url"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
}

func (d *DatabaseConfig) DSN() string {
	if d.URL != "" {
		return d.URL
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode)
}

type RedisConfig struct {
	Addr         string        `yaml:"addr"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	PoolSize     int           `yaml:"pool_size"`
	MinIdleConns int           `yaml:"min_idle_conns"`
}

type StorageConfig struct {
	Type        string `yaml:"type"`
	BasePath    string `yaml:"base_path"`
	BaseURL     string `yaml:"base_url"`
	S3Bucket    string `yaml:"s3_bucket,omitempty"`
	S3Region    string `yaml:"s3_region,omitempty"`
	S3AccessKey string `yaml:"s3_access_key,omitempty"`
	S3SecretKey string `yaml:"s3_secret_key,omitempty"`
}

type EmailConfig struct {
	Provider string     `yaml:"provider"`
	SMTP     SMTPConfig `yaml:"smtp,omitempty"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

type AppConfig struct {
	Environment    string        `yaml:"environment"`
	LogLevel       string        `yaml:"log_level"`
	IdempotencyTTL time.Duration `yaml:"idempotency_ttl"`
	CacheTTL       time.Duration `yaml:"cache_ttl"`
	LockTTL        time.Duration `yaml:"lock_ttl"`
	JWTSecret      string        `yaml:"jwt_secret"`
	JWTExpiration  time.Duration `yaml:"jwt_expiration"`
}

func Load(configPath string) (*Config, error) {
	cfg := &Config{}

	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if len(data) > 0 {
			if err := yaml.Unmarshal(data, cfg); err != nil {
				return nil, fmt.Errorf("failed to parse config file: %w", err)
			}
		}
	}

	overrideWithEnv(cfg)
	setDefaults(cfg)

	return cfg, nil
}

func overrideWithEnv(cfg *Config) {
	if port := os.Getenv("PORT"); port != "" {
		cfg.Server.Port = port
	}
	if host := os.Getenv("HOST"); host != "" {
		cfg.Server.Host = host
	}

	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		cfg.Database.URL = dsn
	} else {
		if host := os.Getenv("DB_HOST"); host != "" {
			cfg.Database.Host = host
		}
		if user := os.Getenv("DB_USER"); user != "" {
			cfg.Database.User = user
		}
		if password := os.Getenv("DB_PASSWORD"); password != "" {
			cfg.Database.Password = password
		}
		if dbname := os.Getenv("DB_NAME"); dbname != "" {
			cfg.Database.DBName = dbname
		}
	}

	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		cfg.Redis.Addr = addr
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		cfg.Redis.Password = password
	}

	if path := os.Getenv("FILE_STORAGE_PATH"); path != "" {
		cfg.Storage.BasePath = path
	}
	if url := os.Getenv("FILE_STORAGE_URL"); url != "" {
		cfg.Storage.BaseURL = url
	}

	if env := os.Getenv("ENVIRONMENT"); env != "" {
		cfg.App.Environment = env
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.App.LogLevel = logLevel
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		cfg.App.JWTSecret = jwtSecret
	}
}

func setDefaults(cfg *Config) {
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 15 * time.Second
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 15 * time.Second
	}
	if cfg.Server.IdleTimeout == 0 {
		cfg.Server.IdleTimeout = 60 * time.Second
	}

	if cfg.Database.Host == "" {
		cfg.Database.Host = "localhost"
	}
	if cfg.Database.Port == 0 {
		cfg.Database.Port = 5432
	}
	if cfg.Database.User == "" {
		cfg.Database.User = "postgres"
	}
	if cfg.Database.Password == "" {
		cfg.Database.Password = "postgres"
	}
	if cfg.Database.DBName == "" {
		cfg.Database.DBName = "loan_db"
	}
	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disable"
	}
	if cfg.Database.MaxOpenConns == 0 {
		cfg.Database.MaxOpenConns = 25
	}
	if cfg.Database.MaxIdleConns == 0 {
		cfg.Database.MaxIdleConns = 5
	}
	if cfg.Database.ConnMaxLifetime == 0 {
		cfg.Database.ConnMaxLifetime = 5 * time.Minute
	}
	if cfg.Database.ConnMaxIdleTime == 0 {
		cfg.Database.ConnMaxIdleTime = 10 * time.Minute
	}

	if cfg.Redis.Addr == "" {
		cfg.Redis.Addr = "localhost:6379"
	}
	if cfg.Redis.PoolSize == 0 {
		cfg.Redis.PoolSize = 10
	}
	if cfg.Redis.MinIdleConns == 0 {
		cfg.Redis.MinIdleConns = 5
	}
	if cfg.Redis.DialTimeout == 0 {
		cfg.Redis.DialTimeout = 5 * time.Second
	}

	if cfg.Storage.Type == "" {
		cfg.Storage.Type = "local"
	}
	if cfg.Storage.BasePath == "" {
		cfg.Storage.BasePath = "./storage"
	}
	if cfg.Storage.BaseURL == "" {
		cfg.Storage.BaseURL = "http://localhost:8080/files"
	}

	if cfg.Email.Provider == "" {
		cfg.Email.Provider = "mock"
	}

	if cfg.App.Environment == "" {
		cfg.App.Environment = "development"
	}
	if cfg.App.LogLevel == "" {
		cfg.App.LogLevel = "info"
	}
	if cfg.App.IdempotencyTTL == 0 {
		cfg.App.IdempotencyTTL = 24 * time.Hour
	}
	if cfg.App.CacheTTL == 0 {
		cfg.App.CacheTTL = 5 * time.Minute
	}
	if cfg.App.LockTTL == 0 {
		cfg.App.LockTTL = 30 * time.Second
	}
	if cfg.App.JWTSecret == "" {
		cfg.App.JWTSecret = "your-secret-key-change-in-production"
	}
	if cfg.App.JWTExpiration == 0 {
		cfg.App.JWTExpiration = 24 * time.Hour
	}
}
