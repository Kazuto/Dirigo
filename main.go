package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Paths map[string]string `yaml:"paths"`
}

func getConfigPath() (string, error) {
	xdg := os.Getenv("XDG_CONFIG_HOME")
	home := os.Getenv("HOME")

	if xdg == "" && home == "" {
		return "", fmt.Errorf("neither XDG_CONFIG_HOME nor HOME are set")
	}

	base := filepath.Join(home, ".config")

	if xdg != "" {
		base = xdg
	}

	return filepath.Join(base, "dirigo", "paths.yml"), nil
}

func ensureConfig() (*Config, string, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, "", err
	}

	configDir := filepath.Dir(configPath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, "", err
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		initial := Config{
			Paths: map[string]string{},
		}

		data, err := yaml.Marshal(initial)
		if err != nil {
			return nil, "", err
		}

		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return nil, "", err
		}

		return &initial, configPath, nil
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, "", err
	}

	var cfg Config
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return nil, "", err
	}

	return &cfg, configPath, nil
}

func addPath(cfg *Config, args []string) {
	key := args[1]
	path := args[2]

	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid path: %v\n", err)
		os.Exit(1)
	}

	cfg.Paths[key] = absPath

	configPath, err := getConfigPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get config path: %v\n", err)
		os.Exit(1)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to serialize config: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Added: %s → %s\n", key, absPath)
}

func deletePath(cfg *Config, args []string) {
	key := args[1]

	if _, exists := cfg.Paths[key]; !exists {
		fmt.Fprintf(os.Stderr, "Key not found: %s\n", key)
		os.Exit(1)
	}

	delete(cfg.Paths, key)

	configPath, err := getConfigPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get config path: %v\n", err)
		os.Exit(1)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to serialize config: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Deleted key: %s\n", key)
}

func listPaths(cfg *Config) {
	if len(cfg.Paths) == 0 {
		fmt.Println("No paths defined.")
		return
	}

	keys := make([]string, 0, len(cfg.Paths))
	for k := range cfg.Paths {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%-10s → %s\n", k, cfg.Paths[k])
	}
}

func main() {
	cfg, _, err := ensureConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	args := os.Args[1:]

	if len(args) == 3 && args[0] == "--add" {
		addPath(cfg, args)

		return
	}

	if len(args) == 2 && args[0] == "--remove" {
		deletePath(cfg, args)

		return
	}

	if len(args) == 0 || (len(args) == 1 && args[0] == "--list") {
		listPaths(cfg)

		return
	}

	key := args[0]
	path, ok := cfg.Paths[key]

	if !ok {
		fmt.Fprintf(os.Stderr, "Key not found in paths: %s\n", key)
		os.Exit(1)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Path does not exist: %s\n", path)
		os.Exit(1)
	}

	fmt.Println(path)
}
