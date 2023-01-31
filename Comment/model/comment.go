package comment

import (
	"database/sql"
	_ "github.com/lib/pq"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/xiwen1/mini-tiktok-comment/Comment/idl/comment"
	"golang.org/x/net/context"
	"log"
	"time"
)

type CommentActionServer struct {
	comment.UnimplementedCommentActionServer
}

var (
	//connStr = "postgres:RpB27iLmDV4z7ZU5tpkn0UPLQWTQx1zFGaUJixDZQhPght7WWLzfZ8PLhZjavGUZ@srv.paracraft.club:31294/nicognaw?sslmode=disable"
	connStr = "postgres://root:zkw030813@127.0.0.1:5432/root?sslmode=disable"
	pool    *sql.DB
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
	return nil
}

func (CommentActionServer) CommentAction(ctx context.Context, request *comment.CommentActionRequest) (response *comment.CommentActionResponse, err error) {
	//token := request.Token

	// todo 检查token
	actionType := request.ActionType

	if actionType == comment.CommentActionRequest_PUBLISH {
		c := Comment{}
		// todo 查询user信息
		c.video_id = request.VideoId
		c.content = request.CommentText
		c.createDate = time.Now().Format("01-02")
		err := c.insertComment()
		if err != nil {
			response.StatusCode = comment.CommentActionResponse_FAIL
			response.StatusMsg = "unable to insert into database"
			log.Fatal(err.Error())
		}
		cc := comment.Comment{}
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

func (CommentActionServer) CommentList(ctx context.Context, request *comment.CommentListRequest) (response *comment.CommentActionResponse, err error) {
	return
}

//func main() {
//	var err error
//	pool, err = sql.Open("postgres", connStr)
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//	ctx := context.Background()
//	err = pool.PingContext(ctx)
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//	println("connected")
//}
