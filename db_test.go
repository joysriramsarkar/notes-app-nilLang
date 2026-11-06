package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStoreCRUDAndAtomicSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "notes.json")
	s, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if len(s.GetAll()) == 0 {
		t.Fatal("expected seeded notes on fresh store")
	}
	created, err := s.Create(Note{Title: "hello <b>", Content: "body", Category: "work"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("created note has empty ID")
	}
	got, ok := s.Get(created.ID)
	if !ok || got.Title != "hello <b>" {
		t.Fatalf("Get mismatch: %+v ok=%v", got, ok)
	}
	updated, err := s.Update(created.ID, Note{Title: "updated", Content: "c2", Category: "work"})
	if err != nil || updated.Title != "updated" {
		t.Fatalf("Update: %v %+v", err, updated)
	}
	// Reload from disk: proves atomic save persisted.
	s2, err := NewStore(path)
	if err != nil {
		t.Fatalf("reload NewStore: %v", err)
	}
	got2, ok := s2.Get(created.ID)
	if !ok || got2.Title != "updated" {
		t.Fatalf("persisted note mismatch: %+v ok=%v", got2, ok)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Fatal("tmp file left behind after atomic save")
	}
	if err := s2.Delete(created.ID); err != nil {
		t.Fatalf("Delete (trash): %v", err)
	}
	trashed, ok := s2.Get(created.ID)
	if !ok || !trashed.IsTrashed {
		t.Fatalf("expected trashed note, got %+v ok=%v", trashed, ok)
	}
	if err := s2.Delete(created.ID); err != nil {
		t.Fatalf("Delete (purge): %v", err)
	}
	if _, ok := s2.Get(created.ID); ok {
		t.Fatal("expected note purged after second delete")
	}
}
