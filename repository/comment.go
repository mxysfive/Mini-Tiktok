package repository

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID         int64  `gorm:"primaryKey" json:"id"`
	Content    string `gorm:"size:255" json:"content"`
	CreateDate string `gorm:"size:30" json:"create_date"`
	UserID     int64  `json:"user_id"`
	VideoID    int64  `json:"video_id"`
	User       User   `json:"user"`
	Video      Video
}

type CommentDao struct {
}

func NewCommentDaoInstance() *CommentDao {
	return &CommentDao{}
}

func (c *CommentDao) CreateOneComment(user *User, videoId int64, text string) *Comment {
	now := time.Now()
	var comment = &Comment{
		ID:         0,
		Content:    text,
		CreateDate: now.Format("2006-01-02"),
		UserID:     user.ID,
		VideoID:    videoId,
		User:       User{},
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(comment).Create(comment).Error; err != nil {
			return err
		}
		comment.User = *user
		if err := tx.Model(&Video{}).
			Where("id = ?", videoId).
			Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {

			return err
		}
		return nil
	})
	if err != nil {
		return &Comment{}
	}
	return comment
}

func (c *CommentDao) DeleteOneComment(commentId int64) *Comment {
	var comment = &Comment{}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(comment, commentId).Error; err != nil {
			return err
		}

		if err := tx.Delete(&Comment{}, commentId).Error; err != nil {
			return err
		}

		if err := tx.Model(&Video{}).
			Where("id = ?", comment.VideoID).
			Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err != nil {

			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	return comment
}

func (c *CommentDao) QueryCommentListBy(videoId int64) *[]Comment {
	//根据给定的videoId 找到Comment
	var comments = make([]Comment, 50)
	db.Model(&comments).Preload("User").Where("video_id = ?", videoId).Find(&comments)

	return &comments

}
