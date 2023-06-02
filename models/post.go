package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	ID          uint           `json:"id,omitempty"`
	Title       string         `json:"title,omitempty" search:"title" filter:"title"`
	Content     string         `json:"content,omitempty" search:"content"`
	Tags        pq.StringArray `gorm:"type:text[]" json:"tags"`
	CreatedByID uint           `json:"created_by_id,omitempty"`
	CreatedBy   User           `json:"-"`
	CreatedAt   time.Time      `json:"created_at,omitempty" filter:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type Comment struct {
	ID          uint           `json:"id,omitempty"`
	Content     string         `json:"content,omitempty" search:"content"`
	PostID      uint           `json:"post_id,omitempty"`
	Post        Post           `json:"-"`
	CreatedByID uint           `json:"created_by_id,omitempty"`
	CreatedBy   User           `json:"-"`
	CreatedAt   time.Time      `json:"created_at,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
