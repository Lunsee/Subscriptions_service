package models

import (
	"time"

	"github.com/google/uuid"
)

// Subscription represents a subscription in the system
// @Description Subscription model for storing subscription details
type Subscriptions struct {
	ID                  int       `json:"id" gorm:"primaryKey"`
	SubscriptionService string    `json:"subscription_service" gorm:"column:subscription_service;not null"`
	UserID              uuid.UUID `json:"user_id" gorm:"not null"`
	Price               int       `json:"price" gorm:"column:price;not null"`
	StartDate           string    `json:"start_date" gorm:"column:start_date;not null"`
	ExpDate             *string   `json:"exp_date,omitempty" gorm:"column:exp_date"`
	CreatedAt           time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt           time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}
