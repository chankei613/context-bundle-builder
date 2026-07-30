package resolve

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/chankei613/context-bundle-builder/internal/db"
)

func TestObsidianNoteRejectsPathTraversal(t *testing.T) {
	vault := t.TempDir()
	if err := os.MkdirAll(filepath.Join(vault, "notes"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(vault, "notes", "ok.md"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	// vaultの外側に「秘密」ファイルを作り、traversalで読めないことを確認する
	if err := os.WriteFile(filepath.Join(filepath.Dir(vault), "secret.md"), []byte("SECRET"), 0o644); err != nil {
		t.Fatal(err)
	}

	opts := Options{ObsidianVaultRoot: vault}

	res := Resolve(context.Background(), []db.ContextRef{
		{Kind: db.RefKindObsidianNote, Ref: "../secret.md"},
		{Kind: db.RefKindObsidianNote, Ref: "notes/ok.md"},
	}, opts)

	if len(res.Errors) != 1 {
		t.Fatalf("expected 1 error (traversal blocked), got %d: %+v", len(res.Errors), res.Errors)
	}
	if res.Errors[0].Ref.Ref != "../secret.md" {
		t.Fatalf("expected the traversal ref to be the one that failed, got %+v", res.Errors[0])
	}
	if len(res.Items) != 1 || res.Items[0].Content != "hello" {
		t.Fatalf("expected the in-vault note to resolve, got %+v", res.Items)
	}
}

func TestObsidianNoteRequiresVaultRoot(t *testing.T) {
	res := Resolve(context.Background(), []db.ContextRef{
		{Kind: db.RefKindObsidianNote, Ref: "notes/ok.md"},
	}, Options{})

	if len(res.Errors) != 1 {
		t.Fatalf("expected an unconfigured-vault error, got %+v", res)
	}
}

func TestTaskOutputSoftFailsWhenUnconfigured(t *testing.T) {
	res := Resolve(context.Background(), []db.ContextRef{
		{Kind: db.RefKindTaskOutput, Ref: "task#1.output_key"},
	}, Options{})

	if len(res.Errors) != 1 || len(res.Items) != 0 {
		t.Fatalf("expected a soft error (not a panic/crash), got errors=%d items=%d", len(res.Errors), len(res.Items))
	}
}

func TestEstimateTokens(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"", 0},
		{"abcd", 1},
		{"abcde", 2},
		{"12345678", 2},
	}
	for _, c := range cases {
		if got := estimateTokens(c.in); got != c.want {
			t.Errorf("estimateTokens(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestFileConnectorMissingFileIsSoftError(t *testing.T) {
	res := Resolve(context.Background(), []db.ContextRef{
		{Kind: db.RefKindFile, Ref: filepath.Join(t.TempDir(), "does-not-exist.md")},
	}, Options{})

	if len(res.Errors) != 1 || len(res.Items) != 0 {
		t.Fatalf("expected a soft error for missing file, got errors=%d items=%d", len(res.Errors), len(res.Items))
	}
}
