package uuid

import (
	"testing"
)

func TestGenerateUUID_Uniqueness(t *testing.T) {
	const iterations = 10000

	seen := make(map[string]struct{}, iterations)

	for i := 0; i < iterations; i++ {
		uuid := New()

		idStr, err := uuid.Generate()
		if err != nil {
			t.Fatal(err)
		}

		if _, ok := seen[idStr]; ok {
			t.Fatalf("Duplicate ULID generated: %s", idStr)
		}

		seen[idStr] = struct{}{}
	}
}

func TestGenerateULID_ValidFormat(t *testing.T) {
	uuid := New()

	idStr, err := uuid.Generate()
	if err != nil {
		t.Fatal(err)
	}

	if len(idStr) != 26 {
		t.Errorf("ULID should be 26 characters, got %d", len(idStr))
	}
}
