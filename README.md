# GPU Toggle

A Chrome extension that lets you enable or disable hardware acceleration with one click, without digging through `chrome://settings/system` every time.

## The problem

Chrome's hardware acceleration toggle is buried four clicks deep in Settings and requires a full browser relaunch to take effect. Users who frequently need to flip it — those with rendering artifacts on certain sites, mixed-DPI setups, older integrated GPUs, or developers testing rendering paths — have no better option today.

## How it works

Chrome extensions can't modify the hardware acceleration setting directly. GPU Toggle pairs a minimal Chrome extension with a small native helper binary that reads and writes Chrome's `Local State` file (where Chrome stores the setting) and triggers a relaunch.

```
Chrome Extension (popup UI)
        │  Native Messaging (stdin/stdout JSON)
        ▼
  GPU Toggle Helper (Go binary)
        │
        ▼
  Chrome Local State file
```

## Current status

**Phase 0 — proof of concept, Windows only.**

The core round-trip works: toggle in the popup → helper writes `Local State` → Chrome relaunches with the new setting applied.

## Project structure

```
extension/          Chrome extension (Manifest V3)
helper/             Native messaging host (Go)
com.gputoggle.helper.json   Native messaging manifest template
install.ps1         Build + registration script (Windows)
```

## Setup (Windows)

**Prerequisites:** Go 1.21+, Chrome

1. Clone the repo.

2. Load the extension in Chrome:
   - Go to `chrome://extensions`
   - Enable **Developer mode**
   - Click **Load unpacked** → select the `extension/` folder
   - Copy the 32-character Extension ID

3. Run the install script in PowerShell:
   ```powershell
   .\install.ps1
   ```
   Paste the Extension ID when prompted. The script builds `helper.exe`, writes the native messaging manifest, and registers it in the Windows registry.

4. Click the GPU Toggle icon in Chrome's toolbar.
