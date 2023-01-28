package model

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	pb "github.com/xiwen1/mini-tiktok-comment/Comment/idl"
	"github.com/xiwen1/mini-tiktok-comment/Comment/snowflake"
	"golang.org/x/net/context"
	"sync"
)

type Comment struct {
	ID         int64
	user       pb.User
	content    string
	createDate string // 格式为 mm-dd
}

type CommentActionServer struct {
	pb.UnimplementedCommentActionServer
	mu sync.Mutex
}

var (
	db neo4j.DriverWithContext
	sf *snowflake.Worker
)

func InitComment(uri, user, password, realm string, node int64) error {
	if db != nil {
		return nil
	}
	var err error
	db, err = neo4j.NewDriverWithContext(
		uri,
		neo4j.BasicAuth(
			user,
			password,
			realm,
		),
	)
	if err != nil {
		return err
	}

	sf, err = snowflake.NewWorker(node)
	if err != nil {
		return err
	}
	return nil
}

func commentToMap(c *Comment) map[string]interface{} {
	return map[string]interface{}{
		"id":          c.ID,
		"user":        c.user.Id,
		"content":     c.content,
		"create_date": c.createDate,
	}
}

func CloseComment(ctx context.Context) error {
	if db == nil {
		return nil
	}
	err := db.Close(ctx)
	if err != nil {
		return err
	}
	db = nil
	return nil
}

func (CommentActionServer) CommentAction(ctx context.Context, req *pb.CommentActionRequest) *pb.CommentActionResponse {
	token := req.Token

	// todo 检查token
	actionType := req.ActionType

	session := db.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	if actionType == pb.CommentActionRequest_PUBLISH {
		id := sf.GetId()

	}
}
