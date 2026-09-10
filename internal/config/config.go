package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	ConfigDirName  = ".forgectl"
	ConfigFileName = "config.yaml"
)

type Config struct {
	v *viper.Viper
}

type RemoteConfig struct {
	Provider string `mapstructure:"provider"`
	Token    string `mapstructure:"token"`
	Username string `mapstructure:"username"`
	BaseURL  string `mapstructure:"base_url"`
	Default  bool   `mapstructure:"default"`
}

type ProjectConfig struct {
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
	License     string `mapstructure:"license"`
	Template    string `mapstructure:"template"`
	Provider    string `mapstructure:"provider"`
	Remote      string `mapstructure:"remote"`
	LocalPath   string `mapstructure:"local_path"`
}

type AppConfig struct {
	DefaultProvider string                      `mapstructure:"default_provider"`
	DefaultLicense  string                      `mapstructure:"default_license"`
	DefaultTemplate string                      `mapstructure:"default_template"`
	Remotes         map[string]RemoteConfig     `mapstructure:"remotes"`
	Projects        map[string]ProjectConfig    `mapstructure:"projects"`
}

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ConfigDirName), nil
}

func ConfigDir() (string, error) {
	return configDir()
}

func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ConfigFileName), nil
}

func New() (*Config, error) {
	dir, err := configDir()
	if err != nil {
		return nil, fmt.Errorf("cannot get config dir: %w", err)
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("cannot create config dir: %w", err)
	}

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(dir)

	v.SetDefault("default_provider", "github")
	v.SetDefault("default_license", "MIT")
	v.SetDefault("default_template", "go")
	v.SetDefault("remotes", map[string]RemoteConfig{})
	v.SetDefault("projects", map[string]ProjectConfig{})

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("cannot read config: %w", err)
		}
		if err := v.WriteConfigAs(filepath.Join(dir, ConfigFileName)); err != nil {
			return nil, fmt.Errorf("cannot write config: %w", err)
		}
	}

	return &Config{v: v}, nil
}

func (c *Config) GetAppConfig() (*AppConfig, error) {
	var cfg AppConfig
	if err := c.v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config: %w", err)
	}
	return &cfg, nil
}

func (c *Config) Set(key string, value interface{}) error {
	c.v.Set(key, value)
	dir, err := configDir()
	if err != nil {
		return err
	}
	return c.v.WriteConfigAs(filepath.Join(dir, ConfigFileName))
}

func (c *Config) Get(key string) interface{} {
	return c.v.Get(key)
}

func (c *Config) Reset() error {
	c.v.Set("default_provider", "github")
	c.v.Set("default_license", "MIT")
	c.v.Set("default_template", "go")
	c.v.Set("remotes", map[string]RemoteConfig{})
	c.v.Set("projects", map[string]ProjectConfig{})

	dir, err := configDir()
	if err != nil {
		return err
	}
	return c.v.WriteConfigAs(filepath.Join(dir, ConfigFileName))
}

func (c *Config) AddRemote(name string, remote RemoteConfig) error {
	c.v.Set(fmt.Sprintf("remotes.%s", name), remote)
	dir, err := configDir()
	if err != nil {
		return err
	}
	return c.v.WriteConfigAs(filepath.Join(dir, ConfigFileName))
}

func (c *Config) GetRemote(name string) (*RemoteConfig, error) {
	key := fmt.Sprintf("remotes.%s", name)
	if !c.v.IsSet(key) {
		return nil, fmt.Errorf("remote %q not found", name)
	}
	var remote RemoteConfig
	if err := c.v.UnmarshalKey(key, &remote); err != nil {
		return nil, err
	}
	return &remote, nil
}

func (c *Config) GetDefaultRemote() (*RemoteConfig, error) {
	cfg, err := c.GetAppConfig()
	if err != nil {
		return nil, err
	}
	for name, r := range cfg.Remotes {
		if r.Default || name == cfg.DefaultProvider {
			return &RemoteConfig{
				Provider: r.Provider,
				Token:    r.Token,
				Username: r.Username,
				BaseURL:  r.BaseURL,
				Default:  true,
			}, nil
		}
	}
	for name, r := range cfg.Remotes {
		_ = name
		return &RemoteConfig{
			Provider: r.Provider,
			Token:    r.Token,
			Username: r.Username,
			BaseURL:  r.BaseURL,
		}, nil
	}
	return nil, fmt.Errorf("no remotes configured. Run: forgectl config set provider <name>")
}

func (c *Config) AddProject(name string, project ProjectConfig) error {
	c.v.Set(fmt.Sprintf("projects.%s", name), project)
	dir, err := configDir()
	if err != nil {
		return err
	}
	return c.v.WriteConfigAs(filepath.Join(dir, ConfigFileName))
}

func (c *Config) GetProject(name string) (*ProjectConfig, error) {
	key := fmt.Sprintf("projects.%s", name)
	if !c.v.IsSet(key) {
		return nil, fmt.Errorf("project %q not found", name)
	}
	var project ProjectConfig
	if err := c.v.UnmarshalKey(key, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

func (c *Config) ListProjects() map[string]ProjectConfig {
	cfg, err := c.GetAppConfig()
	if err != nil {
		return nil
	}
	return cfg.Projects
}

func (c *Config) RemoveProject(name string) error {
	cfg, err := c.GetAppConfig()
	if err != nil {
		return err
	}
	delete(cfg.Projects, name)
	c.v.Set("projects", cfg.Projects)
	dir, err := configDir()
	if err != nil {
		return err
	}
	return c.v.WriteConfigAs(filepath.Join(dir, ConfigFileName))
}
