package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func prefsPath() (string, error) {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return "", fmt.Errorf("LOCALAPPDATA not set")
	}
	return filepath.Join(localAppData, "Google", "Chrome", "User Data", "Local State"), nil
}

func readPrefs() (map[string]any, error) {
	path, err := prefsPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read Preferences: %w", err)
	}
	var prefs map[string]any
	if err := json.Unmarshal(data, &prefs); err != nil {
		return nil, fmt.Errorf("parse Preferences: %w", err)
	}
	return prefs, nil
}

func ReadAccelerationState() (bool, error) {
	prefs, err := readPrefs()
	if err != nil {
		return false, err
	}
	section, ok := prefs["hardware_acceleration_mode"].(map[string]any)
	if !ok {
		// Key absent means Chrome defaults to enabled.
		return true, nil
	}
	enabled, _ := section["enabled"].(bool)
	return enabled, nil
}

func WriteAccelerationState(enabled bool) error {
	path, err := prefsPath()
	if err != nil {
		return err
	}
	prefs, err := readPrefs()
	if err != nil {
		return err
	}

	section, ok := prefs["hardware_acceleration_mode"].(map[string]any)
	if !ok {
		section = make(map[string]any)
	}
	section["enabled"] = enabled
	prefs["hardware_acceleration_mode"] = section
	prefs["hardware_acceleration_mode_previous"] = enabled

	data, err := json.Marshal(prefs)
	if err != nil {
		return fmt.Errorf("marshal Preferences: %w", err)
	}

	// Write atomically: temp file in same directory, then rename.
	tmp := path + ".gputoggle.tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}
