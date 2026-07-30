package resolve

import (
	"path/filepath"
	"strings"
)

func resolveObsidianNote(notePath string, opts Options) (string, error) {
	if opts.ObsidianVaultRoot == "" {
		return "", errObsidianVaultUnset
	}
	if notePath == "" {
		return "", errFileEmpty
	}

	root, err := filepath.Abs(opts.ObsidianVaultRoot)
	if err != nil {
		return "", err
	}
	full := filepath.Join(root, notePath)
	full, err = filepath.Abs(full)
	if err != nil {
		return "", err
	}

	// vaultルート配下から出ないことを検証する（パストラバーサル対策）。
	if full != root && !strings.HasPrefix(full, root+string(filepath.Separator)) {
		return "", errObsidianOutsideRoot
	}

	return resolveFile(full, opts)
}
