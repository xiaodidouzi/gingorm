package dto

// RegisterRequest 注册请求 DTO
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRequest 登录请求 DTO
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterResponse 注册响应 DTO
type RegisterResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

// LoginResponse 登录响应 DTO
type LoginResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
