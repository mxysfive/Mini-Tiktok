package repository

import (
	"fmt"
	"testing"
)

var fDao = NewFavoriteDaoInstance()

func TestFavoriteDao_FavoriteList(t *testing.T) {
	if err := Init(); err != nil {
		panic("connect fail")
	}
	var videos []Video
	var selectedAttr = []string{"videos.id", "play_url", "cover_url", "favorite_count", "comment_count", "title",
		"favorites.is_favorite"}
	db.Model(&videos).Select(selectedAttr).Joins("left join favorites on favorites.video_id = videos.id").
		Joins("User").
		Where("favorites.user_id = ?", 3).
		Where("favorites.is_favorite = ?", true).Find(&videos)

	for _, video := range videos {
		fmt.Println(video)
	}
}

func TestCommentDao_QueryCommentByPKs(t *testing.T) {
	if err := Init(); err != nil {
		panic("connect fail")
	}
	db.AutoMigrate(&Comment{})
}
