package repository

import (
	"errors"
	"gins/data/request"
	"gins/helper"
	"gins/model"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	Db *gorm.DB
}

func NewUsersRepositoryImpl(Db *gorm.DB) UsersRepository {
	return &UserRepositoryImpl{Db: Db}
}

func (u *UserRepositoryImpl) Delete(usersId int) {
	var users model.Users
	result := u.Db.Where("id = ?", usersId).Delete(&users)
	helper.ErrorPanic(result.Error)
}

func (u *UserRepositoryImpl) FindAll() []model.Users {
	var users []model.Users
	result := u.Db.Find(&users)
	helper.ErrorPanic(result.Error)
	return users
}

func (u *UserRepositoryImpl) FindById(usersId int) (model.Users, error) {
	var users model.Users
	result := u.Db.Find(&users, usersId)
	if result != nil {
		return users, nil
	} else {
		return users, errors.New("user not found")
	}
}

func (u *UserRepositoryImpl) FindByUsername(username string) (model.Users, error) {
	var users model.Users
	result := u.Db.Find(&users, "username = ?", username)
	if result.Error != nil {
		return users, errors.New("invalid username or password")
	}
	return users, nil
}

func (u *UserRepositoryImpl) Save(users model.Users) {
	result := u.Db.Create(&users)
	helper.ErrorPanic(result.Error)
}

func (u *UserRepositoryImpl) Update(users model.Users) {
	var updateUsers = request.UpdateUsersRequest{
		Id:       users.Id,
		Username: users.Username,
		Email:    users.Email,
		Password: users.Password,
	}
	result := u.Db.Model(&users).Updates(updateUsers)
	helper.ErrorPanic(result.Error)
}
