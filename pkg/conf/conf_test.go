package conf

import (
	"os"
	"sync"
	"testing"
)

func TestGet(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("PORT", "8080")
	os.Setenv("DB_URL", "test_db_url")

	config := Get()

	if config.Port != "8080" {
		t.Errorf("Expected Port to be '8080', but got '%s'", config.Port)
	}

	if config.DBURL != "test_db_url" {
		t.Errorf("Expected DBURL to be 'test_db_url', but got '%s'", config.DBURL)
	}
}

func TestGet_Singleton(t *testing.T) {
	// Set initial environment variables
	os.Setenv("PORT", "8080")
	os.Setenv("DB_URL", "test_db_url")

	// Get the config for the first time
	config1 := Get()

	// Change environment variables
	os.Setenv("PORT", "9090")
	os.Setenv("DB_URL", "changed_db_url")

	// Get the config for the second time
	config2 := Get()

	// The config should not change, as it should only be loaded once
	if config1.Port != config2.Port || config1.DBURL != config2.DBURL {
		t.Errorf("Config should be loaded only once. Got different configs on multiple calls.")
	}
}

func TestGet_DefaultValues(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("PORT")
	os.Unsetenv("DB_URL")

	// Reset the once.Do function to allow reloading the config
	once = sync.Once{}

	config := Get()

	if config.Port != "" {
		t.Errorf("Expected default Port to be empty, but got '%s'", config.Port)
	}

	if config.DBURL != "" {
		t.Errorf("Expected default DBURL to be empty, but got '%s'", config.DBURL)
	}
}
