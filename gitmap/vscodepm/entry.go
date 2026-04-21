package vscodepm

// Entry mirrors one object in projects.json. Field order and JSON tags
// match the alefragnani.project-manager schema exactly so encoded files
// stay diff-friendly with manual edits and the sample fixture
// (spec/01-vscode-project-manager-sync/sample-projects.json).
type Entry struct {
	Name     string   `json:"name"`
	RootPath string   `json:"rootPath"`
	Paths    []string `json:"paths"`
	Tags     []string `json:"tags"`
	Enabled  bool     `json:"enabled"`
	Profile  string   `json:"profile"`
}

// SyncSummary is returned from Sync to describe what changed.
type SyncSummary struct {
	Added     int
	Updated   int
	Unchanged int
	Total     int
}

// newEntry builds a default Entry for a freshly registered (rootPath, name).
// Tags and Paths are always emitted as non-nil empty slices so the encoded
// JSON contains `[]` rather than `null` (matches the sample fixture).
func newEntry(rootPath, name string) Entry {
	return Entry{
		Name:     name,
		RootPath: rootPath,
		Paths:    []string{},
		Tags:     []string{},
		Enabled:  true,
		Profile:  "",
	}
}

// ensureSlices replaces nil slices with empty ones so re-encoded JSON
// preserves `[]` even on entries that arrived with `null` from disk.
func ensureSlices(e Entry) Entry {
	if e.Paths == nil {
		e.Paths = []string{}
	}

	if e.Tags == nil {
		e.Tags = []string{}
	}

	return e
}
