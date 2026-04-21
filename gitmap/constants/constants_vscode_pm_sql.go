package constants

// SQL for the VSCodeProject table — DB source of truth for the
// VS Code Project Manager sync (v3.38.0+).
//
// `tags` and `paths` are not stored in the DB on purpose; they live only in
// `projects.json` and are preserved verbatim across syncs so user edits in
// the extension UI are never clobbered.

const TableVSCodeProject = "VSCodeProject"

const SQLCreateVSCodeProject = `CREATE TABLE IF NOT EXISTS VSCodeProject (
	VSCodeProjectId INTEGER PRIMARY KEY AUTOINCREMENT,
	RootPath        TEXT NOT NULL,
	Name            TEXT NOT NULL,
	Enabled         INTEGER NOT NULL DEFAULT 1,
	Profile         TEXT NOT NULL DEFAULT '',
	LastSeenAt      TEXT DEFAULT CURRENT_TIMESTAMP,
	CreatedAt       TEXT DEFAULT CURRENT_TIMESTAMP,
	UpdatedAt       TEXT DEFAULT CURRENT_TIMESTAMP
)`

// COLLATE NOCASE so Windows path matching is case-insensitive while
// staying byte-exact on Unix when the user happens to use the same case.
const SQLCreateVSCodeProjectRootPathIndex = `CREATE UNIQUE INDEX IF NOT EXISTS UX_VSCodeProject_RootPath ON VSCodeProject(RootPath COLLATE NOCASE)`

const SQLDropVSCodeProject = "DROP TABLE IF EXISTS VSCodeProject"

const (
	SQLUpsertVSCodeProject = `INSERT INTO VSCodeProject (RootPath, Name, Enabled, Profile, LastSeenAt, UpdatedAt)
		VALUES (?, ?, 1, '', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT(RootPath) DO UPDATE SET
			Name=excluded.Name,
			LastSeenAt=CURRENT_TIMESTAMP,
			UpdatedAt=CURRENT_TIMESTAMP`

	SQLSelectAllVSCodeProjects = `SELECT VSCodeProjectId, RootPath, Name, Enabled, Profile, LastSeenAt, CreatedAt, UpdatedAt
		FROM VSCodeProject ORDER BY UpdatedAt DESC, RootPath ASC`

	SQLSelectVSCodeProjectByPath = `SELECT VSCodeProjectId, RootPath, Name, Enabled, Profile, LastSeenAt, CreatedAt, UpdatedAt
		FROM VSCodeProject WHERE RootPath = ? COLLATE NOCASE`

	SQLRenameVSCodeProject = `UPDATE VSCodeProject
		SET Name = ?, UpdatedAt = CURRENT_TIMESTAMP
		WHERE RootPath = ? COLLATE NOCASE`

	SQLDeleteVSCodeProjectByPath = `DELETE FROM VSCodeProject WHERE RootPath = ? COLLATE NOCASE`
)

// Error messages.
const (
	ErrVSCodePMUpsert       = "failed to upsert VSCodeProject %q: %v"
	ErrVSCodePMList         = "failed to list VSCodeProject rows: %v"
	ErrVSCodePMRename       = "failed to rename VSCodeProject %q: %v"
	ErrVSCodePMDelete       = "failed to delete VSCodeProject %q: %v"
)
