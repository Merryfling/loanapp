package handler

import (
    "loanapp/config"
    "loanapp/model"
    "loanapp/api"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
)

// 创建贷款申请
func ApplyLoan(c *gin.Context) {
    // 从 Context 中获取 userId
    userIdAny, exists := c.Get("userId")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
        return
    }
    // 类型断言为 uint64，不允许一次性处理，应为有两个返回
    userId, ok := userIdAny.(uint64)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
        return
    }
    var req api.ApplicationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    application := model.Application{
        UserId: userId,
        Income: req.Income,
        LoanAmount: req.LoanAmount,
        LoanTerm: req.LoanTerm,
        LoanPurpose: req.LoanPurpose,
        Status: "待审核",
        // Score: TBD
        Remark: "等待审核",
        Comment: "正在审核中",
        SubmitTime: time.Now(),
        UpdateTime: time.Now(),
    }
    if err := config.DB.Create(&application).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create loan application"})
        return
    }

    data := &api.ApplicationData{
        ApplicationId: application.Id,
        ApplicationStatus: application.Status,
        Comment: application.Comment,
        // Score: application.score,
    }

    c.JSON(http.StatusOK, api.ApplicationResponse{
        Status:         "success",
        Message:        "Loan Application submitted successfully",
        Data: data,
    })
}

// 查询贷款状态
func ApplicationStatus(c *gin.Context) {
    // 从 Context 中获取 userID
    userId, exists := c.Get("userId")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
        return
    }

    applicationId := c.Param("application_id")
    var application model.Application
    if err := config.DB.Where("id = ?", applicationId).First(&application).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Loan application not found"})
        return
    }
    
    if userId != application.UserId {
        c.JSON(http.StatusNotFound, gin.H{"error": "Invalid identification"})
        return
    }

    data := &api.ApplicationData{
        ApplicationId: application.Id,
        ApplicationStatus: application.Status,
        Comment: application.Comment,
        // Score: application.score,
    }

    c.JSON(http.StatusOK, api.ApplicationStatusResponse{
        Status:         "success",
        Message:        "Loan application status retrieved",
        Data: data,
    })    
}