package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	DefaultServer   = "https://searx.example.org/"
	DefaultTimeout  = 10 * time.Second
	DefaultLimit    = 10
	DefaultTemplate = "Title={{.Title}} URL={{.URL}} Content={{.Content}}"
)

type ContextConfig struct {
	Server   string `yaml:"server"`
	Timeout  string `yaml:"timeout"`
	Limit    int    `yaml:"limit"`
	Template string `yaml:"template"`
	Auth     Auth   `yaml:"auth"`
}

type Auth struct {
	Mode      string `yaml:"mode"`
	APIHeader string `yaml:"api_key_header"`
	APIKey    string `yaml:"api_key"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

type Config struct {
	APIVersion     string                   `yaml:"apiVersion"`
	Kind           string                   `yaml:"kind"`
	CurrentContext string                   `yaml:"current-context"`
	Contexts       map[string]ContextConfig `yaml:"contexts"`
}

type Overrides struct {
	ConfigPath string
	Context    string
	Server     string
	Timeout    time.Duration
	Limit      *int
	Template   *string
	AuthMode   *string
	AuthHeader *string
	AuthAPIKey *string
	AuthUser   *string
	AuthPass   *string
}

type Effective struct {
	ContextName string
	Server      string
	Timeout     time.Duration
	Limit       int
	Template    string
	AuthMode    string
	AuthHeader  string
	AuthAPIKey  string
	AuthUser    string
	AuthPass    string
}

func DefaultConfig() Config {
	return Config{
		APIVersion:     "searxng-cli/v1",
		Kind:           "Config",
		CurrentContext: "default",
		Contexts: map[string]ContextConfig{
			"default": {
				Server:   DefaultServer,
				Timeout:  DefaultTimeout.String(),
				Limit:    DefaultLimit,
				Template: DefaultTemplate,
				Auth: Auth{
					Mode:      "none",
					APIHeader: "X-API-Key",
				},
			},
		},
	}
}

func DefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "searxng-cli", "config.yml"), nil
}

func LoadEffective(o Overrides) (Effective, Config, error) {
	base := DefaultConfig()
	loaded := base

	path := o.ConfigPath
	if path == "" {
		var err error
		path, err = DefaultConfigPath()
		if err != nil {
			return Effective{}, Config{}, err
		}
	}

	if b, err := os.ReadFile(path); err == nil {
		if err := yaml.Unmarshal(b, &loaded); err != nil {
			return Effective{}, Config{}, fmt.Errorf("parse config: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return Effective{}, Config{}, err
	}

	ctxName := loaded.CurrentContext
	if o.Context != "" {
		ctxName = o.Context
	}
	if ctxName == "" {
		ctxName = "default"
	}

	ctx, ok := loaded.Contexts[ctxName]
	if !ok {
		return Effective{}, Config{}, fmt.Errorf("context %q not found", ctxName)
	}

	eff := Effective{
		ContextName: ctxName,
		Server:      fallbackStr(ctx.Server, DefaultServer),
		Timeout:     parseTimeoutOrDefault(ctx.Timeout),
		Limit:       fallbackInt(ctx.Limit, DefaultLimit),
		Template:    fallbackStr(ctx.Template, DefaultTemplate),
		AuthMode:    fallbackStr(ctx.Auth.Mode, "none"),
		AuthHeader:  fallbackStr(ctx.Auth.APIHeader, "X-API-Key"),
		AuthAPIKey:  ctx.Auth.APIKey,
		AuthUser:    ctx.Auth.Username,
		AuthPass:    ctx.Auth.Password,
	}

	applyAuthFromEnv(&eff)

	if o.Server != "" {
		eff.Server = o.Server
	}
	if o.Timeout > 0 {
		eff.Timeout = o.Timeout
	}
	if o.Limit != nil {
		eff.Limit = *o.Limit
	}
	if o.Template != nil {
		eff.Template = *o.Template
	}
	if o.AuthMode != nil {
		eff.AuthMode = *o.AuthMode
	}
	if o.AuthHeader != nil {
		eff.AuthHeader = *o.AuthHeader
	}
	if o.AuthAPIKey != nil {
		eff.AuthAPIKey = *o.AuthAPIKey
	}
	if o.AuthUser != nil {
		eff.AuthUser = *o.AuthUser
	}
	if o.AuthPass != nil {
		eff.AuthPass = *o.AuthPass
	}

	if loaded.Contexts == nil {
		loaded.Contexts = map[string]ContextConfig{}
	}
	loaded.CurrentContext = eff.ContextName
	loaded.Contexts[eff.ContextName] = ContextConfig{
		Server:   eff.Server,
		Timeout:  eff.Timeout.String(),
		Limit:    eff.Limit,
		Template: eff.Template,
		Auth: Auth{
			Mode:      eff.AuthMode,
			APIHeader: eff.AuthHeader,
			APIKey:    eff.AuthAPIKey,
			Username:  eff.AuthUser,
			Password:  eff.AuthPass,
		},
	}

	return eff, loaded, nil
}

func WriteDefault(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := yaml.Marshal(DefaultConfig())
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func LoadRaw(path string) (Config, error) {
	base := DefaultConfig()
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return base, nil
	}
	if err != nil {
		return Config{}, err
	}
	if err := yaml.Unmarshal(b, &base); err != nil {
		return Config{}, err
	}
	if base.Contexts == nil {
		base.Contexts = map[string]ContextConfig{}
	}
	return base, nil
}

func Save(path string, cfg Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func fallbackStr(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

func fallbackInt(v, def int) int {
	if v <= 0 {
		return def
	}
	return v
}

func parseTimeoutOrDefault(raw string) time.Duration {
	if raw == "" {
		return DefaultTimeout
	}
	v, err := time.ParseDuration(raw)
	if err != nil || v <= 0 {
		return DefaultTimeout
	}
	return v
}

func applyAuthFromEnv(e *Effective) {
	if v := strings.TrimSpace(os.Getenv("SEARXNG_AUTH_MODE")); v != "" {
		e.AuthMode = v
	}
	if v := strings.TrimSpace(os.Getenv("SEARXNG_AUTH_HEADER")); v != "" {
		e.AuthHeader = v
	}
	if v := strings.TrimSpace(os.Getenv("SEARXNG_AUTH_API_KEY")); v != "" {
		e.AuthAPIKey = v
	}
	if v := strings.TrimSpace(os.Getenv("SEARXNG_AUTH_USERNAME")); v != "" {
		e.AuthUser = v
	}
	if v := strings.TrimSpace(os.Getenv("SEARXNG_AUTH_PASSWORD")); v != "" {
		e.AuthPass = v
	}
}
