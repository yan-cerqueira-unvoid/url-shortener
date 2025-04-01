package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server       ServerConfig
	MongoDB      MongoDBConfig
	URLShortener URLShortenerConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type MongoDBConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

type URLShortenerConfig struct {
	DefaultExpiry time.Duration
	CodeLength    int
}

func LoadConfig() *Config {
	port := getEnv("PORT", "8080")

	readTimeout, _ := strconv.Atoi(getEnv("SERVER_READ_TIMEOUT", "5"))
	writeTimeout, _ := strconv.Atoi(getEnv("SERVER_WRITE_TIMEOUT", "10"))
	idleTimeout, _ := strconv.Atoi(getEnv("SERVER_IDLE_TIMEOUT", "120"))

	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	mongoDatabase := getEnv("DB_NAME", "url_shortener")
	mongoTimeout, _ := strconv.Atoi(getEnv("MONGO_TIMEOUT", "10"))

	defaultExpiryDays, _ := strconv.Atoi(getEnv("URL_DEFAULT_EXPIRY_DAYS", "365"))
	codeLength, _ := strconv.Atoi(getEnv("URL_CODE_LENGTH", "6"))

	return &Config{
		Server: ServerConfig{
			Port:         port,
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
			IdleTimeout:  time.Duration(idleTimeout) * time.Second,
		},
		MongoDB: MongoDBConfig{
			URI:      mongoURI,
			Database: mongoDatabase,
			Timeout:  time.Duration(mongoTimeout) * time.Second,
		},
		URLShortener: URLShortenerConfig{
			DefaultExpiry: time.Duration(defaultExpiryDays) * 24 * time.Hour,
			CodeLength:    codeLength,
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func (c *Config) PrintConfig() {
	log.Println("Server Configuration:")
	log.Printf("Port: %s\n", c.Server.Port)
	log.Printf("Read Timeout: %v\n", c.Server.ReadTimeout)
	log.Printf("Write Timeout: %v\n", c.Server.WriteTimeout)
	log.Printf("Idle Timeout: %v\n", c.Server.IdleTimeout)

	log.Println("MongoDB Configuration:")
	log.Printf("Database: %s\n", c.MongoDB.Database)
	log.Printf("Timeout: %v\n", c.MongoDB.Timeout)

	log.Println("URL Shortener Configuration:")
	log.Printf("Default Expiry: %v\n", c.URLShortener.DefaultExpiry)
	log.Printf("Code Length: %d\n", c.URLShortener.CodeLength)
}
