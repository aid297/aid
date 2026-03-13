package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	configpkg "github.com/aid297/aid/simpleDB/config"
)

func runGenKey(args []string, stdout, stderr *os.File) int {
	fs := newFlagSet("gen-key")
	configPath := fs.String("config", defaultConfigPath(), "config file path")
	algorithm := fs.String("algo", "", "encryption algorithm (override config)")
	length := fs.Int("len", 0, "key length in bytes (default 32 for AES-256)")
	format := fs.String("format", "hex", "output format (hex, base64)")

	if err := fs.Parse(args); err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return 2
	}

	// Try to load config to determine default algorithm if not provided via CLI
	targetAlgo := strings.ToLower(strings.TrimSpace(*algorithm))
	if targetAlgo == "" {
		if cfg, err := configpkg.Load(*configPath); err == nil {
			targetAlgo = strings.ToLower(strings.TrimSpace(cfg.Engine.Security.EncryptAlgorithm))
		}
	}

	if targetAlgo == "" {
		_, _ = fmt.Fprintln(stderr, "encryption algorithm not specified (use -algo or set in config file)")
		return 1
	}

	var keyLen int
	if *length > 0 {
		keyLen = *length
	} else {
		switch targetAlgo {
		case "aes":
			keyLen = 32 // Default AES-256
		default:
			_, _ = fmt.Fprintf(stderr, "unknown algorithm '%s', please specify key length with -len\n", targetAlgo)
			return 1
		}
	}

	if keyLen <= 0 {
		_, _ = fmt.Fprintln(stderr, "key length must be positive")
		return 1
	}

	key := make([]byte, keyLen)
	if _, err := rand.Read(key); err != nil {
		_, _ = fmt.Fprintf(stderr, "failed to generate random key: %v\n", err)
		return 1
	}

	var output string
	switch strings.ToLower(*format) {
	case "hex":
		output = hex.EncodeToString(key)
	case "base64":
		output = base64.StdEncoding.EncodeToString(key)
	default:
		_, _ = fmt.Fprintf(stderr, "unsupported output format: %s\n", *format)
		return 1
	}

	_, _ = fmt.Fprintln(stdout, output)
	return 0
}
