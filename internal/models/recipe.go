package models

// swagger:model
type Category struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// swagger:model
type Ingredient struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Unit string `db:"unit" json:"unit"`
}

// swagger:model
type Tag struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// swagger:model
type Recipe struct {
	ID           int64    `db:"id" json:"id"`
	Title        string   `db:"title" json:"title"`
	Description  string   `db:"description" json:"description"`
	Instructions string   `db:"instructions" json:"instructions"`
	ImageURL     string   `db:"image_url" json:"image_url"`
	CategoryID   int64    `db:"category_id" json:"category_id"`
	Tags         []string `json:"tags,omitempty"`
}

// swagger:model
type SearchRequest struct {
	Ingredients []string `json:"ingredients"`
	Tags        []string `json:"tags"`
}

// swagger:model
type User struct {
	ID       int64  `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password,omitempty"` // хранить хэш, не plain
}

// swagger:model
type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// swagger:model
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// swagger:model
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
