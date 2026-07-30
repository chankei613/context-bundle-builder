package resolve

type resolveError struct{ msg string }

func (e *resolveError) Error() string { return e.msg }

var (
	errUnknownKind         = &resolveError{"unknown ref kind"}
	errFileEmpty           = &resolveError{"file ref is empty"}
	errURLEmpty            = &resolveError{"url ref is empty"}
	errObsidianVaultUnset  = &resolveError{"obsidian vault root is not configured"}
	errObsidianOutsideRoot = &resolveError{"note path escapes the vault root"}
	errTaskOutputUnset     = &resolveError{"task_output source is not configured (set an adapter URL)"}
)
