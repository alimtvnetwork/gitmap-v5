package formatter

// Shared warning emitter used by every record-aware writer (WriteJSON,
// WriteCSV) before serializing. The default destination is os.Stderr,
// but tests can swap it via SetValidationSink to capture warnings.

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/alimtvnetwork/gitmap-v5/gitmap/model"
)

// validationSink is the io.Writer that receives per-issue warning lines.
// Guarded by sinkMu so concurrent writers + tests can swap it safely.
var (
	validationSink io.Writer = os.Stderr
	sinkMu         sync.RWMutex
)

// SetValidationSink redirects validation warnings to w. Pass os.Stderr to
// restore the default. Returns the previous sink so tests can defer-restore.
func SetValidationSink(w io.Writer) io.Writer {
	sinkMu.Lock()
	defer sinkMu.Unlock()
	prev := validationSink
	validationSink = w

	return prev
}

// emitValidationWarnings runs the validator over records and writes one
// `gitmap: validation: <issue>` line per finding to the active sink. It
// never returns an error — by policy the write proceeds regardless.
func emitValidationWarnings(records []model.ScanRecord) {
	issues := ValidateRecords(records)
	if len(issues) == 0 {
		return
	}

	sinkMu.RLock()
	w := validationSink
	sinkMu.RUnlock()

	for _, issue := range issues {
		fmt.Fprintf(w, "gitmap: validation: %s\n", issue.String())
	}
}
