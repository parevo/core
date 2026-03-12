package memory

import (
	"bytes"
	"context"
	"io"
	"testing"
)

func TestStore_PutGet(t *testing.T) {
	s := NewStore()
	ctx := context.Background()

	err := s.Put(ctx, "b1", "k1", bytes.NewReader([]byte("hello")), "text/plain")
	if err != nil {
		t.Fatal(err)
	}
	rc, err := s.Get(ctx, "b1", "k1")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rc.Close() }()
	data, _ := io.ReadAll(rc)
	if string(data) != "hello" {
		t.Errorf("want hello, got %s", data)
	}
}

func TestStore_Delete(t *testing.T) {
	s := NewStore()
	ctx := context.Background()
	_ = s.Put(ctx, "b1", "k1", bytes.NewReader([]byte("x")), "")
	err := s.Delete(ctx, "b1", "k1")
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.Get(ctx, "b1", "k1")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestStore_List(t *testing.T) {
	s := NewStore()
	ctx := context.Background()
	_ = s.Put(ctx, "b1", "a/1", bytes.NewReader([]byte("x")), "")
	_ = s.Put(ctx, "b1", "a/2", bytes.NewReader([]byte("y")), "")
	_ = s.Put(ctx, "b1", "b/1", bytes.NewReader([]byte("z")), "")
	infos, err := s.List(ctx, "b1", "a/")
	if err != nil {
		t.Fatal(err)
	}
	if len(infos) != 2 {
		t.Errorf("want 2, got %d", len(infos))
	}
}
