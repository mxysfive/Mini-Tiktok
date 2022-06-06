package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mxysfive/Mini-Tiktok/repository"
	"net/http"
	"strconv"
)

var favoriteDao = repository.NewFavoriteDaoInstance()

type FavoriteListResp struct {
	Response
	VideoList []repository.Video `json:"video_list"`
}

func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")

	if _, exists := onlineUser[token]; !exists {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "not online",
		})
		return
	}
	videoID, err := strconv.ParseInt(video_id, 10, 64)
	if err != nil {
		fmt.Println("videoId is defaultly set by 0")
	}
	user := onlineUser[token]

	if action_type == "1" {
		//表示点赞
		favoriteDao.CreateFavorite(user.ID, videoID)
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "give favorite",
		})
	} else if action_type == "2" {
		//表示取消点赞
		favoriteDao.CancelFavorite(user.ID, videoID)
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "cancel favorite",
		})
	} else {
		//输入了不对的action_type
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Wrong action_type value",
		})
	}

}

func FavoriteList(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	if _, exists := onlineUser[token]; !exists {
		c.JSON(http.StatusOK, Response{
			1,
			"user is not online",
		})
		fmt.Println("wrong in get favorite list")
	}
	userId, _ := strconv.ParseInt(user_id, 10, 64)
	var VideoList = favoriteDao.FavoriteList(userId)
	c.JSON(http.StatusOK, FavoriteListResp{
		Response:  Response{},
		VideoList: VideoList,
	})
}
