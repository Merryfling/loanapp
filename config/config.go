package config

import (
    "fmt"
    "log"
    "loanapp/model"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
    var err error

    // MySQL 数据库连接字符串
    dsn := "loanapp:loanapp@tcp(127.0.0.1:3306)/loanapp?charset=utf8mb4&parseTime=True&loc=Local"

    // 使用 GORM 打开 MySQL 数据库
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect to database: %v", err)
    }

    // 自动迁移模型到数据库
    if err := DB.AutoMigrate(&model.User{}, &model.Application{}); err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }

    fmt.Println("Database connected and migrated successfully!")
}
