package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Transaction object
type Transaction struct {
	gorm.Model
	UserID      uint      `gorm:"not null" json:"user_id"`
	Name        string    `gorm:"not null" json:"name"`
	Amount      float64   `gorm:"not null" json:"amount"`
	TxnType     string    `gorm:"type:enum('income','expense');not null" json:"txn_type"`
	Frequency   string    `gorm:"type:enum('daily','weekly','monthly','quarterly','yearly');not null" json:"frequency"`
	CategoryID  uint      `json:"category_id"`
	TxnDate     time.Time `gorm:"not null" json:"txn_date"`
	Description string    `json:"description"`
}

// Category object
type Category struct {
	gorm.Model
	Name string `gorm:"not null" json:"name"`
	Type string `gorm:"type:enum('income','expense');not null" json:"type"`
}

// Metrics object

type DashboardMetrics struct {
	UserID               uint           `gorm:"primaryKey" json:"user_id"`
	TotalIncome          float64        `json:"total_income"`
	TotalExpense         float64        `json:"total_expense"`
	NetSavings           float64        `json:"net_savings"`
	MonthlyTotals        datatypes.JSON `json:"monthly_totals"`         // JSONB for monthly line chart
	TopExpenseCategories datatypes.JSON `json:"top_expense_categories"` // JSONB for pie chart
	UpdatedAt            time.Time      `json:"updated_at"`
}
