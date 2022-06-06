package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mxysfive/Mini-Tiktok/repository"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type PublishListResp struct {
	Response
	VideoList []repository.Video `json:"video_list"`
}

var videoDao = repository.NewVideoDaoInstance()

// ResourceBase 如果映射的域名和改了，需要更改这个配置
const ResourceBase = "http://5kt3855788.zicp.vip/static"

func PublishVideo(c *gin.Context) {
	token := c.PostForm("token")
	title := c.PostForm("title")
	if _, exists := onlineUser[token]; !exists {
		fmt.Println("not exist?")
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user is not exists",
		})
		return
	}
	data, err := c.FormFile("data")

	if err != nil {
		c.JSON(http.StatusOK, Response{
			1,
			err.Error(),
		})
		return
	}
	filename := data.Filename
	user := onlineUser[token]
	finalName := fmt.Sprintf("%d_%s", user.ID, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "upload success",
		})
	}
	playURL, vErr := joinResourceURL(ResourceBase, finalName)
	if vErr != nil {
		fmt.Printf("wrong join URL")
		fmt.Printf("Wrong URL is: %s", playURL)
	}
	coverURL, cErr := joinResourceURL(ResourceBase, "wetcar.jpg") //test
	if cErr != nil {
		fmt.Printf("wrong join URL")
		fmt.Printf("Wrong URL is: %s", coverURL)
	}
	if err := videoDao.CreateVideoRecord(user.ID, playURL, coverURL, title); err != nil {
		fmt.Println("Error in create video record")
	}
	return
}

func joinResourceURL(baseDomain, resource string) (string, error) {
	// 返回拼接好的视频或封面的URL
	var sb strings.Builder
	_, err := fmt.Fprintf(&sb, "%s/%s", baseDomain, resource)
	if err != nil {
		fmt.Printf("joinResource fail %v", err)
		return "", err
	}
	return sb.String(), nil
}

func PublishList(c *gin.Context) {
	id := c.Query("user_id")
	//因为可以查看别人发出的视频，所以user_id 可能是别人的
	utoken := c.Query("token")
	if _, exists := onlineUser[utoken]; !exists {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user is not online",
		})
		return
	}
	uid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		fmt.Println("wrong uid in publish list")
	}

	user, err := userDaoInstance.QueryUserById(uid)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			1,
			"user is not exists in database",
		})
		return
	}
	var PublishedList = videoDao.QueryByOwner(user.ID)

	c.JSON(http.StatusOK, PublishListResp{
		Response:  Response{0, "Query success!"},
		VideoList: PublishedList,
	})
	return

}
