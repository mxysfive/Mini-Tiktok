package repository

import (
	"fmt"
	"testing"
	"time"
)

var videoDao = NewVideoDaoInstance()

func TestVideoDao_QueryFeedFlow(t *testing.T) {
	if err := Init(); err != nil {
		panic("connect fail")
	}
	var videos = videoDao.QueryFeedFlow(3, time.Now().Unix())
	for _, v := range videos {
		fmt.Println(v)
	}

}

func TestVideoDao_QueryByOwner(t *testing.T) {
	if err := Init(); err != nil {
		panic("connect fail")
	}

	var videos []Video
	db.Where("user_id = ?", 3).Preload("User").Find(&videos)
	for _, v := range videos {
		fmt.Println(v)
	}
	fmt.Println("--------------")
	videoDao.QueryByOwner(3)
}
