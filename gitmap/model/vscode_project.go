package model

// VSCodeProject is one row in the VSCodeProject table — the gitmap-side
// source of truth for entries synced into VS Code Project Manager's
// projects.json. `tags` and `paths` are NOT stored here on purpose; they
// live only inside projects.json and are preserved across syncs.
type VSCodeProject struct {
	ID         int64  `json:"id"`
	RootPath   string `json:"rootPath"`
	Name       string `json:"name"`
	Enabled    bool   `json:"enabled"`
	Profile    string `json:"profile"`
	LastSeenAt string `json:"lastSeenAt"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}
