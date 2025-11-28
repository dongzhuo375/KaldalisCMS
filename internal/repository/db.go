package repository

import (
	"KaldalisCMS/internal/model"
	"sync"
)

// InMemoryDB simulates a database using a map.
var InMemoryDB = struct {
	sync.RWMutex
	posts   map[int]model.Post
	counter int
}{
	posts:   make(map[int]model.Post),
	counter: 0,
}

// init adds some dummy data to the in-memory database.
func init() {
	InMemoryDB.Lock()
	defer InMemoryDB.Unlock()

	// Add a dummy post
	InMemoryDB.counter++
	InMemoryDB.posts[InMemoryDB.counter] = model.Post{
		ID:      InMemoryDB.counter,
		Title:   "Welcome to KaldalisCMS!",
		Content: "This is your first post. Start creating!",
	}
}
