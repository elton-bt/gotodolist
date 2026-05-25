package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"
)

var version = "dev"

type Config struct {
	AppName           string
	Version           string
	Port              string
	DBHost            string
	DBPort            int
	DBName            string
	DBUser            string
	DBPassword        string
	DBSSLMode         string
	CORSAllowOrigin   string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	ShutdownTimeout   time.Duration
}

func Load(appName, defaultPort string) (Config, error) {
	dbPort, err := envInt("DB_PORT", 5432)
	if err != nil {
		return Config{}, fmt.Errorf("parse DB_PORT: %w", err)
	}

	readTimeout, err := envDuration("HTTP_READ_TIMEOUT", 10*time.Second)
	if err != nil {
		return Config{}, fmt.Errorf("parse HTTP_READ_TIMEOUT: %w", err)
	}

	readHeaderTimeout, err := envDuration("HTTP_READ_HEADER_TIMEOUT", 5*time.Second)
	if err != nil {
		return Config{}, fmt.Errorf("parse HTTP_READ_HEADER_TIMEOUT: %w", err)
	}

	writeTimeout, err := envDuration("HTTP_WRITE_TIMEOUT", 15*time.Second)
	if err != nil {
		return Config{}, fmt.Errorf("parse HTTP_WRITE_TIMEOUT: %w", err)
	}

	shutdownTimeout, err := envDuration("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second)
	if err != nil {
		return Config{}, fmt.Errorf("parse HTTP_SHUTDOWN_TIMEOUT: %w", err)
	}

	return Config{
		AppName:           appName,
		Version:           Version(),
		Port:              envString("APP_PORT", defaultPort),
		DBHost:            envString("DB_HOST", "localhost"),
		DBPort:            dbPort,
		DBName:            envString("DB_NAME", "gotodolist"),
		DBUser:            envString("DB_USER", "gotodolist"),
		DBPassword:        envString("DB_PASSWORD", "replace-me"),
		DBSSLMode:         envString("DB_SSLMODE", "disable"),
		CORSAllowOrigin:   envString("CORS_ALLOW_ORIGIN", "*"),
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		ShutdownTimeout:   shutdownTimeout,
	}, nil
}

func Version() string {
	if version == "" {
		return "dev"
	}

	return version
}

func (c Config) Address() string {
	return net.JoinHostPort("", c.Port)
}

func (c Config) DatabaseURL() string {
	return (&url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.DBUser, c.DBPassword),
		Host:   net.JoinHostPort(c.DBHost, strconv.Itoa(c.DBPort)),
		Path:   c.DBName,
		RawQuery: url.Values{
			"sslmode": []string{c.DBSSLMode},
		}.Encode(),
	}).String()
}

func envString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}

	return fallback
}

func envInt(key string, fallback int) (int, error) {
	value := envString(key, strconv.Itoa(fallback))
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return parsed, nil
}

func envDuration(key string, fallback time.Duration) (time.Duration, error) {
	value := envString(key, fallback.String())
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, err
	}

	return parsed, nil
}
