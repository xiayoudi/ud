// Copyright ©2026 xiayoudi. All rights reserved.
// Author: xiayoudi
// Email: ur@xiaud.com

package ud

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
)

func IsPointer(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Pointer
}

func WriteJSON(path string, data any) error {
	if path == "" {
		return Err("path can't be empty")
	}

	if !IsPointer(data) {
		return Err("data must be a pointer")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	tmpPath := path + ".tmp"
	f, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		Wrapf(err, "failed to create tmpPath %s", tmpPath)
	}

	bw := bufio.NewWriter(f)
	enc := json.NewEncoder(bw)
	enc.SetIndent("", "  ")

	if err = enc.Encode(data); err != nil {
		os.Remove(tmpPath)
		return Wrap(err, "failed to encode data")
	}

	err = func() error {
		defer f.Close()

		if err = bw.Flush(); err != nil {
			return Wrap(err, "failed to flush buffer")
		}

		if err = f.Sync(); err != nil {
			return Wrap(err, "failed to sync file")
		}

		return nil
	}()

	if err != nil {
		os.Remove(tmpPath)
		return Wrap(err, "failed to write file")
	}

	if err = os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return Wrap(err, "failed to rename file")
	}

	return nil
}

func ReadJSON(path string, data any) error {
	if path == "" {
		return Err("path can't be empty")
	}

	if !IsPointer(data) {
		return Err("data must be a pointer")
	}

	f, err := os.Open(path)
	if err != nil {
		return Wrap(err, "failed to open file")
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()

	if err = dec.Decode(data); err != nil {
		return Wrap(err, "failed to decode data")
	}

	return nil
}
