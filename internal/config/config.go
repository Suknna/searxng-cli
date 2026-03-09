package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	DefaultServer   = "https://searxng.searxng.orb.local"
	DefaultTimeout  = 10 * time.Second
	DefaultLimit    = 10
	DefaultTemplate = "Title={{.Title}} URL={{.URL}} Content={{.Content}}"
)

type ContextConfig struct {
	Server   string `yaml:"server"`
	Timeout  string `yaml:"timeout"`
	Limit    int    `yaml:"limit"`
	Template string `yaml:"template"`
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
}

type Effective struct {
	ContextName string
	Server      string
	Timeout     time.Duration
	Limit       int
	Template    string
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
	}

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

	if loaded.Contexts == nil {
		loaded.Contexts = map[string]ContextConfig{}
	}
	loaded.CurrentContext = eff.ContextName
	loaded.Contexts[eff.ContextName] = ContextConfig{
		Server:   eff.Server,
		Timeout:  eff.Timeout.String(),
		Limit:    eff.Limit,
		Template: eff.Template,
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
