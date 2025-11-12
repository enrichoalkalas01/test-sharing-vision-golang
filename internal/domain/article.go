package domain

import (
	"errors"
	"time"
)

type Article struct {
	ID          uint
	Title       string
	Content     string
	Category    string
	CreatedDate time.Time
	UpdatedDate time.Time
	Status      string
}

type ArticleStatus string

const (
	StatusPublish ArticleStatus = "Publish"
	StatusDraft   ArticleStatus = "Draft"
	StatusTrash   ArticleStatus = "Thrash"
)

func (s ArticleStatus) IsValid() bool {
	switch s {
	case StatusPublish, StatusDraft, StatusTrash:
		return true
	default:
		return false
	}
}

func (a *Article) Validate() error {
	if a.Title == "" {
		return errors.New("title is required")
	}

	if len(a.Title) < 3 || len(a.Title) > 200 {
		return errors.New("title must be between 3 and 200 characters")
	}

	if a.Content == "" {
		return errors.New("content is required")
	}

	if len(a.Content) < 10 {
		return errors.New("content must be at least 10 characters")
	}

	if a.Category == "" {
		return errors.New("category is required")
	}

	validStatuses := map[string]bool{
		"Publish": true,
		"Draft":   true,
		"Thrash":  true,
	}
	if !validStatuses[a.Status] {
		return errors.New("invalid status, must be one of: Publish, Draft, Thrash")
	}

	return nil
}

func (a *Article) IsPublished() bool {
	return a.Status == string(StatusPublish)
}

func (a *Article) IsDraft() bool {
	return a.Status == string(StatusDraft)
}

func (a *Article) IsInTrash() bool {
	return a.Status == string(StatusTrash)
}
