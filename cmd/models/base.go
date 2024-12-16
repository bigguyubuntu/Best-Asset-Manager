package mdls

// contains creted at and updated at fields
type dateTracking struct {
	CreatedAt string `json:"created_at,omitempty"` //todo update this to timestame type
	UpdatedAt string `json:"updated_at,omitempty"` //todo update this to timestame type
}
