package service

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/xiwen1/mini-tiktok-comment/Comment/idl/auth"
	"github.com/xiwen1/mini-tiktok-comment/Comment/idl/comment"
	"github.com/xiwen1/mini-tiktok-comment/Comment/idl/user"
	"google.golang.org/grpc"
	"log"
)

type CommentActionServer struct {
	comment.UnimplementedCommentActionServer
}

var (
	connStrUser = "postgres://postgres:RpB27iLmDV4z7ZU5tpkn0UPLQWTQx1zFGaUJixDZQhPght7WWLzfZ8PLhZjavGUZ@srv.paracraft.club:31294/nicognaw?sslmode=disable"
	connStr     = "postgres://root:zkw030813@127.0.0.1:5432/root?sslmode=disable"
	pool        *sql.DB
	poolUser    *sql.DB
	consulUser  = "consul://0.0.0.0:14514/bawling-minidouyin-user-grpc"
	consulAuth  = "consul://0.0.0.0:14514/bawling-minidouyin-auth-grpc"
	connUser    *grpc.ClientConn
	connAuth    *grpc.ClientConn
	clientAuth  auth.AuthServiceClient
	clientUser  user.UserServiceClient
)

func InitComment() error {
	var err error
	pool, err = sql.Open("pq", connStr)
	if err != nil {
		log.Fatal("unable to use data source name", err)
		return nil
	}
	poolUser, err = sql.Open("pq", connStrUser)
	if err != nil {
		fmt.Println("unable to connect to user")
		return nil
	}
	connUser, err = grpc.Dial(consulUser, grpc.EmptyDialOption{})
	if err != nil {
		log.Fatal(err.Error())
	}
	connAuth, err = grpc.Dial(consulAuth, grpc.EmptyDialOption{})
	if err != nil {
		log.Fatal(err.Error())
	}
	//clientAuth = auth.NewAuthServiceClient(connAuth)
	//clientUser = user.NewUserServiceClient(connUser)
	return nil
}

func CloseComment() error {
	if pool == nil {
		return nil
	}
	err := pool.Close()
	if err != nil {
		return err
	}
	pool = nil
	err = connUser.Close()
	if err != nil {
		return err
	}
	return nil
}
//
//func (CommentActionServer) CommentAction(ctx context.Context, request *comment.CommentActionRequest) (response *comment.CommentActionResponse, err error) {
//	token := request.Token
//	clientAuth := auth.NewAuthServiceClient(connAuth)
//	clientUser := user.NewUserServiceClient(connUser)
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//	authResp, err := clientAuth.Auth(ctx, &auth.AuthRequest{Token: token})
//	if authResp.StatusCode != auth.AuthResponse_SUCCESS {
//		response.StatusCode = comment.CommentActionResponse_FAIL
//		response.StatusMsg = "auth fail"
//		return
//	}
//
//	userId := authResp.UserId
//	actionType := request.ActionType
//
//	userResp, err := clientUser.GetInfo(ctx, &user.UserInfoRequest{UserId: userId, Token: token})
//	if userResp.StatusCode == user.UserInfoResponse_AUTH_FAIL {
//		response.StatusCode = comment.CommentActionResponse_FAIL
//		response.StatusMsg = "userinfo auth fail"
//		return
//	}
//	if userResp.StatusCode == user.UserInfoResponse_UNSPECIFIED {
//		response.StatusCode = comment.CommentActionResponse_FAIL
//		response.StatusMsg = "userid unspecified"
//		return
//	}
//
//	u := comment.User{
//		Id:            userId,
//		FollowCount:   userResp.FollowCount,
//		FollowerCount: userResp.FollowerCount,
//		Name:          userResp.Username,
//		IsFollow:      userResp.IsFollow,
//	}
//
//	if actionType == comment.CommentActionRequest_PUBLISH {
//		c := CommentData{}
//		c.User = userId
//		c.Video_id = request.VideoId
//		c.Content = request.CommentText
//		c.CreateDate = time.Now().Format("01-02")
//		Id, err := c.insertComment()
//		if err != nil {
//			response.StatusCode = comment.CommentActionResponse_FAIL
//			response.StatusMsg = "unable to insert into database"
//			log.Fatal(err.Error())
//		}
//		cc := comment.Comment{Id: Id, Content: c.Content, CreateDate: c.CreateDate, User: &u}
//		response.Comment = &cc
//		response.StatusMsg = "success"
//		response.StatusCode = comment.CommentActionResponse_SUCCESS
//	} else {
//		err := Delete(request.CommentId)
//		if err != nil {
//			response.StatusCode = comment.CommentActionResponse_FAIL
//			response.StatusMsg = "unable to delete from database"
//			log.Fatal(err.Error())
//		}
//		response.StatusMsg = "success"
//		response.StatusCode = comment.CommentActionResponse_SUCCESS
//	}
//	return
//}
//
//func (CommentActionServer) CommentList(ctx context.Context, request *comment.CommentListRequest) (response *comment.CommentListResponse, err error) {
//	token := request.Token
//	clientAuth := auth.NewAuthServiceClient(connAuth)
//	authResp, err := clientAuth.Auth(ctx, &auth.AuthRequest{Token: token})
//	if authResp.StatusCode != auth.AuthResponse_SUCCESS {
//		response.StatusCode = comment.CommentListResponse_FAIL
//		response.StatusMsg = "unable to auth while commentList"
//		return
//	}
//	video_id := request.VideoId
//	c, err := SearchComment(uint32(video_id))
//	if err != nil {
//		response.StatusCode = comment.CommentListResponse_FAIL
//		response.StatusMsg = "unable to search for comment"
//		log.Fatal(err.Error())
//		return
//	}
//	response.StatusMsg = "success"
//	response.StatusCode = comment.CommentListResponse_SUCCESS
//	response.CommentList = c
//
//	return
//}
