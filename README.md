# APP后端接口

## API接口
数据交互均为`JSON`格式

### USER

#### **1. 注册接口**
`{{base_url}}/api/v1/register`
- **接口说明**：用户使用手机号注册。
- **请求方式**：`POST`
- **接口地址**：`/api/v1/register`
- **请求参数**：
  - `phone`（String，必填）：用户手机号。
  - `password`（String，必填）：用户密码。
  - `captcha`（String，选填）：验证码。等待后续的验证码接入
- 请求示例：
 ```JSON
 {
  "phone": "222",
  "password": "222",
  "captcha": "456"
 }
 ```
- **返回示例**：
  ```json
{
	"status": "success",
	"message": "User registered successfully",
	"data": {
		"user_id": 2,
		"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzE5NDc5NDgsInVzZXJfaWQiOjJ9.I4B-5GJqmmwvBaKqQmVCzDjJYXp1u8mOe3umYIT8X34"
	}
}
  ```

---

#### **2. 登录接口**
`{{base_url}}/api/v1/login`
- **接口说明**：用户使用手机号和密码登录。
- **请求方式**：`POST`
- **接口地址**：`/api/v1/login`
- **请求参数**：
  - `phone`（String，必填）：用户手机号。
  - `password`（String，必填）：用户密码。
- 请求示例：

```JSON
{
    "phone": "222",
    "password": "222"
}
```

- **返回示例**：

```json
{
    "status": "success",
    "message": "Login successful",
    "data": {
        "user_id": 2,
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzE5NDc5NTYsInVzZXJfaWQiOjJ9.E3I0d7YtRr41epU649jcKLLnA21hzWsUnZKHTcSS7hk"
    }
}
```

---

### LOAN

下列接口均需要登录时得到的`token`作为验证
#### **3. 创建贷款申请**
`{{base_url}}/api/v1/loan/apply`
- **接口说明**：用户提交贷款申请。
- **请求方式**：`POST`
- **接口地址**：`/api/v1/loan/apply`
- **请求头**：`Authorization: Bearer {token}`
- **请求参数**：
  - `name`（String，必填）：用户姓名。
  - `id_number`（String，必填）：身份证号。实际生产请对此进行验证
  - `income`（Float，必填）：月收入。
  - `loan_amount`（Float，必填）：申请贷款金额。
  - `loan_term`（Integer，必填）：贷款期限（单位：月）。
  - `loan_purpose`（String，选填）：贷款用途。
- 请求示例：

```JSON
{
    "name": "phone222",
    "id_number": "phone222'sID",
    "income": 100000,
    "loan_amount": 10000,
    "loan_term": 2,
    "loan_purpose": "test:first application"
}
```

- **返回示例**：

```json
{
    "status": "success",
    "message": "Loan Application submitted successfully",
    "data": {
        "application_id": 1,
        "application_status": "待审核",
        "comment": "正在审核中"
    }
}
```

---

#### **4. 查询贷款申请状态**
`{{base_url}}/api/v1/loan/status/:application_id`
- **接口说明**：用户查询贷款申请的状态。
- **请求方式**：`GET`
- **接口地址**：`/api/v1/loan/status/{application_id}`
- **请求头**：`Authorization: Bearer {token}`
- **请求参数**：无

- **返回示例**：

```json
{
    "status": "success",
    "message": "Loan application status retrieved",
    "data": {
        "application_id": 1,
        "application_status": "待审核",
        "comment": "正在审核中"
    }
}
```

---

### PROFILE

#### **5. 查看历史贷款记录**
`{{base_url}}/api/v1/user/history`
- **接口说明**：用户查看所有历史贷款申请记录。
- **请求方式**：`GET`
- **接口地址**：`/api/v1/user/history`
- **请求头**：`Authorization: Bearer {token}`
- **请求参数**：无

- **返回示例**：

```json
{
    "status": "success",
    "message": "Get Application History successful",
    "data": [
        {
            "application_id": 1,
            "application_status": "待审核",
            "loan_amount": 10000,
            "loan_term": 2,
            "submission_time": "2024-11-18 23:46",
            "comment": "正在审核中"
        },
        {
            "application_id": 2,
            "application_status": "待审核",
            "loan_amount": 10000,
            "loan_term": 2,
            "submission_time": "2024-11-18 23:49",
            "comment": "正在审核中"
        }
    ]
}
```

---

#### **6. 查看用户资料**
`{{base_url}}/api/v1/user/profile`
- **接口说明**：用户查看自己的资料信息。
- **请求方式**：`GET`
- **接口地址**：`/api/v1/user/profile`
- **请求头**：`Authorization: Bearer {token}`
- **请求参数**：无

- **返回示例**：

  ```json
  {
    "status": "success",
    "message": "User profile retrieved",
    "data": {
      "user_id": "12345",
      "name": "John Doe",
      "phone": "1234567890",
      "id_number": "123456789012345678"
    }
  }
  ```

---

#### **7. 更新用户资料**
`{{base_url}}/api/v1/user/profile`
- **接口说明**：用户可以更新自己的资料。
- **请求方式**：`POST`
- **接口地址**：`/api/v1/user/profile`
- **请求头**：`Authorization: Bearer {token}`
- **请求参数**：
  - `name`（String，选填）：用户姓名。
  - `phone`（String，选填）：手机号。
  - `password`（String，选填）：密码。
  - `id_number`（String，选填）：身份证号。请在发送时进行检验
- 请求示例：

```JSON
{
    "name": "zhang san",
    "phone": "111",
    "password": "111",
    "id_number": "111111111111111111"
}
```

- **返回示例**：

```json
{
    "status": "success",
    "message": "GetProfile successful",
    "data": {
        "user_id": 1,
        "name": "zhang san",
        "phone": "111",
        "id_number": "111************111"
    }
}
```
