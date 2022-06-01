package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mxysfive/Mini-Tiktok/repository"
	"net/http"
	"strconv"
	"strings"
)

var onlineUser = map[string]*User{}

const MaxUsernameLen = 32
const MaxPasswordLen = 32

type RegisterResp struct {
	Response
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}
type UserResp struct {
	Response
	User User `json:"user"`
}

var userDaoInstance = repository.NewUserDaoInstance()

func Register(c *gin.Context) {
	userName := c.Query("username")
	password := c.Query("password")
	uErr := checkUserName(userName)
	pErr := checkPassword(password)
	if uErr != nil {
		c.JSON(http.StatusOK, RegisterResp{
			Response: Response{StatusCode: 1,
				StatusMsg: uErr.Error()},
		})
		return
	}
	if pErr != nil {
		c.JSON(http.StatusOK, RegisterResp{
			Response: Response{StatusCode: 1,
				StatusMsg: pErr.Error()},
		})
		return
	}
	//调用Dao层
	var user = userDaoInstance.CreateByNameAndPassword(userName, password)

	var tokenSb strings.Builder
	fmt.Fprintf(&tokenSb, "%s%s", userName, password)
	c.JSON(http.StatusOK, RegisterResp{
		Response: Response{
			StatusCode: 0,
		},
		UserId: user.ID, //不知道该怎么写了
		Token:  tokenSb.String(),
	})
	return
}

func Login(c *gin.Context) {
	userName := c.Query("username")
	password := c.Query("password")
	var user = userDaoInstance.QueryLoginInfo(userName, password)
	var tokenSb strings.Builder
	fmt.Fprintf(&tokenSb, "%s%s", userName, password)
	c.JSON(http.StatusOK, RegisterResp{
		Response: Response{0, ""},
		UserId:   user.ID,
		Token:    tokenSb.String(),
	})

	var loginUser = &User{
		UserId:        user.ID,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      false,
	}

	//加入到online表里
	onlineUser[tokenSb.String()] = loginUser
}

func UserInfo(c *gin.Context) {
	qid := c.Query("user_id")
	utoken := c.Query("token") //不知道utoken有什么用
	if _, exists := onlineUser[utoken]; !exists {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user is not online",
		})
		return
	}

	userId, err := strconv.Atoi(qid)
	if err != nil {
		fmt.Printf("Function of atoi in UserInfo fail %v", err)
	}
	var userEntity = userDaoInstance.QueryUserById(int64(userId))
	fmt.Println("entity is: ", userEntity)
	loginUser := &User{
		UserId:        userEntity.ID,
		Name:          userEntity.Name,
		FollowCount:   userEntity.FollowCount,
		FollowerCount: userEntity.FollowerCount,
		IsFollow:      false,
	}
	c.JSON(http.StatusOK, UserResp{
		Response: Response{0, ""},
		User:     *loginUser,
	})

	return
}

func checkUserName(userName string) error {
	if len(userName) > MaxUsernameLen {
		return errors.New("username is too long")
	}

	return nil
}

func checkPassword(passWord string) error {
	if len(passWord) > MaxPasswordLen {
		return errors.New("password is too long")
	}
	return nil
}
