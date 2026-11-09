package findings

import "testing"

func TestStatusVocabulary(t *testing.T) {
	for _, s := range []string{"Open", "In Progress", "Resolved"} {
		if s == "" {
			t.Fatal("valid status")
		}
	}
	if "Invalid" == "Open" || "Invalid" == "In Progress" || "Invalid" == "Resolved" {
		t.Fatal("invalid status accepted")
	}
}
