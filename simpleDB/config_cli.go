package main

import (
	"fmt"
	"os"

	configpkg "github.com/aid297/aid/simpleDB/config"
	"gopkg.in/yaml.v3"
)

func defaultConfigPath() string {
	return configpkg.DefaultConfigPath
}

func runPrintConfig(args []string, stdout, stderr *os.File) int {
	fs := newFlagSet("print-config")
	configPath := fs.String("config", defaultConfigPath(), "config file path")
	if err := fs.Parse(args); err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return 2
	}

	config, err := configpkg.Load(*configPath)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "load config failed: %v\n", err)
		return 1
	}
	config.ApplyDefaults()
	payload, err := yaml.Marshal(config)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "marshal config failed: %v\n", err)
		return 1
	}
	_, _ = stdout.Write(payload)
	return 0
}

func runInitConfig(args []string, stdout, stderr *os.File) int {
	fs := newFlagSet("init-config")
	configPath := fs.String("config", defaultConfigPath(), "config file path")
	force := fs.Bool("force", false, "overwrite existing config")
	if err := fs.Parse(args); err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return 2
	}

	if !*force {
		if _, err := os.Stat(*configPath); err == nil {
			_, _ = fmt.Fprintf(stderr, "config already exists: %s\n", *configPath)
			return 1
		}
	}

	if err := configpkg.Save(*configPath, configpkg.Default()); err != nil {
		_, _ = fmt.Fprintf(stderr, "write config failed: %v\n", err)
		return 1
	}
	_, _ = fmt.Fprintf(stdout, "config written to %s\n", *configPath)
	return 0
}
