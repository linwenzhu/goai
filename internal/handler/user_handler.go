package handler

import (
	"github.com/gin-gonic/gin"
	"goAi/internal/service"
	"goAi/pkg/response"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{svc: service.NewUserService()}
}

// RegisterRequest 注册请求结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"required,email"`
}

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误："+err.Error())
		return
	}
	if err := h.svc.Register(req.Username, req.Password, req.Email); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误："+err.Error())
	}

	token, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Success(c, gin.H{"token": token})

}
