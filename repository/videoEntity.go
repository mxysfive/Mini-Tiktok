package repository

import (
	"errors"
	"fmt"
	"time"
)

const MaxListLength = 30

type Video struct {
	ID            int64  `gorm:"primaryKey" json:"id"`
	PlayUrl       string `gorm:"size:64" json:"play_url"`
	CoverUrl      string `gorm:"size:64" json:"cover_url"`
	FavoriteCount int64  `gorm:"column:favorite_count" json:"favorite_count"`
	CommentCount  int64  `gorm:"column:comment_count" json:"comment_count"`
	Title         string `gorm:"column:title;size:32" json:"title"`
	IsFavorite    bool   `gorm:"column:is_favorite" json:"is_favorite"`
	UserID        int64  `json:"user_id"`
	User          User   `json:"author"`
	CreateTime    int64  `json:"create_time"`
}

type VideoDao struct {
}

func NewVideoDaoInstance() *VideoDao {
	return &VideoDao{}
}

func (d *VideoDao) QueryByOwner(ownerId int64) []Video {
	//在用户查看自己的发布视频时使用，feed接口不用这个

	var videos = make([]Video, 30)
	db.Where("user_id = ?", ownerId).Preload("User").Find(&videos)

	packFavoriteIn(&videos, ownerId)
	return videos
}

func (d *VideoDao) CreateVideoRecord(userId int64, playURL string, coverURL string, title string) error {
	var video = &Video{
		ID:            0,
		PlayUrl:       playURL,
		CoverUrl:      coverURL,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
		UserID:        userId,
		CreateTime:    time.Now().Unix(),
		IsFavorite:    false,
	}
	db.Model(video).Create(video)
	if video.ID == 0 {
		return errors.New("failure in create video record")
	}
	return nil

}

func (d *VideoDao) QueryFeedFlow(requestUID, latestTime int64) []Video {
	//Feed 流接口使用本方法
	var videos []Video
	db.Order("create_time desc").Limit(MaxListLength).Preload("User").Find(&videos)

	for idx, _ := range videos {

		videos[idx].IsFavorite = favoriteDao.Judge(requestUID, videos[idx].ID)
	}

	return videos
}

func packAuthorIn(video_list *[]Video) {
	//给 给定的videoList 中每个video 填装Author
	videoList := *video_list
	if err := db.Model(&videoList).Association("User").Error; err != nil {
		fmt.Println("assosiation wrong with publish list")
	}
	// packing
	for i, _ := range videoList {
		db.Model(&videoList[i]).Association("User").Find(&videoList[i].User)
	}
}

func packFavoriteIn(video_list *[]Video, ownerId int64) {
	//给 给定的videoList 中每个video 填装IsFavorite字段
	videoList := *video_list
	for i, _ := range videoList {
		videoList[i].IsFavorite = favoriteDao.Judge(ownerId, videoList[i].ID)
	}
}
