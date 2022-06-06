package repository

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserFollow struct {
	UserID     int64 `gorm:"primaryKey;autoIncrement:false"`
	FollowedID int64 `gorm:"primaryKey;autoIncrement:false"`
	IsFollow   bool
}

type RelationDao struct {
}

func NewRelationDaoInstance() *RelationDao {
	return &RelationDao{}
}

func (d *RelationDao) Follow(from_user_id int64, to_user_id int64) {
	var r = &UserFollow{
		UserID:     from_user_id,
		FollowedID: to_user_id,
		IsFollow:   true,
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(r).Error; err != nil {
			return err
		}
		if err := tx.Model(&User{ID: from_user_id}).
			Update("follow_count", gorm.Expr("follow_count + ?", 1)).Error; err != nil {
			return err
		}
		if err := tx.Model(&User{ID: to_user_id}).
			Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println("transaction fail")
		return
	}
	return
}

func (d *RelationDao) UnFollow(from_user_id int64, to_user_id int64) {
	var r = &UserFollow{
		UserID:     from_user_id,
		FollowedID: to_user_id,
		IsFollow:   false,
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(r).Error; err != nil {
			return err
		}
		if err := tx.Model(&User{ID: from_user_id}).
			Update("follow_count", gorm.Expr("follow_count - ?", 1)).Error; err != nil {
			return err
		}
		if err := tx.Model(&User{ID: to_user_id}).
			Update("follower_count", gorm.Expr("follow_count - ?", 1)).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println("transaction fail")
		return
	}
	return
}

func (d *RelationDao) FollowListOf(userID int64) *[]User {
	//找到 userID 关注的人
	var users = make([]User, 30)
	db.Model(&users).Joins("left join user_follows on users.id = user_follows.followed_id").
		Where("user_follows.user_id = ?", userID).
		Where("user_follows.is_follow = ?", true).
		Find(&users)

	for i, _ := range users {
		users[i].IsFollow = true
	}
	return &users
}

func (d *RelationDao) FansOf(userID int64) *[]User {
	var users = make([]User, 30)
	db.Model(&users).Joins("left join user_follows on users.id = user_follows.user_id").
		Where("user_follows.followed_id = ?", userID).
		Find(&users)

	//装填是否关注
	for i, _ := range users {
		users[i].IsFollow = Isfollow(users[i].ID, userID)
	}
	return &users
}

func Isfollow(fromId, toId int64) bool {
	var f = &UserFollow{
		UserID:     fromId,
		FollowedID: toId,
		IsFollow:   false,
	}
	db.Model(f).FirstOrCreate(f)
	if f.IsFollow {
		return true
	} else {
		return false
	}
	return false
}
