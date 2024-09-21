package models

import "time"

type Wallet struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`                                // Primary key, auto-increment
	Address   string    `json:"address" gorm:"size:42;not null;unique"`                            // Ethereum address (unique)
	Amount    float64   `json:"amount" gorm:"type:decimal(20,8);not null"`                         // Monetary amount
	Timestamp time.Time `json:"timestamp" gorm:"type:datetime;not null;default:current_timestamp"` // Transaction or record timestamp
}

type Prize struct {
	Address string  `json:address`
	Prize   float64 `json:prize`
}

type PrizeList struct {
	PrizeName   string    `json:"prize_name" gorm:"size:255;not null"`       // Name of the prize
	Amount      float64   `json:"amount" gorm:"type:decimal(20,8);not null"` // Monetary amount for the prize
	Probability int       `json:"probability" gorm:"not null"`               // Probability of winning the prize
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`          // Automatically set at creation
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`          // Automatically set at update
}
