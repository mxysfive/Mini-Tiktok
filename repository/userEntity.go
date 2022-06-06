package repository

import "fmt"

type User struct {
	ID            int64  `gorm:"primary_key" json:"id"`
	Name          string `gorm:"column:name;size:32;not null" json:"name"`
	Password      string `gorm:"column:password;size:32;" json:"password"`
	FollowCount   int64  `gorm:"column:follow_count" json:"follow_count"`
	FollowerCount int64  `gorm:"column:follower_count" json:"follower_count"`
	IsFollow      bool   `gorm:"column:is_follow" json:"is_follow"`
}

type UserDao struct {
}

func NewUserDaoInstance() *UserDao {
	return &UserDao{}
}

func (u *UserDao) QueryUserById(userId int64) (*User, error) {
	var user = &User{}
	err := db.First(user, userId).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserDao) CreateByNameAndPassword(name, password string) *User {
	var user = &User{
		ID:            0,
		Name:          name,
		Password:      password,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}
	db.Create(user)
	fmt.Println(user.ID)
	return user
}

func (u *UserDao) QueryLoginInfo(name, password string) *User {
	var user = &User{}
	db.Where("name = ? and password = ?", name, password).Find(user)
	return user
}
