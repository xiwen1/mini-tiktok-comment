package comment

import (
	"database/sql"
	_ "github.com/lib/pq"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/xiwen1/mini-tiktok-comment/Comment/idl/auth"
	"github.com/xiwen1/mini-tiktok-comment/Comment/idl/comment"
	"github.com/xiwen1/mini-tiktok-comment/Comment/idl/user"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"time"
)

type CommentActionServer struct {
	comment.UnimplementedCommentActionServer
}

var (
	//connStr = "postgres:RpB27iLmDV4z7ZU5tpkn0UPLQWTQx1zFGaUJixDZQhPght7WWLzfZ8PLhZjavGUZ@srv.paracraft.club:31294/nicognaw?sslmode=disable"
	connStr    = "postgres://root:zkw030813@127.0.0.1:5432/root?sslmode=disable"
	pool       *sql.DB
	consul     = "consul://"
	conn       *grpc.ClientConn
	clientAuth auth.AuthServiceClient
	clientUser user.UserServiceClient
)

func InitComment(node int64) error {
	if pool != nil {
		return nil
	}
	var err error
	pool, err = sql.Open("pq", connStr)
	if err != nil {
		log.Fatal("unable to use data source name", err)
		return nil
	}
	conn, err = grpc.Dial(consul, grpc.EmptyDialOption{})
	clientAuth = auth.NewAuthServiceClient(conn)
	clientUser = user.NewUserServiceClient(conn)
	if err != nil {
		log.Fatal(err.Error())
	}
	return nil
}

func CloseComment(ctx context.Context) error {
	if pool == nil {
		return nil
	}
	err := pool.Close()
	if err != nil {
		return err
	}
	pool = nil
	err = conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (CommentActionServer) CommentAction(ctx context.Context, request *comment.CommentActionRequest) (response *comment.CommentActionResponse, err error) {
	token := request.Token
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	authResp, err := clientAuth.Auth(ctx, &auth.AuthRequest{Token: token})
	// todo 检查状态码
	userId := authResp.UserId
	actionType := request.ActionType

	userResp, err := clientUser.GetInfo(ctx, &user.UserInfoRequest{UserId: userId, Token: token})
	// todo 检查状态码
	u := comment.User{
		Id:            userId,
		FollowCount:   userResp.FollowCount,
		FollowerCount: userResp.FollowerCount,
		Name:          userResp.Username,
		IsFollow:      userResp.IsFollow,
	}

	if actionType == comment.CommentActionRequest_PUBLISH {
		c := Comment{}
		c.user = userId
		c.video_id = request.VideoId
		c.content = request.CommentText
		c.createDate = time.Now().Format("01-02")
		err := c.insertComment()
		if err != nil {
			response.StatusCode = comment.CommentActionResponse_FAIL
			response.StatusMsg = "unable to insert into database"
			log.Fatal(err.Error())
		}
		cc := comment.Comment{Id: c.ID, Content: c.content, CreateDate: c.createDate, User: &u}
		response.Comment = &cc
		response.StatusMsg = "success"
		response.StatusCode = comment.CommentActionResponse_SUCCESS
	} else {
		err := delete(request.CommentId)
		if err != nil {
			response.StatusCode = comment.CommentActionResponse_FAIL
			response.StatusMsg = "unable to delete from database"
			log.Fatal(err.Error())
		}
		response.StatusMsg = "success"
		response.StatusCode = comment.CommentActionResponse_SUCCESS
	}
	return
}

func (CommentActionServer) CommentList(ctx context.Context, request *comment.CommentListRequest) (response *comment.CommentListResponse, err error) {
	token := request.Token
	authResp, err := clientAuth.Auth(ctx, &auth.AuthRequest{Token: token})
	if authResp.StatusCode != auth.AuthResponse_SUCCESS {
		response.StatusCode = comment.CommentListResponse_FAIL
		return
	}
	video_id := request.VideoId
	c, err := searchComment(uint32(video_id), token)
	if err != nil {
		response.StatusCode = comment.CommentListResponse_FAIL
		log.Fatal(err.Error())
		return
	}
	response.StatusMsg = "success"
	response.StatusCode = comment.CommentListResponse_SUCCESS
	response.CommentList = c
	// todo 查询结果到comment结构体的map映射
	return
}
