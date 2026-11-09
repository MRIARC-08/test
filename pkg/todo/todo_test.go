package todo_test

import (
	"testing"

	"github.com/MRIARC-08/test/pkg/todo"
)

func TestTodoLifecycle(t *testing.T) {
	store := todo.NewMemoryStore()

	// 1. Test Create
	item, err := store.Create("Deploy Biscuit Watchdog", "Ensure port 8005 webhook works")
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	if item.ID != 1 {
		t.Errorf("expected ID 1, got %d", item.ID)
	}
	if item.Completed {
		t.Errorf("expected new item to be incomplete")
	}

	// 2. Test Get
	got, err := store.Get(1)
	if err != nil {
		t.Fatalf("failed to get item: %v", err)
	}
	if got.Title != "Deploy Biscuit Watchdog" {
		t.Errorf("title mismatch: %s", got.Title)
	}

	// 3. Test Complete
	completed, err := store.Complete(1)
	if err != nil {
		t.Fatalf("failed to complete item: %v", err)
	}
	if !completed.Completed {
		t.Errorf("expected item to be marked completed")
	}

	// 4. Test List
	items := store.List()
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	// 5. Test Delete
	if err := store.Delete(1); err != nil {
		t.Fatalf("failed to delete item: %v", err)
	}
	if len(store.List()) != 0 {
		t.Errorf("expected 0 items after deletion")
	}
}

func TestValidationErrors(t *testing.T) {
	store := todo.NewMemoryStore()

	_, err := store.Create("", "empty title")
	if err != todo.ErrEmptyTitle {
		t.Errorf("expected ErrEmptyTitle, got %v", err)
	}

	_, err = store.Get(999)
	if err == nil {
		t.Errorf("expected error for non-existent item")
	}
}
