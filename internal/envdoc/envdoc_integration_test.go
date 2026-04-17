package envdoc_test

import (
	"testing"

	"github.com/user/envdiff/internal/envdoc"
)

func TestParse_FullPipeline(t *testing.T) {
	content := `# Application port
PORT=3000

# Database URL
DB_URL=postgres://localhost/mydb # primary db

SECRET=hunter2
`
	path := writeTempEnv(t, content)

	entries, err := envdoc.Parse(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	docs := envdoc.ToMap(entries)

	if docs["PORT"] != "Application port" {
		t.Errorf("PORT comment: got %q", docs["PORT"])
	}
	if docs["DB_URL"] != "primary db" {
		t.Errorf("DB_URL comment: got %q", docs["DB_URL"])
	}
	if docs["SECRET"] != "" {
		t.Errorf("SECRET should have no comment, got %q", docs["SECRET"])
	}

	// Verify values are intact
	for _, e := range entries {
		if e.Key == "DB_URL" && e.Value != "postgres://localhost/mydb" {
			t.Errorf("DB_URL value: got %q", e.Value)
		}
	}
}
