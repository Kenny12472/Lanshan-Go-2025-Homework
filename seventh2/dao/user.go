package dao

import (
	"errors"
	"seventh2/model"

	"seventh2/utils"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(username, password string) error {

	var count int64
	DB.Model(&models.User{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return errors.New("用户已存在")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Username:     username,
		PasswordHash: string(hash),
	}

	return DB.Create(&user).Error
}

func LoginUser(username, password string) (string, error) {

	var user models.User

	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return "", errors.New("用户不存在")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return "", errors.New("密码错误")
	}

	token, err := utils.GenerateToken(username)
	if err != nil {
		return "", err
	}

	return token, nil
}
