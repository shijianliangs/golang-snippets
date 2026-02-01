package atomicfile

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// WriteFile writes data to filename atomically.
//
// It writes to a temp file in the same directory and then renames it over the
// destination path. On POSIX filesystems, os.Rename is atomic.
//
// If perm is 0, it uses 0o644.
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	return WriteFileFunc(filename, perm, func(w io.Writer) error {
		_, err := w.Write(data)
		return err
	})
}

// WriteFileFunc atomically writes a file by invoking write(w) to stream the contents.
//
// The temp file is removed on error.
func WriteFileFunc(filename string, perm os.FileMode, write func(w io.Writer) error) error {
	if write == nil {
		return errors.New("atomicfile: write func is nil")
	}
	if perm == 0 {
		perm = 0o644
	}

	dir := filepath.Dir(filename)
	base := filepath.Base(filename)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("atomicfile: mkdir: %w", err)
	}

	f, err := os.CreateTemp(dir, "."+base+".*.tmp")
	if err != nil {
		return fmt.Errorf("atomicfile: create temp: %w", err)
	}
	name := f.Name()
	defer func() {
		_ = os.Remove(name)
	}()

	// Set file mode before rename.
	if err := f.Chmod(perm); err != nil {
		_ = f.Close()
		return fmt.Errorf("atomicfile: chmod: %w", err)
	}

	if err := write(f); err != nil {
		_ = f.Close()
		return fmt.Errorf("atomicfile: write: %w", err)
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return fmt.Errorf("atomicfile: sync: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("atomicfile: close: %w", err)
	}

	if err := os.Rename(name, filename); err != nil {
		return fmt.Errorf("atomicfile: rename: %w", err)
	}
	// Best-effort directory sync to improve durability.
	if d, err := os.Open(dir); err == nil {
		_ = d.Sync()
		_ = d.Close()
	}
	return nil
}
