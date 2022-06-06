package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mxysfive/Mini-Tiktok/repository"
	"net/http"
	"strconv"
)

var commentDao = repository.NewCommentDaoInstance()

type CommentResp struct {
	Response
	Comment repository.Comment `json:"comment"`
}
type CommentListResp struct {
	Response
	CommentList []repository.Comment `json:"comment_list"`
}

func CommentAction(c *gin.Context) {
	token := c.Query("token")               //鉴权token
	video_id := c.Query("video_id")         //视频id
	action_type := c.Query("action_type")   //1-发布评论，2-删除评论
	comment_text := c.Query("comment_text") //用户填写的评论内容，在action_type=1的时候使用
	comment_id := c.Query("comment_id")     //要删除的评论id，在action_type=2的时候使用

	if _, exists := onlineUser[token]; !exists {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "not online",
		})
		return
	}
	user := onlineUser[token]
	videoID, err := strconv.ParseInt(video_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "检查输入的视频id",
		})
		return
	}
	commentID, err := strconv.ParseInt(comment_id, 10, 64)
	if err != nil {
		// do nothing
	}
	fmt.Println(comment_text)
	if action_type == "1" {
		//发布评论
		comment := commentDao.CreateOneComment(user, videoID, comment_text)
		c.JSON(http.StatusOK, CommentResp{
			Response: Response{0, "comment success!"},
			Comment:  *comment,
		})
		return
	} else if action_type == "2" {
		//删除评论
		comment := commentDao.DeleteOneComment(commentID)
		c.JSON(http.StatusOK, CommentResp{
			Response: Response{},
			Comment:  *comment,
		})
		return
	} else {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "非法的action_type",
		})
		return
	}
	return

}

func CommentList(c *gin.Context) {
	token := c.Query("token")
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "check video_id",
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
	comments := commentDao.QueryCommentListBy(videoId)
	c.JSON(http.StatusOK, CommentListResp{
		Response:    Response{},
		CommentList: *comments,
	})
}
