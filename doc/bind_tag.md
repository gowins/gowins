# Gin 框架 Struct Tag 总结

Gin 在数据绑定时支持多种 struct tag，用于控制数据绑定和验证行为。以下是完整的 tag 使用总结：

## 一、基础绑定 Tag

### 1. 数据源指定
```go
type User struct {
    ID       string `form:"id" json:"id" uri:"id" xml:"id"` // 多数据源
    Username string `form:"username"`  // 仅从表单绑定
    Email    string `json:"email"`     // 仅从JSON绑定
    Role     string `query:"role"`     // 仅从查询参数绑定
}
```

### 2. 字段映射
```go
type Product struct {
    ProductName string `json:"product_name"`  // 映射JSON中的product_name
    Price       float64 `json:"price,string"` // JSON字符串转float
}
```

## 二、验证 Tag

### 1. 必填验证
```go
type LoginForm struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}
```

### 2. 类型验证
```go
type Request struct {
    Age      int     `binding:"numeric"`   // 必须为数字
    Email    string  `binding:"email"`     // 必须为邮箱格式
    URL      string  `binding:"url"`       // 必须为URL
    UUID     string  `binding:"uuid"`      // 必须为UUID
    IP       string  `binding:"ip"`        // 必须为IP地址
    Phone    string  `binding:"e164"`      // E.164格式电话号码
}
```

### 3. 范围验证
```go
type Config struct {
    Port     int    `binding:"min=1024,max=65535"`  // 端口范围
    Discount float64 `binding:"gte=0,lte=1"`        // 0-1之间
    Age      int    `binding:"oneof=18 20 25"`     // 只能是给定值之一
}
```

### 4. 字符串验证
```go
type User struct {
    Username string `binding:"alphanum"`        // 只允许字母数字
    Password string `binding:"contains=!@#?"`   // 必须包含特殊字符
    Code     string `binding:"len=6"`           // 固定长度6
    Bio      string `binding:"max=500"`         // 最大长度500
}
```

### 5. 时间验证
```go
type Event struct {
    Start time.Time `binding:"required" time_format:"2006-01-02"`
    End   time.Time `binding:"required,gtfield=Start" time_format:"2006-01-02"`
}
```

## 三、特殊控制 Tag

### 1. 忽略字段
```go
type User struct {
    ID       string `json:"-"`                  // 完全忽略该字段
    Password string `json:"password,omitempty"` // 空值时忽略
}
```

### 2. 嵌套结构
```go
type Order struct {
    User    User   `json:"user"`                // 嵌套结构
    Items   []Item `json:"items" binding:"dive"` // 对数组元素验证
}

type Item struct {
    ID    string `json:"id" binding:"required"`
    Count int    `json:"count" binding:"min=1"`
}
```

### 3. 自定义错误消息
```go
type RegisterForm struct {
    Email string `binding:"required,email" msg:"邮箱不能为空且格式必须正确"`
}
```

## 四、高级 Tag 用法

### 1. 条件绑定
```go
type Request struct {
    Type     string `json:"type"`
    // 当Type为"vip"时要求VIPCode必填
    VIPCode  string `json:"vip_code" binding:"required_if=Type vip"`
}
```

### 2. 跨字段验证
```go
type ChangePassword struct {
    NewPassword     string `binding:"required"`
    ConfirmPassword string `binding:"required,eqfield=NewPassword"`
}

type DateRange struct {
    Start time.Time `time_format:"2006-01-02"`
    End   time.Time `time_format:"2006-01-02" binding:"gtefield=Start"`
}
```

### 3. 自定义验证器
```go
// 注册自定义验证
if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
    v.RegisterValidation("bookabledate", func(fl validator.FieldLevel) bool {
        date, ok := fl.Field().Interface().(time.Time)
        return ok && date.After(time.Now())
    })
}

// 使用
type Booking struct {
    CheckIn time.Time `binding:"required,bookabledate" time_format:"2006-01-02"`
}
```

## 五、完整示例

```go
type UserCreateRequest struct {
    // 基本信息
    Username  string `json:"username" binding:"required,alphanum,min=3,max=20"`
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=8,containsany=!@#$%"`

    // 个人信息
    Age       int    `json:"age" binding:"omitempty,min=18,max=100"`
    Gender    string `json:"gender" binding:"omitempty,oneof=male female other"`

    // 地址信息
    Address   struct {
        City    string `json:"city" binding:"required_with=Address"`
        ZipCode string `json:"zip_code" binding:"required_with=Address,len=6"`
    } `json:"address" binding:"dive"`

    // 系统信息
    SignupIP  string `json:"-"` // 不绑定客户端数据
    CreatedAt time.Time `json:"created_at" binding:"-"` // 忽略绑定
}
```

## 六、注意事项

1. **Tag 优先级**：当同一个字段有多个绑定时（如同时有json和form），Gin会按照以下顺序尝试：
    - 显式指定的绑定方法（如ShouldBindJSON）优先
    - URI > Query > Form > JSON > XML (对于ShouldBind)

2. **零值问题**：使用指针类型区分"未设置"和"零值"
   ```go
   type Config struct {
       Enable *bool `json:"enable"` // 可以区分false和未设置
   }
   ```

3. **性能考虑**：复杂的验证逻辑会影响性能，对于高频API应简化验证规则

4. **错误消息**：默认错误消息不友好，建议使用自定义错误处理：
   ```go
   if err := c.ShouldBind(&input); err != nil {
       if verr, ok := err.(validator.ValidationErrors); ok {
           for _, fe := range verr {
               c.JSON(400, gin.H{
                   "field": fe.Field(),
                   "error": fe.Tag(), // 可以映射为友好消息
               })
               return
           }
       }
   }
   ```

