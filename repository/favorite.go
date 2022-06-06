package repository

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Favorite struct {
	UserId     int64 `gorm:"primaryKey;autoIncrement:false"`
	VideoId    int64 `gorm:"primaryKey;autoIncrement:false"`
	IsFavorite bool  `gorm:"column:is_favorite"`
}

var favoriteDao = NewFavoriteDaoInstance()

type FavoriteDao struct {
}

func NewFavoriteDaoInstance() *FavoriteDao {
	return &FavoriteDao{}
}

func (f *FavoriteDao) Judge(uId, vId int64) bool {
	var favorite = &Favorite{
		UserId:     uId,
		VideoId:    vId,
		IsFavorite: false,
	}
	if err := db.Model(favorite).FirstOrCreate(favorite).Error; err != nil {
		fmt.Println("no record found in favorite")
	} else if favorite.IsFavorite {
		return true
	}
	return false
}

func (f *FavoriteDao) CreateFavorite(uid int64, vid int64) {
	var favorite = &Favorite{
		UserId:     uid,
		VideoId:    vid,
		IsFavorite: true,
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(favorite).Error; err != nil {
			return err
		}
		if err := tx.Model(&Video{}).
			Where("id = ?", vid).
			Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Println("transaction wrong")
		return
	}

}

func (f *FavoriteDao) CancelFavorite(uid int64, vid int64) {
	var favorite = &Favorite{
		UserId:     uid,
		VideoId:    vid,
		IsFavorite: false,
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(favorite).Error; err != nil {
			return err
		}

		if err := db.Model(&Video{}).
			Where("id = ?", vid).
			Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {

			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println("transaction fail")
		return
	}

}

func (f *FavoriteDao) FavoriteList(userId int64) []Video {

	var videos = make([]Video, 30)
	var selectedAttr = []string{"videos.id", "play_url", "cover_url", "favorite_count", "comment_count", "title",
		"favorites.is_favorite"}
	db.Model(&videos).Select(selectedAttr).Joins("left join favorites on favorites.video_id = videos.id").
		Joins("User").
		Where("favorites.user_id = ?", userId).
		Where("favorites.is_favorite = ?", true).Find(&videos)

	return videos
}
