package dto

import "errors"

var (
	// Validation errors
	ErrTitleRequired       = errors.New("title is required")
	ErrTitleLength         = errors.New("title must be between 3 and 200 characters")
	ErrContentRequired     = errors.New("content is required")
	ErrContentTooShort     = errors.New("content must be at least 10 characters")
	ErrCategoryRequired    = errors.New("category is required")
	ErrInvalidStatus       = errors.New("invalid status, must be one of: Publish, Draft, Thrash")
	ErrInvalidFilterStatus = errors.New("invalid filter status")

	// Database errors
	ErrArticleNotFound = errors.New("article not found")
	ErrArticleExists   = errors.New("article already exists")

	// Business logic errors
	ErrFailedCreateArticle = errors.New("failed to create article")
	ErrFailedUpdateArticle = errors.New("failed to update article")
	ErrFailedDeleteArticle = errors.New("failed to delete article")
)

type ErrorCode string

const (
	// Validation error codes
	ErrCodeValidation       ErrorCode = "VALIDATION_ERROR"
	ErrCodeTitleRequired    ErrorCode = "TITLE_REQUIRED"
	ErrCodeTitleInvalid     ErrorCode = "TITLE_INVALID"
	ErrCodeContentRequired  ErrorCode = "CONTENT_REQUIRED"
	ErrCodeContentInvalid   ErrorCode = "CONTENT_INVALID"
	ErrCodeCategoryRequired ErrorCode = "CATEGORY_REQUIRED"
	ErrCodeStatusInvalid    ErrorCode = "STATUS_INVALID"

	// Database error codes
	ErrCodeNotFound ErrorCode = "NOT_FOUND"
	ErrCodeConflict ErrorCode = "CONFLICT"
	ErrCodeDBError  ErrorCode = "DATABASE_ERROR"

	// Business logic error codes
	ErrCodeCreateFailed ErrorCode = "CREATE_FAILED"
	ErrCodeUpdateFailed ErrorCode = "UPDATE_FAILED"
	ErrCodeDeleteFailed ErrorCode = "DELETE_FAILED"

	// General error codes
	ErrCodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden     ErrorCode = "FORBIDDEN"
	ErrCodeInternalError ErrorCode = "INTERNAL_SERVER_ERROR"
)

func MapErrorToCode(err error) ErrorCode {
	if err == nil {
		return ""
	}

	switch err {
	case ErrTitleRequired:
		return ErrCodeTitleRequired
	case ErrTitleLength:
		return ErrCodeTitleInvalid
	case ErrContentRequired:
		return ErrCodeContentRequired
	case ErrContentTooShort:
		return ErrCodeContentInvalid
	case ErrCategoryRequired:
		return ErrCodeCategoryRequired
	case ErrInvalidStatus, ErrInvalidFilterStatus:
		return ErrCodeStatusInvalid
	case ErrArticleNotFound:
		return ErrCodeNotFound
	case ErrArticleExists:
		return ErrCodeConflict
	case ErrFailedCreateArticle:
		return ErrCodeCreateFailed
	case ErrFailedUpdateArticle:
		return ErrCodeUpdateFailed
	case ErrFailedDeleteArticle:
		return ErrCodeDeleteFailed
	default:
		return ErrCodeInternalError
	}
}
