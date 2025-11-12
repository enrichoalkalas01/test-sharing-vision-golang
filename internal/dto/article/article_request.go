package dto

type CreateArticleRequest struct {
	Title    string `json:"title" validate:"required,min=3,max=200"`
	Content  string `json:"content" validate:"required,min=10"`
	Category string `json:"category" validate:"required,min=3,max=100"`
	Status   string `json:"status" validate:"required,oneof=Publish Draft Thrash"`
}

type UpdateArticleRequest struct {
	Title    string `json:"title" validate:"omitempty,min=3,max=200"`
	Content  string `json:"content" validate:"omitempty,min=10"`
	Category string `json:"category" validate:"omitempty,min=3,max=100"`
	Status   string `json:"status" validate:"omitempty,oneof=Publish Draft Thrash"`
}

func (r *CreateArticleRequest) Validate() error {
	if r.Title == "" {
		return ErrTitleRequired
	}
	if len(r.Title) < 3 || len(r.Title) > 200 {
		return ErrTitleLength
	}

	if r.Content == "" {
		return ErrContentRequired
	}
	if len(r.Content) < 10 {
		return ErrContentTooShort
	}

	if r.Category == "" {
		return ErrCategoryRequired
	}

	validStatuses := map[string]bool{
		"Publish": true,
		"Draft":   true,
		"Thrash":  true,
	}
	if !validStatuses[r.Status] {
		return ErrInvalidStatus
	}

	return nil
}

func (r *UpdateArticleRequest) Validate() error {
	if r.Title != "" {
		if len(r.Title) < 3 || len(r.Title) > 200 {
			return ErrTitleLength
		}
	}

	if r.Content != "" {
		if len(r.Content) < 10 {
			return ErrContentTooShort
		}
	}

	if r.Status != "" {
		validStatuses := map[string]bool{
			"Publish": true,
			"Draft":   true,
			"Thrash":  true,
		}
		if !validStatuses[r.Status] {
			return ErrInvalidStatus
		}
	}

	return nil
}
