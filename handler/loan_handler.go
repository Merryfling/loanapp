package handler

import (
    "loanapp/config"
    "loanapp/model"
    "loanapp/api"
    "net/http"
    "fmt"
    "github.com/gin-gonic/gin"
)

// 创建贷款申请
func ApplyLoan(c *gin.Context) {
    var req api.LoanApplicationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    application := model.LoanApplication{
        UserID:      1,                          // 示例值
        Income:      float64(req.Income),        // 类型转换
        LoanAmount:  float64(req.LoanAmount),    // 类型转换
        LoanTerm:    int(req.LoanTerm),          // 类型转换
        LoanPurpose: req.LoanPurpose,
        Status:      "Pending",
    }
    if err := config.DB.Create(&application).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create loan application"})
        return
    }

    c.JSON(http.StatusOK, api.LoanApplicationResponse{
        Status:         "success",
        Message:        "Loan application submitted",
        ApplicationId:  fmt.Sprintf("%d", application.ID),  // 转换为字符串
        ApplicationStatus: application.Status,
    })
}

// 查询贷款状态
func GetLoanStatus(c *gin.Context) {
    // 从 Context 中获取 userID
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
        return
    }
    // TODO: 修改proto，get不需要参数query/form
    // TODO: 金额修改成整数配合固定精度

    applicationID := c.Param("application_id")
    var application model.LoanApplication
    if err := config.DB.Where("id = ?", applicationID).First(&application).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Loan application not found"})
        return
    }
    
    if userID != application.UserID {
        c.JSON(http.StatusNotFound, gin.H{"error": "Invalid identification"})
        return
    }

    c.JSON(http.StatusOK, api.LoanStatusResponse{
        Status:         "success",
        Message:        "Loan application status retrieved",
        ApplicationId:  fmt.Sprintf("%d", application.ID),  // 转换为字符串
        LoanStatus:     application.Status,
        Score:          int32(application.Score),           // 转换为 int32
        Comments:       "Loan application is under review",
    })    
}

func GetLoanHistory(c *gin.Context) {
    // 从上下文中获取用户 ID（通过中间件设置）
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
        return
    }

    // 查询用户的贷款记录
    var loans []model.LoanApplication
    if err := config.DB.Where("user_id = ?", userID).Find(&loans).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch loan history"})
        return
    }

    // 构造响应数据
    var records []*api.LoanHistoryRecord
    for _, loan := range loans {
        record := &api.LoanHistoryRecord{
            ApplicationId:   fmt.Sprintf("%d", loan.ID), // 转为字符串
            LoanStatus:      loan.Status,
            LoanAmount:      float32(loan.LoanAmount),
            LoanTerm:        int32(loan.LoanTerm),
            SubmissionDate:  loan.CreatedAt.Format("2006-01-02"), // 格式化日期
        }
        records = append(records, record)
    }

    // 返回响应
    c.JSON(http.StatusOK, &api.LoanHistoryResponse{
        Records: records,
    })
}