package models

import "time"

type Transaction struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Amount      float64   `json:"amount" gorm:"type:decimal(12,2);not null"`
	Type        string    `json:"type" gorm:"size:10;not null"`
	Description string    `json:"description" gorm:"size:100"`
	Date        time.Time `json:"date" gorm:"not null"`
	UserID      uint      `json:"-" gorm:"not null"`
	User        User      `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	CategoryID  *uint     `json:"-"`
	Category    *Category `json:"category" gorm:"constraint:OnDelete:SET NULL"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}
