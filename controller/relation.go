package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mxysfive/Mini-Tiktok/repository"
	"net/http"
	"strconv"
)

type RelationResp struct {
	Response
	UserList []repository.User `json:"user_list"`
}

var relationDao = repository.NewRelationDaoInstance()

func FollowAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")
	toUserID, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "check input to_user_id",
		})
		return
	}
	if _, exists := onlineUser[token]; !exists {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "not online",
		})
		return
	}

	user := onlineUser[token]
	if actionType == "1" {
		//点关注
		relationDao.Follow(user.ID, toUserID)
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "ok!",
		})
	} else if actionType == "2" {
		//取消关注
		relationDao.UnFollow(user.ID, toUserID)
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "ok!",
		})
	} else {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "illegal action_type",
		})
	}
	return
}

func FollowList(c *gin.Context) {
	token := c.Query("token")
	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "check input to_user_id",
		})
		return
	}
	if _, exists := onlineUser[token]; !exists {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "not online",
		})
		return
	}

	var userList = relationDao.FollowListOf(userID)
	c.JSON(http.StatusOK, RelationResp{
		Response: Response{},
		UserList: *userList,
	})
}

func FollowedList(c *gin.Context) {
	token := c.Query("token")
	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "check input to_user_id",
		})
		return
	}
	if _, exists := onlineUser[token]; !exists {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "not online",
		})
		return
	}

	var userList = relationDao.FansOf(userID)
	c.JSON(http.StatusOK, RelationResp{
		Response: Response{},
		UserList: *userList,
	})
}
