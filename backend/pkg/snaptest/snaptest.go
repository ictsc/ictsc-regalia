package snaptest

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/google/go-cmp/cmp"
)

func Match[T any](t *testing.T, got T, opts ...cmp.Option) {
	t.Helper()

	snapPath := path.Join("__snapshots__", t.Name()+".snap.json")
	file, err := os.Open(snapPath) //nolint:gosec
	if err != nil && errors.Is(err, os.ErrNotExist) {
		if !checkUpdate() {
			t.Error("Snapshot file does not exist. Update snapshot with UPDATE_SNAPSHOT=1")
			return
		}
		if err := updateSnapshot(snapPath, got); err != nil {
			t.Errorf("failed to update snapshot: %v", err)
		}
		return
	} else if err != nil {
		t.Errorf("Failed to open snapshot file: %v", err)
		return
	}
	defer file.Close() //nolint:errcheck

	var want T
	if err := json.NewDecoder(file).Decode(&want); err != nil {
		t.Errorf("Failed to decode snapshot file: %v", err)
		return
	}

	if diff := cmp.Diff(want, got); diff != "" {
		if !checkUpdate() {
			t.Errorf("Mismatch (-want +got):\n%s\nUpdate snapshot with UPDATE_SNAPSHOT=1", diff)
			return
		} else {
			if err := updateSnapshot(snapPath, got); err != nil {
				t.Errorf("Failed to update snapshot: %v", err)
			}
		}
	}
}

func checkUpdate() bool {
	return os.Getenv("UPDATE_SNAPSHOT") == "1"
}

const (
	snapshotFilePerm = os.FileMode(0o644)
	snapshotDirPerm  = os.FileMode(0o755)
)

func updateSnapshot(snapPath string, got any) error {
	gotJSON, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal snapshot")
	}

	if err := os.MkdirAll(path.Dir(snapPath), snapshotDirPerm); err != nil {
		return errors.Wrap(err, "failed to create snapshot directory")
	}
	if err := os.WriteFile(snapPath, gotJSON, snapshotFilePerm); err != nil {
		return errors.Wrap(err, "failed to write snapshot file")
	}

	return nil
}
