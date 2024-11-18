package handler

import (
    "loanapp/config"
    "loanapp/model"
    "loanapp/api"
    "net/http"
    "fmt"
    "github.com/gin-gonic/gin"
)

// 用户注册
func Register(c *gin.Context) {
    var req api.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    user := model.User{
        Phone:    req.Phone,
        Password: req.Password,  // 实际开发中要对密码加密
    }
    // 先检查phone存不存在，创建哪怕失败，id都已经自增
    var existingUser model.User
    err := config.DB.Where("phone = ?", user.Phone).First(&existingUser).Error
    if err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already registered"})
        return
    }

    if err := config.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    // 生成 token
    token, err := config.GenerateToken(user.Id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    data := &api.RegisterData{
        UserId: user.Id,
        Token: token,
    }

    c.JSON(http.StatusOK, api.RegisterResponse{
        Status:  "success",
        Message: "User registered successfully",
        Data: data,
    })
}

// 用户登录
func Login(c *gin.Context) {
    var req api.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    var user model.User
    if err := config.DB.Where("phone = ? AND password = ?", req.Phone, req.Password).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // 生成 token
    token, err := config.GenerateToken(user.Id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    data := &api.LoginData{
        UserId: user.Id,
        Token: token,
    }

    c.JSON(http.StatusOK, api.LoginResponse{
        Status:  "success",
        Message: "Login successful",
        Data: data,
    })
}

// 用户资料查询
func GetProfile(c *gin.Context) {
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

    var user model.User
    // 根据 userID 查询用户信息
    if err := config.DB.Where("id = ?", userId).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    data := &api.UserInfo{
        UserId: user.Id,
        Name: user.Name,
        Phone: user.Phone,
        IdNumber: user.IdNumber[:3]+"************"+user.IdNumber[len(user.IdNumber)-3:],
    }
    // [start:end] 左闭右开，索引号base 0
    c.JSON(http.StatusOK, api.GetUserProfileResponse{
        Status:  "success",
        Message: "GetProfile successful",
        Data: data,
    })
}

// 用户资料更新
func UpdateProfile(c *gin.Context) {
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

    var req api.UpdateUserProfileRequest
    // 解析 JSON 请求
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    var user model.User
    // 查询用户信息
    if err := config.DB.Where("id = ?", userId).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // 更新用户信息
    if req.Name != "" {
        user.Name = req.Name
    }
    if req.Phone != "" {
        user.Phone = req.Phone
        var existingPhone User
        err := db.Where("phone = ?", req.Phone).First(&existingPhone).Error
        if err == nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "This phone number is already used"})
            return
        }
    }
    if req.Password != "" {
        user.Password = req.Password
    }
    if req.IdNumber != "" {
        user.IdNumber = req.IdNumber
    }

    // 保存更新
    if err := config.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
        return
    }

    // 返回成功响应
    // 省略了userinfo，即data确实不需要
    c.JSON(http.StatusOK, api.UpdateUserProfileResponse{
        Status:  "success",
        Message: "User profile updated successfully",
    })
}

func ApplicationHistory(c *gin.Context) {
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

    var applications []model.Application
    // 根据 userId 查询用户信息
    if err := config.DB.Where("user_id = ?", userId).Find(&applications).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User has zero application"})
        return
    }
    fmt.Println("application:", len(applications))
    if len(applications) == 0 {
        c.JSON(http.StatusOK, api.ApplicationHistoryResponse{
            Status:  "success",
            Message: "Get Application History successful",
            Data:    []*api.ApplicationData{}, // 空数组
        })
        return
    }

    var data []*api.ApplicationData
    for _, application := range applications {
        data = append(data, &api.ApplicationData{
            ApplicationId: application.Id,
            ApplicationStatus: application.Status,
            LoanAmount: application.LoanAmount,
            LoanTerm: application.LoanTerm,
            SubmissionTime: application.SubmitTime.Format("2006-01-02 15:04"),
            Comment: application.Comment,
            Score: application.Score,
        })
    }

    c.JSON(http.StatusOK, api.ApplicationHistoryResponse{
        Status:  "success",
        Message: "Get Application History successful",
        Data: data,
    })
}