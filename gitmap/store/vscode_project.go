package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v5/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v5/gitmap/model"
)

// UpsertVSCodeProject inserts or updates a row keyed by RootPath
// (case-insensitive). Bumps Name + LastSeenAt + UpdatedAt on conflict.
func (db *DB) UpsertVSCodeProject(rootPath, name string) error {
	if _, err := db.conn.Exec(constants.SQLUpsertVSCodeProject, rootPath, name); err != nil {
		return fmt.Errorf(constants.ErrVSCodePMUpsert, rootPath, err)
	}

	return nil
}

// ListVSCodeProjects returns every row in the VSCodeProject table.
func (db *DB) ListVSCodeProjects() ([]model.VSCodeProject, error) {
	rows, err := db.conn.Query(constants.SQLSelectAllVSCodeProjects)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrVSCodePMList, err)
	}
	defer rows.Close()

	return scanVSCodeProjectRows(rows)
}

// FindVSCodeProjectByPath returns the row matching RootPath (case-insensitive)
// or sql.ErrNoRows when missing.
func (db *DB) FindVSCodeProjectByPath(rootPath string) (model.VSCodeProject, error) {
	row := db.conn.QueryRow(constants.SQLSelectVSCodeProjectByPath, rootPath)
	p, err := scanOneVSCodeProject(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.VSCodeProject{}, sql.ErrNoRows
		}

		return model.VSCodeProject{}, fmt.Errorf(constants.ErrVSCodePMList, err)
	}

	return p, nil
}

// RenameVSCodeProjectByPath updates the Name column for the matching RootPath.
// Returns the number of rows affected so callers can detect "no match".
func (db *DB) RenameVSCodeProjectByPath(rootPath, newName string) (int64, error) {
	res, err := db.conn.Exec(constants.SQLRenameVSCodeProject, newName, rootPath)
	if err != nil {
		return 0, fmt.Errorf(constants.ErrVSCodePMRename, rootPath, err)
	}

	affected, _ := res.RowsAffected()

	return affected, nil
}

// DeleteVSCodeProjectByPath removes a row by RootPath.
func (db *DB) DeleteVSCodeProjectByPath(rootPath string) error {
	if _, err := db.conn.Exec(constants.SQLDeleteVSCodeProjectByPath, rootPath); err != nil {
		return fmt.Errorf(constants.ErrVSCodePMDelete, rootPath, err)
	}

	return nil
}

// scanOneVSCodeProject reads a single VSCodeProject row.
func scanOneVSCodeProject(row interface{ Scan(dest ...any) error }) (model.VSCodeProject, error) {
	var (
		p       model.VSCodeProject
		enabled int64
	)

	err := row.Scan(&p.ID, &p.RootPath, &p.Name, &enabled, &p.Profile,
		&p.LastSeenAt, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return p, err
	}

	p.Enabled = enabled != 0

	return p, nil
}

// scanVSCodeProjectRows reads VSCodeProject values from query result rows.
func scanVSCodeProjectRows(rows interface {
	Next() bool
	Scan(dest ...any) error
}) ([]model.VSCodeProject, error) {
	var results []model.VSCodeProject

	for rows.Next() {
		p, err := scanOneVSCodeProject(rows)
		if err != nil {
			return nil, fmt.Errorf(constants.ErrVSCodePMList, err)
		}

		results = append(results, p)
	}

	return results, nil
}
