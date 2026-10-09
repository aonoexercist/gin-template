package user

import "gorm.io/gorm"

// NewModule wires the feature's repository, service and handler.
func NewModule(db *gorm.DB) *Handler {
	return NewHandler(NewService(NewRepository(db)))
}
