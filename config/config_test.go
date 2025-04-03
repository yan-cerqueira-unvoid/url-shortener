package config

import (
	"os"
	"testing"
)

func TestLoadConfigDefaults(t *testing.T) {
	err := os.Unsetenv("PORT")
	if err != nil {
		t.Errorf("Error unsetting PORT environment variable: %v", err)
	}

	err = os.Unsetenv("SERVER_READ_TIMEOUT")
	if err != nil {
		t.Errorf("Error unsetting SERVER_READ_TIMEOUT environment variable: %v", err)
	}

	err = os.Unsetenv("SERVER_WRITE_TIMEOUT")
	if err != nil {
		t.Errorf("Error unsetting SERVER_WRITE_TIMEOUT environment variable: %v", err)
	}

	err = os.Unsetenv("SERVER_IDLE_TIMEOUT")
	if err != nil {
		t.Errorf("Error unsetting SERVER_IDLE_TIMEOUT environment variable: %v", err)
	}

	err = os.Unsetenv("MONGO_URI")
	if err != nil {
		t.Errorf("Error unsetting MONGO_URI environment variable: %v", err)
	}

	err = os.Unsetenv("DB_NAME")
	if err != nil {
		t.Errorf("Error unsetting DB_NAME environment variable: %v", err)
	}

	err = os.Unsetenv("MONGO_TIMEOUT")
	if err != nil {
		t.Errorf("Error unsetting MONGO_TIMEOUT environment variable: %v", err)
	}

	err = os.Unsetenv("URL_DEFAULT_EXPIRY_DAYS")
	if err != nil {
		t.Errorf("Error unsetting URL_DEFAULT_EXPIRY_DAYS environment variable: %v", err)
	}

	err = os.Unsetenv("URL_CODE_LENGTH")
	if err != nil {
		t.Errorf("Error unsetting URL_CODE_LENGTH environment variable: %v", err)
	}
}

func TestLoadConfigCustomValues(t *testing.T) {
	err := os.Setenv("PORT", "9090")
	if err != nil {
		t.Errorf("Error setting PORT environment variable: %v", err)
	}

	err = os.Setenv("SERVER_READ_TIMEOUT", "15")
	if err != nil {
		t.Errorf("Error setting SERVER_READ_TIMEOUT environment variable: %v", err)
	}

	err = os.Setenv("SERVER_WRITE_TIMEOUT", "20")
	if err != nil {
		t.Errorf("Error setting SERVER_WRITE_TIMEOUT environment variable: %v", err)
	}

	err = os.Setenv("SERVER_IDLE_TIMEOUT", "180")
	if err != nil {
		t.Errorf("Error setting SERVER_IDLE_TIMEOUT environment variable: %v", err)
	}

	err = os.Setenv("MONGO_URI", "mongodb://testhost:27017")
	if err != nil {
		t.Errorf("Error setting MONGO_URI environment variable: %v", err)
	}

	err = os.Setenv("DB_NAME", "test_db")
	if err != nil {
		t.Errorf("Error setting DB_NAME environment variable: %v", err)
	}

	err = os.Setenv("MONGO_TIMEOUT", "30")
	if err != nil {
		t.Errorf("Error setting MONGO_TIMEOUT environment variable: %v", err)
	}

	err = os.Setenv("URL_DEFAULT_EXPIRY_DAYS", "180")
	if err != nil {
		t.Errorf("Error setting URL_DEFAULT_EXPIRY_DAYS environment variable: %v", err)
	}

	err = os.Setenv("URL_CODE_LENGTH", "8")
	if err != nil {
		t.Errorf("Error setting URL_CODE_LENGTH environment variable: %v", err)
	}

	err = os.Unsetenv("PORT")
	if err != nil {
		t.Errorf("Error unsetting PORT environment variable: %v", err)
	}

	err = os.Unsetenv("SERVER_READ_TIMEOUT")
	if err != nil {
		t.Errorf("Error unsetting SERVER_READ_TIMEOUT environment variable: %v", err)
	}

	err = os.Unsetenv("SERVER_WRITE_TIMEOUT")
	if err != nil {
		t.Errorf("Error unsetting SERVER_WRITE_TIMEOUT environment variable: %v", err)
	}

	err = os.Unsetenv("SERVER_IDLE_TIMEOUT")
	if err != nil {
		t.Errorf("Error unsetting SERVER_IDLE_TIMEOUT environment variable: %v", err)
	}

	err = os.Unsetenv("MONGO_URI")
	if err != nil {
		t.Errorf("Error unsetting MONGO_URI environment variable: %v", err)
	}

	err = os.Unsetenv("DB_NAME")
	if err != nil {
		t.Errorf("Error unsetting DB_NAME environment variable: %v", err)
	}

	err = os.Unsetenv("MONGO_TIMEOUT")
	if err != nil {
		t.Errorf("Error unsetting MONGO_TIMEOUT environment variable: %v", err)
	}

	err = os.Unsetenv("URL_DEFAULT_EXPIRY_DAYS")
	if err != nil {
		t.Errorf("Error unsetting URL_DEFAULT_EXPIRY_DAYS environment variable: %v", err)
	}

	err = os.Unsetenv("URL_CODE_LENGTH")
	if err != nil {
		t.Errorf("Error unsetting URL_CODE_LENGTH environment variable: %v", err)
	}
}

func TestGetEnv(t *testing.T) {
	err := os.Setenv("TEST_KEY", "test_value")
	if err != nil {
		t.Errorf("Error setting TEST_KEY environment variable: %v", err)
	}

	value := getEnv("TEST_KEY", "default_value")
	if value != "test_value" {
		t.Errorf("Expected value to be test_value, got %s", value)
	}

	err = os.Unsetenv("TEST_KEY")
	if err != nil {
		t.Errorf("Error unsetting TEST_KEY environment variable: %v", err)
	}

	value = getEnv("TEST_KEY", "default_value")
	if value != "default_value" {
		t.Errorf("Expected value to be default_value, got %s", value)
	}
}
