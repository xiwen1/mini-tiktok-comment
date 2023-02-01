package Comment

import (
	"context"
	"database/sql"
	"github.com/xiwen1/mini-tiktok-comment/Comment/idl/comment"
	"github.com/xiwen1/mini-tiktok-comment/Comment/idl/user"
	"log"
)

type Comment struct {
	ID         int64
	user       uint32
	content    string
	createDate string // 格式为 mm-dd
	video_id   int64
}

type User struct {
	ID             uint32
	name           string
	follow_count   int64
	follower_count int64
	is_follow      bool
}

type CommentActionRequest struct {
	token        string
	video_id     int64
	type_code    int
	comment_text string
	comment_id   int64
}

type CommentActionResponse struct {
	status_code int
	status_msg  string
	comment     Comment
}

func delete(commentId int64) error {
	statement, err := pool.Prepare("DELETE FROM public.comment WHERE id=$1")
	checkerr(err)
	_, err = statement.Exec(commentId)
	return err
}

func checkerr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (c *Comment) insertComment() (err error) {
	var statement *sql.Stmt
	statement, err = pool.Prepare("INSERT INTO public.comment(user_id, content, video_id) VALUES ($1,$2, $3);")
	_, err = statement.Exec(c.user, c.content, c.video_id)
	if err != nil {
		println("error")
		log.Fatal(err.Error())
	}
	return
}

func (u *User) insertUser() error {
	statement, err := pool.Prepare("INSERT INTO public.users(name, follow_count, follower_count, is_follow) VALUES ($1, $2, $3, $4);")
	checkerr(err)
	_, err = statement.Exec(u.name, u.follow_count, u.follower_count, u.is_follow)
	checkerr(err)
	return err
}

func searchComment(id uint32, token string) (c []*comment.Comment, err error) {
	var rows *sql.Rows
	var userId uint32
	rows, err = pool.Query("SELECT id, content, create_date, user_id FROM public.comment WHERE id=$1", id)
	for rows.Next() {
		cmt := comment.Comment{}
		err = rows.Scan(&cmt.Id, &cmt.Content, &cmt.CreateDate, &userId)
		ctx := context.Background()
		userResp, err := clientUser.GetInfo(ctx, &user.UserInfoRequest{UserId: userId, Token: token})
		// todo 检查状态码
		u := comment.User{
			Id:            userId,
			FollowCount:   userResp.FollowCount,
			FollowerCount: userResp.FollowerCount,
			Name:          userResp.Username,
			IsFollow:      userResp.IsFollow,
		}
		cmt.User = &u
		if err != nil {
			log.Fatal(err.Error())
		}
		c = append(c, &cmt)
	}
	return
}
