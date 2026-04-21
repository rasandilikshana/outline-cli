package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	URL    string `mapstructure:"url" yaml:"url"`
	APIKey string `mapstructure:"api_key" yaml:"api_key"`

	// Source is populated by Load(). It reports where credentials came from:
	// the path of the file used, "env" if pulled from env vars only, or "" if nothing was found.
	Source string `mapstructure:"-" yaml:"-"`
}

func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".outline-cli")
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.yaml")
}

// ProjectConfigPath returns the preferred project-local config location.
// Loader also accepts ./outline.yaml at the project root for teams that
// prefer to keep the file visible.
func ProjectConfigPath() string {
	return filepath.Join(".outline", "config.yaml")
}

// Load reads configuration with this precedence (highest wins):
//  1. OUTLINE_URL / OUTLINE_API_KEY environment variables
//  2. ./.outline/config.yaml   (project-scoped, preferred for teams)
//  3. ./outline.yaml           (project-scoped alternative)
//  4. ~/.outline-cli/config.yaml (home, the original location)
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	v.SetEnvPrefix("OUTLINE")
	_ = v.BindEnv("url", "OUTLINE_URL")
	_ = v.BindEnv("api_key", "OUTLINE_API_KEY")
	v.AutomaticEnv()

	// Search project-scoped then home. Viper picks the first that exists.
	v.SetConfigName("config")
	v.AddConfigPath(".outline")
	v.AddConfigPath(ConfigDir())
	_ = v.ReadInConfig()

	// Fallback: ./outline.yaml (a common "visible" convention).
	if v.ConfigFileUsed() == "" {
		if _, err := os.Stat("outline.yaml"); err == nil {
			v.SetConfigFile("outline.yaml")
			_ = v.ReadInConfig()
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	switch {
	case v.ConfigFileUsed() != "":
		cfg.Source = v.ConfigFileUsed()
	case cfg.URL != "" || cfg.APIKey != "":
		cfg.Source = "env"
	}
	return &cfg, nil
}

// Save writes the given config to the home-dir location. Project-local files
// are expected to be managed by the team (checked into git or templated) and
// aren't overwritten here.
func Save(cfg *Config) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	v := viper.New()
	v.Set("url", cfg.URL)
	v.Set("api_key", cfg.APIKey)
	return v.WriteConfigAs(ConfigPath())
}

func (c *Config) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("URL not configured. Run: outline config --url=<your-outline-url> (or set OUTLINE_URL)")
	}
	if c.APIKey == "" {
		return fmt.Errorf("API key not configured. Run: outline config --api-key=<your-api-key> (or set OUTLINE_API_KEY)")
	}
	return nil
}
