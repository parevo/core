package export

import (
	"encoding/json"
	"testing"
)

func TestToJSON(t *testing.T) {
	p := NewPayload("u1")
	p.Profile = map[string]any{"email": "a@b.com"}
	p.Sessions = []map[string]any{{"id": "s1"}}

	b, err := ToJSON(p)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Payload
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.UserID != "u1" {
		t.Errorf("user_id: got %q", decoded.UserID)
	}
	if decoded.Profile["email"] != "a@b.com" {
		t.Errorf("profile: got %v", decoded.Profile)
	}
	if decoded.ExportedAt.IsZero() {
		t.Error("exported_at should be set")
	}
}

func TestToReader(t *testing.T) {
	p := NewPayload("u1")
	r, err := ToReader(p)
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fatal("reader is nil")
	}
}
