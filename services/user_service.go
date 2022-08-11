package services

import (
	"github.com/kataras/iris/v12/core/errgroup"
	"golang.org/x/crypto/bcrypt"
	"imooc-Product/datamodels"
	"imooc-Product/repositories"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool)
	AddUser(user *datamodels.User) (userId int64, err error)
}

type UserService struct {
	UserRepository repositories.IUser
}

func NewUserService(repository repositories.IUser) IUserService {
	return &UserService{repository}
}
func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool) {

	user, err := u.UserRepository.Select(userName)
	if err != nil {
		return
	}
	isOk, _ = ValidatePassword(pwd, user.HashPassword)

	if !isOk {
		return &datamodels.User{}, false
	}

	return
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error) {
	pwdByte, errPwd := GeneratePassword(user.HashPassword)
	if errPwd != nil {
		return userId, errPwd
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

/*
	GenerateFromPassword 以给定的代价返回密码的 bcrypt 哈希值。
	如果给定的成本小于 MinCost，则成本将设置为 DefaultCost。
	使用此包中定义的 CompareHashAndPassword 将返回的散列密码与其明文版本进行比较。

	密码数据库加密
*/
func GeneratePassword(userPassword string) ([]byte, error) {

	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

// ValidatePassword 查询密码是否正确
func ValidatePassword(userPassword string, hashed string) (isOK bool, err error) {
	//解密数据库密码 并比较
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errgroup.New("密码比对错误！")
	}
	return true, nil

}
