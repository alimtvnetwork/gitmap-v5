package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v5/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v5/gitmap/store"
	"github.com/alimtvnetwork/gitmap-v5/gitmap/vscodepm"
)

// runCode implements `gitmap code [alias] [path]`.
//
// Resolution order for the target rootPath:
//  1. Positional `path` arg, if supplied.
//  2. `git rev-parse --show-toplevel` if invoked inside a Git repo.
//  3. The current working directory.
//
// The alias defaults to the folder basename and can be overridden by the
// first positional arg. After resolving the path the command upserts the
// (rootPath, name) pair into both the gitmap DB and projects.json, then
// launches VS Code on the path.
func runCode(args []string) {
	checkHelp(constants.CmdCode, args)
	alias, pathArg := parseCodeArgs(args)

	rootPath, err := resolveCodeRootPath(pathArg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if alias == "" {
		alias = filepath.Base(rootPath)
	}

	upsertCodeEntry(rootPath, alias)
	syncCodeEntry(rootPath, alias)
	openInVSCode(rootPath)
}

// parseCodeArgs returns (alias, path) from the positional args.
// Both are optional. Extra args are an error.
func parseCodeArgs(args []string) (string, string) {
	if len(args) == 0 {
		return "", ""
	}

	if len(args) == 1 {
		return args[0], ""
	}

	if len(args) == 2 {
		return args[0], args[1]
	}

	fmt.Fprintln(os.Stderr, "usage: gitmap code [alias] [path]")
	os.Exit(2)

	return "", ""
}

// resolveCodeRootPath picks the rootPath per the documented precedence.
func resolveCodeRootPath(pathArg string) (string, error) {
	if pathArg != "" {
		return absoluteExisting(pathArg)
	}

	if root, err := gitTopLevel(); err == nil {
		return root, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot determine current directory: %w", err)
	}

	return absoluteExisting(cwd)
}

// absoluteExisting cleans the path, returns it absolute, and verifies it
// exists. Non-existent paths are an error so we never write garbage rows.
func absoluteExisting(p string) (string, error) {
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", fmt.Errorf("cannot resolve absolute path %q: %w", p, err)
	}

	if _, err := os.Stat(abs); err != nil {
		return "", fmt.Errorf("path does not exist: %s", abs)
	}

	return abs, nil
}

// upsertCodeEntry persists the row into the VSCodeProject table.
func upsertCodeEntry(rootPath, name string) {
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgDBUpsertFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgDBUpsertFailed, err)
		os.Exit(1)
	}

	if err := db.UpsertVSCodeProject(rootPath, name); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

// syncCodeEntry pushes the (rootPath, name) tuple into projects.json.
// Soft-fails when VS Code or the extension is missing.
func syncCodeEntry(rootPath, name string) {
	summary, err := vscodepm.Sync([]vscodepm.Pair{{RootPath: rootPath, Name: name}})
	if err != nil {
		reportVSCodePMSoftError(err)

		return
	}

	fmt.Printf(constants.MsgVSCodePMSyncSummary,
		summary.Added, summary.Updated, summary.Unchanged, summary.Total)
}
