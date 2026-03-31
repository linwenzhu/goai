package service

import (
	"errors"
	"goAi/internal/model"
	"goAi/internal/repository"
	"goAi/pkg/jwt"
	"goAi/pkg/util"
	"gorm.io/gorm"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{repo: repository.NewUserRepository()}
}

func (s *UserService) Register(username, password, email string) error {
	// 1. 检查用户名是否已存在
	_, err := s.repo.FindByUsername(username)
	if err == nil {
		return errors.New("用户名已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	// 2. 检查邮箱是否已存在
	_, err = s.repo.FindByEmail(email)
	if err == nil {
		return errors.New("邮箱已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 3. 加密密码
	hashed, err := util.HashPassword(password)
	if err != nil {
		return err
	}

	// 4. 创建用户
	user := &model.User{
		Username: username,
		Password: hashed,
		Email:    email,
	}
	return s.repo.Create(user)

}

func (s *UserService) Login(username string, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", errors.New("用户名或密码错误")
	}
	if err != nil {
		return "", err
	}

	if !util.CheckPassword(password, user.Password) {
		return "", errors.New("用户名或密码错误")
	}

	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}
