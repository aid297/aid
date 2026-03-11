package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/aid297/aid/debugLogger"
)

func run(title string, fn func() error) {
	debugLogger.Print("\n===== %s =====\n", title)
	if err := fn(); err != nil {
		debugLogger.Print("ERROR: %v\n", err)
		return
	}
	debugLogger.Print("DONE\n")
}

func printJSON(value any) {
	payload, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		debugLogger.Print("marshal error：%v", err)
		return
	}
	debugLogger.Print("%s", string(payload))
}

func cleanupAll() {
	_ = os.RemoveAll(demoDatabase)
}

func cleanupTable(table string) {
	_ = os.RemoveAll(filepath.Join(demoDatabase, table))
}
