package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfigDefaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("SERVER_READ_TIMEOUT")
	os.Unsetenv("SERVER_WRITE_TIMEOUT")
	os.Unsetenv("SERVER_IDLE_TIMEOUT")
	os.Unsetenv("MONGO_URI")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("MONGO_TIMEOUT")
	os.Unsetenv("URL_DEFAULT_EXPIRY_DAYS")
	os.Unsetenv("URL_CODE_LENGTH")

	config := LoadConfig()

	if config.Server.Port != "8080" {
		t.Errorf("Expected default port to be 8080, got %s", config.Server.Port)
	}

	if config.Server.ReadTimeout != 5*time.Second {
		t.Errorf("Expected default read timeout to be 5s, got %v", config.Server.ReadTimeout)
	}

	if config.Server.WriteTimeout != 10*time.Second {
		t.Errorf("Expected default write timeout to be 10s, got %v", config.Server.WriteTimeout)
	}

	if config.Server.IdleTimeout != 120*time.Second {
		t.Errorf("Expected default idle timeout to be 120s, got %v", config.Server.IdleTimeout)
	}

	if config.MongoDB.URI != "mongodb://localhost:27017" {
		t.Errorf("Expected default MongoDB URI to be mongodb://localhost:27017, got %s", config.MongoDB.URI)
	}

	if config.MongoDB.Database != "url_shortener" {
		t.Errorf("Expected default database name to be url_shortener, got %s", config.MongoDB.Database)
	}

	if config.MongoDB.Timeout != 10*time.Second {
		t.Errorf("Expected default MongoDB timeout to be 10s, got %v", config.MongoDB.Timeout)
	}

	if config.URLShortener.DefaultExpiry != 365*24*time.Hour {
		t.Errorf("Expected default URL expiry to be 365 days, got %v", config.URLShortener.DefaultExpiry)
	}

	if config.URLShortener.CodeLength != 6 {
		t.Errorf("Expected default code length to be 6, got %d", config.URLShortener.CodeLength)
	}
}

func TestLoadConfigCustomValues(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("SERVER_READ_TIMEOUT", "15")
	os.Setenv("SERVER_WRITE_TIMEOUT", "20")
	os.Setenv("SERVER_IDLE_TIMEOUT", "180")
	os.Setenv("MONGO_URI", "mongodb://testhost:27017")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("MONGO_TIMEOUT", "30")
	os.Setenv("URL_DEFAULT_EXPIRY_DAYS", "180")
	os.Setenv("URL_CODE_LENGTH", "8")

	config := LoadConfig()

	if config.Server.Port != "9090" {
		t.Errorf("Expected custom port to be 9090, got %s", config.Server.Port)
	}

	if config.Server.ReadTimeout != 15*time.Second {
		t.Errorf("Expected custom read timeout to be 15s, got %v", config.Server.ReadTimeout)
	}

	if config.Server.WriteTimeout != 20*time.Second {
		t.Errorf("Expected custom write timeout to be 20s, got %v", config.Server.WriteTimeout)
	}

	if config.Server.IdleTimeout != 180*time.Second {
		t.Errorf("Expected custom idle timeout to be 180s, got %v", config.Server.IdleTimeout)
	}

	if config.MongoDB.URI != "mongodb://testhost:27017" {
		t.Errorf("Expected custom MongoDB URI to be mongodb://testhost:27017, got %s", config.MongoDB.URI)
	}

	if config.MongoDB.Database != "test_db" {
		t.Errorf("Expected custom database name to be test_db, got %s", config.MongoDB.Database)
	}

	if config.MongoDB.Timeout != 30*time.Second {
		t.Errorf("Expected custom MongoDB timeout to be 30s, got %v", config.MongoDB.Timeout)
	}

	if config.URLShortener.DefaultExpiry != 180*24*time.Hour {
		t.Errorf("Expected custom URL expiry to be 180 days, got %v", config.URLShortener.DefaultExpiry)
	}

	if config.URLShortener.CodeLength != 8 {
		t.Errorf("Expected custom code length to be 8, got %d", config.URLShortener.CodeLength)
	}

	os.Unsetenv("PORT")
	os.Unsetenv("SERVER_READ_TIMEOUT")
	os.Unsetenv("SERVER_WRITE_TIMEOUT")
	os.Unsetenv("SERVER_IDLE_TIMEOUT")
	os.Unsetenv("MONGO_URI")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("MONGO_TIMEOUT")
	os.Unsetenv("URL_DEFAULT_EXPIRY_DAYS")
	os.Unsetenv("URL_CODE_LENGTH")
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	value := getEnv("TEST_KEY", "default_value")
	if value != "test_value" {
		t.Errorf("Expected value to be test_value, got %s", value)
	}

	os.Unsetenv("TEST_KEY")
	value = getEnv("TEST_KEY", "default_value")
	if value != "default_value" {
		t.Errorf("Expected value to be default_value, got %s", value)
	}
}
