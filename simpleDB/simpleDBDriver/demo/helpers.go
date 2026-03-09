package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func run(title string, fn func() error) {
	fmt.Printf("\n===== %s =====\n", title)
	if err := fn(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	fmt.Println("DONE")
}

func printJSON(value any) {
	payload, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("marshal error:", err)
		return
	}
	fmt.Println(string(payload))
}

func cleanupAll() {
	_ = os.RemoveAll(demoDatabase)
}

func cleanupTable(table string) {
	_ = os.RemoveAll(filepath.Join(demoDatabase, table))
}
