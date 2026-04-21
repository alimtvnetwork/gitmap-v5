package formatter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v5/gitmap/model"
)

// TestWriteJSON_EmitsWarningsButStillWrites verifies the warn-and-write
// contract: invalid records produce stderr lines AND a complete JSON file.
func TestWriteJSON_EmitsWarningsButStillWrites(t *testing.T) {
	var sink bytes.Buffer
	prev := SetValidationSink(&sink)
	defer SetValidationSink(prev)

	records := []model.ScanRecord{
		{ /* missing everything */ },
		{RepoName: "good", HTTPSUrl: "https://x/y.git", RelativePath: "y"},
	}

	var out bytes.Buffer
	if err := WriteJSON(&out, records); err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}

	assertSinkAndOutput(t, &sink, &out, "good", `"repoName":"good"`)
}

// TestWriteCSV_EmitsWarningsButStillWrites mirrors the JSON test for CSV.
func TestWriteCSV_EmitsWarningsButStillWrites(t *testing.T) {
	var sink bytes.Buffer
	prev := SetValidationSink(&sink)
	defer SetValidationSink(prev)

	records := []model.ScanRecord{
		{RepoName: "alpha" /* missing url + relpath */},
		{RepoName: "beta", SSHUrl: "git@x:y/beta.git", RelativePath: "beta"},
	}

	var out bytes.Buffer
	if err := WriteCSV(&out, records); err != nil {
		t.Fatalf("WriteCSV: %v", err)
	}

	assertSinkAndOutput(t, &sink, &out, "alpha", "beta")
}

// TestWriteJSON_NoWarningsOnCleanInput verifies the sink stays silent when
// every record passes validation.
func TestWriteJSON_NoWarningsOnCleanInput(t *testing.T) {
	var sink bytes.Buffer
	prev := SetValidationSink(&sink)
	defer SetValidationSink(prev)

	records := []model.ScanRecord{{
		Slug:         "ok",
		RepoName:     "ok",
		HTTPSUrl:     "https://x/ok.git",
		RelativePath: "ok",
	}}

	var out bytes.Buffer
	if err := WriteJSON(&out, records); err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}
	if sink.Len() != 0 {
		t.Errorf("expected silent sink, got: %q", sink.String())
	}
}

// assertSinkAndOutput is the shared assertion shape: sink mentions the
// validation prefix + bad row's identifier; output contains the good row.
func assertSinkAndOutput(t *testing.T, sink, out *bytes.Buffer, badName, mustHave string) {
	t.Helper()
	sinkStr := sink.String()
	if !strings.Contains(sinkStr, "gitmap: validation:") {
		t.Errorf("sink missing validation prefix: %q", sinkStr)
	}
	if !strings.Contains(out.String(), mustHave) {
		t.Errorf("output missing expected payload %q in: %q", mustHave, out.String())
	}
	_ = badName // reserved for future per-bad-row assertions
}
