package migrations

import "testing"

func TestLoadMigrationsIsOrderedAndChecksummed(t *testing.T) {
	items, err := load()
	if err != nil {
		t.Fatalf("load migrations: %v", err)
	}
	if len(items) == 0 {
		t.Fatal("expected at least one migration")
	}
	for i, item := range items {
		if item.checksum == "" {
			t.Fatalf("migration %d has no checksum", item.version)
		}
		if i > 0 && items[i-1].version >= item.version {
			t.Fatalf("migrations are not strictly ordered")
		}
	}
}
