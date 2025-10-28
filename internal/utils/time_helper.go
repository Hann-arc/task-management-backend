package utils

import (
	"time"

	"gorm.io/gorm"
)

// Helper function converts gorm.DeletedAt to *time.Time
func ToTimePtr(deletedAt gorm.DeletedAt) *time.Time {
	if deletedAt.Time.IsZero() {
		return nil
	}
	return &deletedAt.Time
}
