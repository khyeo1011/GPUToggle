package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func findChromeExe() (string, error) {
	candidates := []string{
		filepath.Join(os.Getenv("PROGRAMFILES"), "Google", "Chrome", "Application", "chrome.exe"),
		filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "Google", "Chrome", "Application", "chrome.exe"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Google", "Chrome", "Application", "chrome.exe"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("chrome.exe not found in standard locations")
}

// RestartChrome writes the new acceleration state, force-kills Chrome (preventing
// any graceful-shutdown pref write from overwriting us), then relaunches.
func RestartChrome(enabled bool) error {
	if err := WriteAccelerationState(enabled); err != nil {
		return err
	}

	chromePath, err := findChromeExe()
	if err != nil {
		chromePath = ""
	}

	// /f skips Chrome's graceful shutdown so it cannot rewrite Local State.
	exec.Command("taskkill", "/f", "/im", "chrome.exe").Run()
	time.Sleep(500 * time.Millisecond)

	if chromePath == "" {
		return exec.Command("cmd", "/c", "start", "", "chrome").Start()
	}
	return exec.Command(chromePath).Start()
}
