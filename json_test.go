// Copyright ©2026 xiayoudi. All rights reserved.
// Author: xiayoudi
// Email: ur@xiaud.com

package ud

import (
	"os"
	"path/filepath"
	"testing"
)

type Error struct {
	Code int
	Msg  string
}

func TestJSON(t *testing.T) {
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "config.json")

	input := &Error{Code: 2026, Msg: "json"}

	t.Run("Write", func(t *testing.T) {
		err := WriteJSON(jsonPath, input)
		if err != nil {
			t.Fatalf("failed to write: %v", err)
		}
	})

	t.Run("Read", func(t *testing.T) {
		var e Error
		err := ReadJSON(jsonPath, &e)
		if err != nil {
			t.Fatalf("failed to read: %v", err)
		}

		t.Logf("Read Success: %+v", e)
	})

	t.Cleanup(func() {
		os.Remove(jsonPath)
		t.Log("temporary file removed.")
	})
}
