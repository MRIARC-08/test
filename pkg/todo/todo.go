package todo

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrNotFound      = errors.New("todo item not found")
	ErrEmptyTitle    = errors.New("todo title cannot be empty")
	ErrAlreadyDone   = errors.New("todo is already marked as completed")
)

// Item represents a single todo entry.
type Item struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Store defines thread-safe storage for Todo items.
type Store interface {
	Create(title, description string) (*Item, error)
	Get(id int) (*Item, error)
	List() []*Item
	Complete(id int) (*Item, error)
	Delete(id int) error
}

// MemoryStore provides an in-memory, thread-safe Store implementation.
type MemoryStore struct {
	mu     sync.RWMutex
	items  map[int]*Item
	nextID int
}

// NewMemoryStore creates a new in-memory todo store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		items:  make(map[int]*Item),
		nextID: 1,
	}
}

// Create inserts a new todo item into the store.
func (s *MemoryStore) Create(title, description string) (*Item, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	item := &Item{
		ID:          s.nextID,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	s.items[item.ID] = item
	s.nextID++
	return item, nil
}

// Get retrieves an item by its ID.
func (s *MemoryStore) Get(id int) (*Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.items[id]
	if !exists {
		return nil, fmt.Errorf("%w: id %d", ErrNotFound, id)
	}
	return item, nil
}

// List returns all todo items.
func (s *MemoryStore) List() []*Item {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*Item, 0, len(s.items))
	for _, item := range s.items {
		result = append(result, item)
	}
	return result
}

// Complete marks an existing item as completed.
func (s *MemoryStore) Complete(id int) (*Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.items[id]
	if !exists {
		return nil, fmt.Errorf("%w: id %d", ErrNotFound, id)
	}
	if item.Completed {
		return nil, ErrAlreadyDone
	}

	item.Completed = true
	item.UpdatedAt = time.Now().UTC()
	return item, nil
}

// Delete removes an item by ID.
func (s *MemoryStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[id]; !exists {
		return fmt.Errorf("%w: id %d", ErrNotFound, id)
	}
	delete(s.items, id)
	return nil
}
