package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Name         string        `mapstructure:"name"`
	Mode         string        `mapstructure:"mode"`
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type DatabaseConfig struct {
	URI            string        `mapstructure:"uri"`
	Database       string        `mapstructure:"database"`
	ConnectTimeout time.Duration `mapstructure:"connect_timeout"`
	MaxPoolSize    uint64        `mapstructure:"max_pool_size"`
}

type LogConfig struct {
	Level  string        `mapstructure:"level"`
	Format string        `mapstructure:"format"`
	Output string        `mapstructure:"output"`
	File   LogFileConfig `mapstructure:"file"`
}

type LogFileConfig struct {
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type JWTConfig struct {
	Secret string        `mapstructure:"secret"`
	Issuer string        `mapstructure:"issuer"`
	Expire time.Duration `mapstructure:"expire"`
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetEnvPrefix("SHUYOU")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}
	if c.Database.URI == "" {
		return fmt.Errorf("database uri is required")
	}
	if c.Database.Database == "" {
		return fmt.Errorf("database name is required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("jwt secret is required")
	}
	if c.JWT.Expire <= 0 {
		return fmt.Errorf("jwt expire must be positive")
	}
	return nil
}
