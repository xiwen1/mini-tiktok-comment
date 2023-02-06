package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/xiwen1/mini-tiktok-comment/Comment/idl/auth"
	"github.com/xiwen1/mini-tiktok-comment/Comment/idl/user"
	"log"
	"time"
)

type CommentData struct {
	ID         int64
	User       uint32
	Content    string
	Video_id   int64
	CreateDate string // 格式为 mm-dd
}

type CommentInfo struct {
	Id         int64  `json:"id"`
	User       User   `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}
type User struct {
	ID            uint32
	name          string
	followCount   uint32
	followerCount uint32
	isFollow      bool
}

type CommentActionRequest struct {
	token        string
	video_id     int64
	type_code    int
	comment_text string
	comment_id   int64
}

//type CommentActionResponse struct {
//	status_code int
//	status_msg  string
//	comment     CommentData
//}

func Delete(commentId int64) error {
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

func (c *CommentData) insertComment() (lastInserId int64, err error) {
	var statement *sql.Stmt
	statement, err = pool.Prepare("INSERT INTO public.comment(user_id, content, video_id) VALUES ($1,$2, $3);")
	result, err := statement.Exec(c.User, c.Content, c.Video_id)
	if err != nil {
		println("error")
		log.Fatal(err.Error())
	}
	lastInserId, err = result.LastInsertId()
	return
}

func (u *User) insertUser() (int64, error) {
	statement, err := pool.Prepare("INSERT INTO public.users(name, follow_count, follower_count, is_follow) VALUES ($1, $2, $3, $4);")
	checkerr(err)
	result, err := statement.Exec(u.name, u.followCount, u.followerCount, u.isFollow)
	checkerr(err)
	lastInsertId, err := result.LastInsertId()
	return lastInsertId, err
}

func Send(data CommentData, token string) (CommentInfo, error) {
	fmt.Println("service: send is running")
	Id, err := data.insertComment()
	if err != nil {
		return CommentInfo{}, err
	}
	clientUser := user.NewUserServiceClient(connUser)
	ctx := context.Background()
	userResp, err := clientUser.GetInfo(ctx, &user.UserInfoRequest{UserId: data.User, Token: token})
	if userResp.StatusCode != user.UserInfoResponse_AUTH_FAIL {
		return CommentInfo{}, err
	}
	u := User{
		ID:            data.User,
		followCount:   userResp.FollowCount,
		followerCount: userResp.FollowerCount,
		name:          userResp.Username,
		isFollow:      userResp.IsFollow,
	}

	return CommentInfo{
		User:       u,
		CreateDate: data.CreateDate,
		Content:    data.Content,
		Id:         Id,
	}, err

}

func Auth(token string) (statusCode int, userId uint32, err error) {
	clientAuth := auth.NewAuthServiceClient(connAuth)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	authResp, err := clientAuth.Auth(ctx, &auth.AuthRequest{Token: token})
	return int(authResp.StatusCode), authResp.UserId, err
}


func SearchComment(videoId uint32, token string) (c *[]CommentInfo, err error) {
	var rows *sql.Rows
	var userId uint32
	ctx := context.Background()
	rows, err = pool.Query("SELECT id, content, create_date, user_id FROM public.comment WHERE video_id=$1", videoId)
	clientAuth := auth.NewAuthServiceClient(connAuth)
	authResp, err := clientAuth.Auth(ctx, &auth.AuthRequest{Token: token})
	if authResp.StatusCode == auth.AuthResponse_SUCCESS {
		for rows.Next() {
			cmt := CommentInfo{}
			err = rows.Scan(&cmt.Id, &cmt.Content, &cmt.CreateDate, &userId)
			clientUser := user.NewUserServiceClient(connUser)
			userResp, err := clientUser.GetInfo(ctx, &user.UserInfoRequest{UserId: userId, Token: token})
			if err != nil {
				log.Fatal(err.Error())
				return nil, err
			}
			statusCode := int(userResp.StatusCode)
			if statusCode != 1{
				u := User{
					ID:            userId,
					followCount:   userResp.FollowCount,
					followerCount: userResp.FollowerCount,
					name:          userResp.Username,
					isFollow:      userResp.IsFollow,
				}
				cmt.User = u
				*c = append(*c, cmt)
			} else {
				fmt.Println("getting userInfo fail")
				return
			}
		}
	} else {
		for rows.Next() {
			cmt := CommentInfo{}
			err = rows.Scan(&cmt.Id, &cmt.Content, &cmt.CreateDate, &userId)
			if err != nil {
				log.Fatal(err.Error())
				return nil, err
			}
			var u User
			err = poolUser.QueryRow("SELECT id, follow_count, follower_count, name FROM public.users WHERE id=$1", userId,
			).Scan(&u.ID, &u.followCount, &u.followerCount, &u.name)
			if err != nil {
				fmt.Println("query for userinfo from db fail")
			}
			u.isFollow = false
			cmt.User = u
			*c = append(*c, cmt)

		}
	}

	return
}
