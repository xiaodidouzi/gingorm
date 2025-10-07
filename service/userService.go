package service

import (
	"errors"
	"gingorm/global"
	"gingorm/models"
	"gingorm/service/dto"
	"gingorm/utils"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type UserService struct {
	DB  *gorm.DB
	Rdb *redis.Client
}

// 构造函数
func NewUserService(db *gorm.DB, rdb *redis.Client) *UserService {
	return &UserService{
		DB:  db,
		Rdb: rdb,
	}
}

// 注册
func (s *UserService) Register(req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	//检查用户名是否已存在
	var existUser models.User
	if err := global.DB.Where("username=?", req.Username).First(&existUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	//哈希密码
	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	//保存到数据库
	user := models.User{
		Username: req.Username,
		Password: hashPassword,
	}
	if err := global.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	//返回结果
	res := &dto.RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
	}
	return res, nil
}

// 登录
func (s *UserService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	//检查用户名是否存在
	var user models.User
	if err := global.DB.Where("username=?", req.Username).First(&user).Error; err != nil {
		return nil, errors.New("用户名不存在")
	}
	//验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("密码错误")
	}
	//生成JWT
	token, err := utils.GenerateJWT(int(user.ID), user.Username)
	if err != nil {
		return nil, err
	}
	//返回结果
	res := &dto.LoginResponse{
		Username: user.Username,
		Token:    token,
	}
	return res, nil
}
