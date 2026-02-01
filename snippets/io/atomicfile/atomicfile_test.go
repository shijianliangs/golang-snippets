package atomicfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "cfg", "app.conf")

	if err := WriteFile(p, []byte("v=1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "v=1\n" {
		t.Fatalf("got %q", string(b))
	}
	st, err := os.Stat(p)
	if err != nil {
		t.Fatal(err)
	}
	if st.Mode().Perm() != 0o600 {
		t.Fatalf("perm=%o", st.Mode().Perm())
	}

	// Overwrite
	if err := WriteFile(p, []byte("v=2\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	b, _ = os.ReadFile(p)
	if string(b) != "v=2\n" {
		t.Fatalf("got %q", string(b))
	}
}
