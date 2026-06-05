package main

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"os"
)

type inMsg struct {
	Type    string `json:"type"`
	Enabled *bool  `json:"enabled,omitempty"`
}

type outMsg struct {
	Type    string `json:"type"`
	Enabled *bool  `json:"enabled,omitempty"`
	Message string `json:"message,omitempty"`
}

func readMessage(r io.Reader) ([]byte, error) {
	var length uint32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}
	buf := make([]byte, length)
	_, err := io.ReadFull(r, buf)
	return buf, err
}

func writeMessage(w io.Writer, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, uint32(len(data))); err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func boolPtr(b bool) *bool { return &b }

func main() {
	for {
		raw, err := readMessage(os.Stdin)
		if err != nil {
			// Chrome closed the pipe — normal exit.
			return
		}

		var msg inMsg
		if err := json.Unmarshal(raw, &msg); err != nil {
			writeMessage(os.Stdout, outMsg{Type: "error", Message: "invalid JSON"})
			continue
		}

		switch msg.Type {
		case "getState":
			enabled, err := ReadAccelerationState()
			if err != nil {
				writeMessage(os.Stdout, outMsg{Type: "error", Message: err.Error()})
				continue
			}
			writeMessage(os.Stdout, outMsg{Type: "state", Enabled: boolPtr(enabled)})

		case "restart":
			// enabled is required — it is written to disk after Chrome exits.
			if msg.Enabled == nil {
				writeMessage(os.Stdout, outMsg{Type: "error", Message: "missing enabled field"})
				continue
			}
			writeMessage(os.Stdout, outMsg{Type: "ok"})
			RestartChrome(*msg.Enabled)
			return

		default:
			writeMessage(os.Stdout, outMsg{Type: "error", Message: "unknown type"})
		}
	}
}
