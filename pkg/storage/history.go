package storage

import (
	"encoding/json"
	"os"
	"time"
)

type InviteEntry struct {
	Timestamp string `json:"timestamp"`
	Keyword   string `json:"keyword"`
	Action    string `json:"action"`
	Name      string `json:"name"`
}

// LogAction appends a new action to the history.json file
func LogAction(keyword, action string, name string) error {
	filename := "history.json"
	var history []InviteEntry

	// 1. Read existing file
	data, err := os.ReadFile(filename)
	if err == nil {
		json.Unmarshal(data, &history)
	}

	// 2. Add new entry
	newEntry := InviteEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Keyword:   keyword,
		Action:    action,
		Name:      name,
	}
	history = append(history, newEntry)

	// 3. Write back to file
	updatedData, _ := json.MarshalIndent(history, "", "  ")
	return os.WriteFile(filename, updatedData, 0644)
}
