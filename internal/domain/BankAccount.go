package domain

import "time"

type BankAccount struct {
	ID                uint      `json:"id" gorm:"primary_key"`
	UserID            uint      `json:"user_id"`
	BankAccountNumber uint      `json:"bank_account_number" gorm:"index;unique;not null"`
	SwiftCode         string    `json:"swift_code"`
	PaymentType       string    `json:"payment_type"`
	CreatedAt         time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
