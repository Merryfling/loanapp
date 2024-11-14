package model

import "time"

// User 用户模型
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Phone    string `gorm:"unique"`
    Password string
    Name     string
    IDNumber string
}

// LoanApplication 贷款申请模型
type LoanApplication struct {
    ID          uint      `gorm:"primaryKey"`
    UserID      uint
    Income      float64
    LoanAmount  float64
    LoanTerm    int
    LoanPurpose string
    Status      string
    Score       int
    CreatedAt   time.Time
}
