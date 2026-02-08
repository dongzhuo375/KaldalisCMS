package entity

// Tag represents a post tag.
//
// Note: database-specific mapping is intentionally omitted for now.
type Tag struct {
	ID   uint
	Name string
}
