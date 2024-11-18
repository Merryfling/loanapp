package model

import "time"

// User 用户模型
type User struct {
    Id       uint64   `gorm:"primaryKey"`
    Phone    string `gorm:"unique"`
    Password string
    Name     string
    IdNumber string
}

// Application 贷款申请模型
type Application struct {
    Id          uint64      `gorm:"primaryKey"`
    UserId      uint64
    Income      uint64
    LoanAmount  uint64
    LoanTerm    uint64
    LoanPurpose string
    Status      string
    Score       uint64
    Remark      string // 这个记录admin想法
    Comment     string // 这个给用户
	SubmitTime  time.Time
	UpdateTime  time.Time
}