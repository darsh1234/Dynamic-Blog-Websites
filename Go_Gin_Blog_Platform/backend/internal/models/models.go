package models

import "time"

type Role string

type PostStatus string

const (
	RoleAdmin  Role = "admin"
	RoleAuthor Role = "author"
	RoleReader Role = "reader"
)

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
)

type User struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	Role         Role      `gorm:"type:text;not null"`
	CreatedAt    time.Time `gorm:"not null;default:now()"`
	UpdatedAt    time.Time `gorm:"not null;default:now()"`
}

type Post struct {
	ID        string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	AuthorID  string     `gorm:"type:uuid;not null;index"`
	Title     string     `gorm:"not null"`
	Content   string     `gorm:"not null"`
	Status    PostStatus `gorm:"type:text;not null;default:published"`
	CreatedAt time.Time  `gorm:"not null;default:now()"`
	UpdatedAt time.Time  `gorm:"not null;default:now()"`
}

type RefreshToken struct {
	ID        string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string     `gorm:"type:uuid;not null;index"`
	TokenHash string     `gorm:"not null"`
	ExpiresAt time.Time  `gorm:"not null"`
	RevokedAt *time.Time `gorm:"index"`
	CreatedAt time.Time  `gorm:"not null;default:now()"`
}

type PasswordResetToken struct {
	ID        string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string     `gorm:"type:uuid;not null;index"`
	TokenHash string     `gorm:"not null"`
	ExpiresAt time.Time  `gorm:"not null"`
	UsedAt    *time.Time `gorm:"index"`
	CreatedAt time.Time  `gorm:"not null;default:now()"`
}
