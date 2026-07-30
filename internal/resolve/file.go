package resolve

import "os"

func resolveFile(path string, opts Options) (string, error) {
	if path == "" {
		return "", errFileEmpty
	}

	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if info.Size() > opts.MaxBytes {
		f, err := os.Open(path)
		if err != nil {
			return "", err
		}
		defer func() { _ = f.Close() }()
		buf := make([]byte, opts.MaxBytes)
		n, err := f.Read(buf)
		if err != nil {
			return "", err
		}
		return string(buf[:n]) + "\n\n[...truncated, exceeds size limit...]", nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
